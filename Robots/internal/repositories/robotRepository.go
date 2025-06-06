package repositories

import (
	"RobotService/internal/entities"
	"context"

	"github.com/jackc/pgx/v5"
)

type RobotRepositories struct {
	DataBase *pgx.Conn
}

func (repo *RobotRepositories) CreateRobot(robot entities.Robot) (entities.Robot, error) {
	query := "INSERT INTO robots (name, xcord, ycord, zcord) VALUES($1, $2, $3, $4) returning id"
	err := repo.DataBase.QueryRow(context.Background(), query, robot.Name, robot.XCord, robot.YCord, robot.ZCord).Scan(&robot.ID)
	if err != nil {
		return robot, err
	}
	return robot, nil
}

func (repo *RobotRepositories) GetRobotInfo(id int) (*entities.Robot, error) {
	robot := &entities.Robot{ID: id}
	query := "SELECT name, xcord, ycord, zcord FROM robots WHERE id = $1"
	err := repo.DataBase.QueryRow(context.Background(), query, id).Scan(&robot.Name, &robot.XCord, &robot.YCord, &robot.ZCord)
	if err != nil {
		return nil, err
	}
	return robot, nil
}

func (repo *RobotRepositories) UpdateRobotCords(id int, newCords entities.RobotCord) error {
	query := "UPDATE robots SET xcord = $1, ycord = $2, zcord = $3 WHERE id = $4"
	_, err := repo.DataBase.Exec(context.Background(), query, newCords.XCord, newCords.YCord, newCords.ZCord, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *RobotRepositories) UpdateRobotName(id int, newName string) error {
	query := "UPDATE robots SET name = $1 WHERE id = $2"
	_, err := repo.DataBase.Exec(context.Background(), query, newName, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *RobotRepositories) DeleteRobot(id int) error {
	query := "DELETE FROM robots WHERE id = $1"
	_, err := repo.DataBase.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
