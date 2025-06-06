package services

import (
	"RobotService/internal/dto"
	"RobotService/internal/entities"
	"RobotService/internal/rabbit"
	"RobotService/internal/repositories"
	"RobotService/internal/sorrage"
	"fmt"

	"log"
	"strconv"
	"time"
)

const (
	keyadd         = "robots.Add"
	keyget         = "robots.Get"
	keyupdatecords = "robots.UpdateCord"
	keyupdatename  = "robots.UpdateName"
	keyupdatetype  = "robots.UpdateType"
	keydel         = "robots.Del"
)

type RbtSrvic struct {
	RobotRepository repositories.RobotRepositories
	Redis           *sorrage.RdsCache
	Rabbit          *rabbit.Publisher
}

func (srvc *RbtSrvic) CreateRobot(dto dto.CreateRobotDTO) (int, error) {
	robot := entities.Robot{
		Name:  dto.Name,
		Type:  dto.Type,
		XCord: dto.XCord,
		YCord: dto.YCord,
		ZCord: dto.ZCord,
	}
	createdRobot, err := srvc.RobotRepository.CreateRobot(robot)
	// После создания робота закидываем его данные в редиску
	_ = srvc.Redis.SetRobotData(strconv.Itoa(createdRobot.ID), createdRobot, 5*time.Minute)
	srvc.publishToRabbitWithStruct(&createdRobot, keyadd)
	return createdRobot.ID, err
}

func (serv *RbtSrvic) GetRobotInfo(id int) (*entities.Robot, error) {
	idStr := strconv.Itoa(id)
	// Пытаемся получить данные из кэша, если они есть - получаем ошибку и идём дальше по коду, если данные есть то ретёрним их
	robotdata, err := serv.Redis.GetRobotData(idStr)
	if err == nil {
		serv.publishToRabbitWithStruct(robotdata, keyget)
		return robotdata, nil
	}
	// Если данных в кэше нет, то обращаемся к репозиторию и получаем данные из БД
	robotdata, err = serv.RobotRepository.GetRobotInfo(id)
	if err != nil {
		return nil, err
	}
	// Добавляем полученные данные в кэш на пять минут
	err = serv.Redis.SetRobotData(idStr, *robotdata, 5*time.Minute)
	serv.publishToRabbitWithStruct(robotdata, keyget)
	return robotdata, err
}

func (srv *RbtSrvic) UpdateRobotCords(updateData dto.UpdateRobotCordDTO) error {
	newCord := entities.RobotCord{XCord: updateData.XCord, YCord: updateData.YCord, ZCord: updateData.ZCord}
	robotID := updateData.ID
	err := srv.RobotRepository.UpdateRobotCords(robotID, newCord)
	if err != nil {
		return err
	}
	// Удаление кэша после обновления
	_ = srv.Redis.DeleteRobotData(strconv.Itoa(robotID))
	msgToRabbit := fmt.Sprintf("Координаты робота с ID: %d были изменены на  X:%d, Y:%d, Z:%d", robotID, newCord.XCord, newCord.YCord, newCord.ZCord)
	srv.publishToRabbitWithText(msgToRabbit, keyupdatecords)
	return nil
}

func (sv *RbtSrvic) UpdateRobotName(updateData dto.UpdateRobotNameDTO) error {
	newName := updateData.Name
	robotID := updateData.ID
	err := sv.RobotRepository.UpdateRobotName(robotID, newName)
	if err != nil {
		return err
	}
	// Удаление кэша после обновления
	_ = sv.Redis.DeleteRobotData(strconv.Itoa(robotID))
	msgToRabbit := fmt.Sprintf("Имя робота с ID: %d было изменены на %s", robotID, newName)
	sv.publishToRabbitWithText(msgToRabbit, keyupdatename)
	return nil
}

func (ssrv *RbtSrvic) ChangeRobotType(updateData dto.ChangeTypeDTO) error {
	newType := updateData.Type
	robotID := updateData.ID
	err := ssrv.RobotRepository.ChangeRobotType(robotID, newType)
	if err != nil {
		return err
	}
	// Удаление кэша после обновления
	_ = ssrv.Redis.DeleteRobotData(strconv.Itoa(robotID))
	msgToRabbit := fmt.Sprintf("Тип робота с ID: %d был изменен на %s", robotID, newType)
	ssrv.publishToRabbitWithText(msgToRabbit, keyupdatetype)
	return nil
}

func (srv *RbtSrvic) DeleteRobot(id int) error {
	msgToRabbit := fmt.Sprintf("Робота с ID: %d был уничтожен. Помянем...", id)
	srv.publishToRabbitWithText(msgToRabbit, keydel)
	_ = srv.Redis.DeleteRobotData(strconv.Itoa(id))
	return srv.RobotRepository.DeleteRobot(id)
}

// Отправка в реббит сообщения со струтурой робота
func (srv *RbtSrvic) publishToRabbitWithStruct(robot *entities.Robot, routingKey string) {
	if err := srv.Rabbit.Publish(robot, routingKey); err != nil {
		log.Println("Не получилось отправить, сорян")
	}
}

// Отправка в реббит сообщения с текстом
func (srv *RbtSrvic) publishToRabbitWithText(msg string, routingKey string) {
	if err := srv.Rabbit.PublishWithText(msg, routingKey); err != nil {
		log.Println("Не получилось отправить, сорян")
	}
}
