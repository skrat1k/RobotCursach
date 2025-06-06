package dto

type CreateRobotDTO struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	XCord int    `json:"xCord"`
	YCord int    `json:"yCord"`
	ZCord int    `json:"zCord"`
}
