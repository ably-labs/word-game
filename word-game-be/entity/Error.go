package entity

type Error struct {
	Err string `json:"err"`
}

var (
	ErrInvalidInput     = Error{Err: "Invalid Input"}
	ErrBadSession       = Error{Err: "Bad Session"}
	ErrUserCreateFailed = Error{Err: "Could not create user"}
	ErrUserExists       = Error{Err: "A user with that name already exists"}
	ErrUserDoesNotExist = Error{Err: "User does not exist, you must register first"}
	ErrDatabaseError    = Error{Err: "Database Error"}
	ErrForbidden        = Error{Err: "Forbidden"}
	ErrUnauthorised     = Error{Err: "Unauthorised"}
	ErrInvalidLobby     = Error{Err: "Invalid Lobby ID"}
	ErrLobbyNotFound    = Error{Err: "Lobby Not Found"}
	ErrLobbyNotJoined   = Error{Err: "You are not a member of this lobby"}
	ErrLobbyJoined      = Error{Err: "You are already a member of this lobby"}
	ErrLobbyFull        = Error{Err: "Lobby is full"}
	ErrLoggedIn         = Error{Err: "You must be logged out to use this endpoint"}
	ErrSpectating       = Error{Err: "You must be an active participant in the game to use this endpoint"}
	ErrNotYourTurn      = Error{Err: "It is not currently your turn"}
	ErrTileOccupied     = Error{Err: "Tile is already occupied"}
	ErrInvalidPlay      = Error{Err: "Invalid tile placement"}
	ErrGameInProgress   = Error{Err: "Game is already in progress"}
)
