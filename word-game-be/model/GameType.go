package model

type BonusTilePattern string

const (
	None    BonusTilePattern = "none"
	Regular BonusTilePattern = "regular"
	Stripe  BonusTilePattern = "stripe"
	Border  BonusTilePattern = "border"
)

type GameType struct {
	ID   uint32 `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
	// The number of tiles in the bag at the start of the game
	TileBagSize int `gorm:"default:100" json:"tileBagSize"`
	// The size of a player's tile deck
	PlayerDeckSize int `gorm:"default:9" json:"playerDeckSize"`
	// The number of tiles a player gets replenished to at the end of a turn
	PlayerTileCount int `gorm:"default:7" json:"playerTileCount"`
	// If blank tiles are given out
	EnableBlankTiles bool `gorm:"default:true" json:"enableBlankTiles"`
	// The width of the main board
	BoardWidth int `gorm:"default:15" json:"boardWidth"`
	// The height of the main board
	BoardHeight int `gorm:"default:15" json:"boardHeight"`
	// The pattern to generate bonus tiles in
	BonusTilePattern BonusTilePattern `gorm:"default:'regular'" json:"bonusTilePattern"`
	// If the first player has to start in the center
	StartAnywhere bool `gorm:"default:false" json:"startAnywhere"`
	// If the GameType is selectable by others
	Visible bool `gorm:"default:true" json:"-"`
}
