package services

import (
	"RobotService/internal/dto"
	"RobotService/internal/entities"
	"RobotService/internal/rabbit"
	"RobotService/internal/repositories"
	"RobotService/internal/storage"
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

type RobotService struct {
	RobotRepository repositories.RobotRepositories
	Redis           *storage.RedisCache
	Rabbit          *rabbit.Publisher
}

func (service *RobotService) CreateRobot(dto dto.CreateRobotDTO) (int, error) {
	robot := entities.Robot{
		Name:  dto.Name,
		Type:  dto.Type,
		XCord: dto.XCord,
		YCord: dto.YCord,
		ZCord: dto.ZCord,
	}
	createdRobot, err := service.RobotRepository.CreateRobot(robot)
	// После создания робота закидываем его данные в редиску
	_ = service.Redis.SetRobotData(strconv.Itoa(createdRobot.ID), createdRobot, 5*time.Minute)
	service.publishToRabbitWithStruct(&createdRobot, keyadd)
	return createdRobot.ID, err
}

func (service *RobotService) GetRobotInfo(id int) (*entities.Robot, error) {
	idStr := strconv.Itoa(id)
	// Пытаемся получить данные из кэша, если они есть - получаем ошибку и идём дальше по коду, если данные есть то ретёрним их
	robotdata, err := service.Redis.GetRobotData(idStr)
	if err == nil {
		service.publishToRabbitWithStruct(robotdata, keyget)
		return robotdata, nil
	}
	// Если данных в кэше нет, то обращаемся к репозиторию и получаем данные из БД
	robotdata, err = service.RobotRepository.GetRobotInfo(id)
	if err != nil {
		return nil, err
	}
	// Добавляем полученные данные в кэш на пять минут
	err = service.Redis.SetRobotData(idStr, *robotdata, 5*time.Minute)
	service.publishToRabbitWithStruct(robotdata, keyget)
	return robotdata, err
}

func (service *RobotService) UpdateRobotCords(updateData dto.UpdateRobotCordDTO) error {
	newCord := entities.RobotCord{XCord: updateData.XCord, YCord: updateData.YCord, ZCord: updateData.ZCord}
	robotID := updateData.ID
	err := service.RobotRepository.UpdateRobotCords(robotID, newCord)
	if err != nil {
		return err
	}
	// Удаление кэша после обновления
	_ = service.Redis.DeleteRobotData(strconv.Itoa(robotID))
	msgToRabbit := fmt.Sprintf("Координаты робота с ID: %d были изменены на  X:%d, Y:%d, Z:%d", robotID, newCord.XCord, newCord.YCord, newCord.ZCord)
	service.publishToRabbitWithText(msgToRabbit, keyupdatecords)
	return nil
}

func (service *RobotService) UpdateRobotName(updateData dto.UpdateRobotNameDTO) error {
	newName := updateData.Name
	robotID := updateData.ID
	err := service.RobotRepository.UpdateRobotName(robotID, newName)
	if err != nil {
		return err
	}
	// Удаление кэша после обновления
	_ = service.Redis.DeleteRobotData(strconv.Itoa(robotID))
	msgToRabbit := fmt.Sprintf("Имя робота с ID: %d было изменены на %s", robotID, newName)
	service.publishToRabbitWithText(msgToRabbit, keyupdatename)
	return nil
}

func (service *RobotService) ChangeRobotType(updateData dto.ChangeTypeDTO) error {
	newType := updateData.Type
	robotID := updateData.ID
	err := service.RobotRepository.ChangeRobotType(robotID, newType)
	if err != nil {
		return err
	}
	// Удаление кэша после обновления
	_ = service.Redis.DeleteRobotData(strconv.Itoa(robotID))
	msgToRabbit := fmt.Sprintf("Тип робота с ID: %d был изменен на %s", robotID, newType)
	service.publishToRabbitWithText(msgToRabbit, keyupdatetype)
	return nil
}

func (service *RobotService) DeleteRobot(id int) error {
	msgToRabbit := fmt.Sprintf("Робота с ID: %d был уничтожен. Помянем...", id)
	service.publishToRabbitWithText(msgToRabbit, keydel)
	_ = service.Redis.DeleteRobotData(strconv.Itoa(id))
	return service.RobotRepository.DeleteRobot(id)
}

// Отправка в реббит сообщения со струтурой робота
func (s *RobotService) publishToRabbitWithStruct(robot *entities.Robot, routingKey string) {
	if err := s.Rabbit.Publish(robot, routingKey); err != nil {
		log.Println("Не получилось отправить, сорян")
	}
}

// Отправка в реббит сообщения с текстом
func (s *RobotService) publishToRabbitWithText(msg string, routingKey string) {
	if err := s.Rabbit.PublishWithText(msg, routingKey); err != nil {
		log.Println("Не получилось отправить, сорян")
	}
}
