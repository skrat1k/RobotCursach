package entities

type Robot struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	XCord int    `json:"xCord"`
	YCord int    `json:"yCord"`
	ZCord int    `json:"zCord"`
}
