package controller

import (
	_ "embed"
	"fmt"
	"github.com/ably-labs/word-game/word-game-be/constant"
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/middleware"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/ably-labs/word-game/word-game-be/util"
	"github.com/ably/ably-go/ably"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"sort"
)

type GameController struct {
	db   *gorm.DB
	ably *ably.Realtime
}

func NewGameController(e *echo.Group, db *gorm.DB, ably *ably.Realtime) *GameController {
	gc := &GameController{
		db:   db,
		ably: ably,
	}

	g := e.Group("/:id", middleware.RequireAuthorisation, middleware.ValidateLobby, middleware.RequireLobbyMember)

	g.GET("/boards", gc.GetBoards)
	g.GET("/boards/deck", gc.GetDeck)
	g.PATCH("/boards", gc.PatchBoard)
	g.POST("/boards", gc.EndTurn, middleware.RequireTurn)

	return gc
}

func (gc *GameController) GetBoards(c echo.Context) error {
	lobby := c.Get("lobby").(*model.Lobby)

	boards := make(map[string]entity.SquareSet)
	boards["main"] = lobby.Board

	lobbyMember := c.Get("lobbyMember").(*model.LobbyMember)
	if lobbyMember.MemberType != constant.MemberTypeSpectator {
		boards["deck"] = lobbyMember.Deck
	}

	return c.JSON(200, boards)
}

func (gc *GameController) GetDeck(c echo.Context) error {
	lobbyMember := c.Get("lobbyMember").(*model.LobbyMember)
	if lobbyMember.MemberType == constant.MemberTypeSpectator {
		return c.JSON(403, entity.ErrSpectating)
	}
	return c.JSON(200, lobbyMember.Deck)
}

func (gc *GameController) PatchBoard(c echo.Context) error {
	moveTile := entity.MoveTile{}
	err := c.Bind(&moveTile)
	if err != nil {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	lobby := c.Get("lobby").(*model.Lobby)
	lobbyMember := c.Get("lobbyMember").(*model.LobbyMember)

	var fromBoard entity.SquareSet
	var toBoard entity.SquareSet

	boardUpdate := false

	if moveTile.From == "main" {
		fromBoard = lobby.Board
		boardUpdate = true
	} else {
		fromBoard = lobbyMember.Deck
	}

	if moveTile.To == "main" {
		toBoard = lobby.Board
		boardUpdate = true
	} else {
		toBoard = lobbyMember.Deck
	}

	// Disallow this move if it affects the board when it's not the users turn
	if boardUpdate && *lobby.PlayerTurnID != lobbyMember.UserID {
		fmt.Println(lobby.PlayerTurnID, lobbyMember.UserID)
		return c.JSON(403, entity.ErrNotYourTurn)
	}

	// Make sure the tiles actually exist
	if !validatePos(fromBoard, moveTile.FromIndex) || !validatePos(toBoard, moveTile.ToIndex) {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	// No tile exists where we're taking from
	if (*fromBoard.Squares)[moveTile.FromIndex].Tile == nil {
		return c.JSON(400, entity.ErrInvalidPlay)
	}

	// Trying to use a non-blank tile as a blank tile
	if moveTile.Letter != "" && !(*fromBoard.Squares)[moveTile.FromIndex].Tile.Blank {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	if (*toBoard.Squares)[moveTile.ToIndex].Tile != nil {
		return c.JSON(400, entity.ErrTileOccupied)
	}

	(*toBoard.Squares)[moveTile.ToIndex].Tile = (*fromBoard.Squares)[moveTile.FromIndex].Tile

	if moveTile.Letter != "" {
		(*toBoard.Squares)[moveTile.ToIndex].Tile.Letter = moveTile.Letter
	}

	(*fromBoard.Squares)[moveTile.FromIndex].Tile = nil

	if boardUpdate {
		fmt.Println("Updating board")
		_ = util.PublishLobbyMessage(gc.ably, lobby.ID, "moveTile", map[string]interface{}{
			"move": moveTile,
			"tile": (*toBoard.Squares)[moveTile.ToIndex].Tile,
		})
	}

	gc.db.Save(&lobby)
	gc.db.Save(&lobbyMember)

	return c.NoContent(204)
}

func validatePos(board entity.SquareSet, position int) bool {
	return position >= 0 && position < len(*board.Squares)
}

func (gc *GameController) EndTurn(c echo.Context) error {
	lobby := c.Get("lobby").(*model.Lobby)
	lobbyMember := c.Get("lobbyMember").(*model.LobbyMember)

	score, err := util.ValidateBoard(lobby.Board)

	if err != nil {
		return c.JSON(400, entity.Error{Err: err.Error()})
	}

	gc.db.Order("joined_at").Where("member_type = 'player' AND lobby_id = ?", lobby.ID).Find(&lobby.Members)

	// don't bother updating anything if this was a pass
	if score > 0 {
		lobbyMember.Score += score
		_ = util.PublishLobbyMessage(gc.ably, lobby.ID, "message", entity.ChatSent{
			Message: fmt.Sprintf("<@%d> scored %d points", lobbyMember.UserID, score),
			Author:  "system",
		})
		_ = util.PublishLobbyMessage(gc.ably, lobby.ID, "scoreUpdate", lobbyMember)
		remainingTiles := lobbyMember.Deck.TileCount()

		if len(*lobby.Bag.Squares) > 0 {
			lobbyMember.Deck.AddTiles(util.TakeFromBag(7-remainingTiles, &lobby.Bag))
		} else if remainingTiles == 0 {
			lobby.State = entity.LobbyRoundOver
			sort.Slice(lobby.Members, func(i, j int) bool {
				return lobby.Members[i].Score < lobby.Members[j].Score
			})
			// Set winner to current turn
			*lobby.PlayerTurnID = lobby.Members[0].UserID
			_ = util.PublishLobbyMessage(gc.ably, lobby.ID, "message", entity.ChatSent{
				Message: fmt.Sprintf("Game over, <@%d> wins!", *lobby.PlayerTurnID),
				Author:  "system",
			})
		}

		indices := util.GetPlacedTileIndices(lobby.Board)
		for _, index := range indices {
			(*lobby.Board.Squares)[index].Tile.Draggable = false
		}

	} else {
		_ = util.PublishLobbyMessage(gc.ably, lobby.ID, "message", entity.ChatSent{
			Message: fmt.Sprintf("<@%d> passed", lobbyMember.UserID),
			Author:  "system",
		})
	}

	// If we're still in play, set the next turn
	if lobby.State == entity.LobbyInGame {
		for i, member := range lobby.Members {
			if member.UserID == *lobby.PlayerTurnID {
				if i == len(lobby.Members)-1 {
					*lobby.PlayerTurnID = lobby.Members[0].UserID
				} else {
					*lobby.PlayerTurnID = lobby.Members[i+1].UserID
				}
				break
			}
		}
		fmt.Println("It is now this players turn ", *lobby.PlayerTurnID)
	}

	_ = util.PublishLobbyMessage(gc.ably, lobby.ID, "lobbyUpdate", lobby)

	gc.db.Save(&lobbyMember)
	gc.db.Save(&lobby)

	return c.NoContent(204)
}
