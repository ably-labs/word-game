package entity

type MoveTile struct {
	From      string `json:"from"`
	FromIndex int    `json:"fromIndex"`
	To        string `json:"to"`
	ToIndex   int    `json:"toIndex"`
	Letter    string `json:"letter,omitempty"`
}
