package entities

type Robot struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	XCord int    `json:"xCord"`
	YCord int    `json:"yCord"`
	ZCord int    `json:"zCord"`
}
