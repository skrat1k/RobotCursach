package dto

type UpdateRobotCordDTO struct {
	ID    int `json:"id"`
	XCord int `json:"xCord"`
	YCord int `json:"yCord"`
	ZCord int `json:"zCord"`
}
