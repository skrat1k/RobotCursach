package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Robot struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	XCord int    `json:"xCord"`
	YCord int    `json:"yCord"`
	ZCord int    `json:"zCord"`
}

const (
	rabbitURL    = "amqp://guest:guest@localhost:5672/"
	exchangeName = "robots"
)

var queueBindings = []struct {
	QueueName  string
	RoutingKey string
}{
	{"robot_add_queue", "robots.Add"},
	{"robot_get_queue", "robots.Get"},
	{"robot_updatecord_queue", "robots.UpdateCord"},
	{"robot_updatename_queue", "robots.UpdateName"},
	{"robot_delete_queue", "robots.Del"},
}

func main() {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal("Cannot connect to rabbit", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Cannot open channel", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error create exchange", err)
	}

	for _, binding := range queueBindings {

		q, err := ch.QueueDeclare(
			binding.QueueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Error create queue %s: %v", binding.QueueName, err)
		}

		err = ch.QueueBind(
			q.Name,
			binding.RoutingKey,
			exchangeName,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Error binding queue %s: %v", binding.QueueName, err)
		}

		msg, err := ch.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Error consuming from queue %s: %v", binding.QueueName, err)
		}

		log.Println("NotificationService listening queue", q.Name)

		go func(queueName string, messages <-chan amqp.Delivery) {
			for d := range messages {
				var robot Robot
				var msg string
				switch d.ContentType {
				case "application/json":
					if err := json.Unmarshal(d.Body, &robot); err != nil {
						log.Printf("[%s] JSON unmarshal error: %v", queueName, err)
						continue
					}
				case "text/plain":
					msg = string(d.Body)
				}
				switch queueName {
				case "robot_add_queue":
					log.Printf("[%s] ADD Robot: ID=%d, Name=%s", queueName, robot.ID, robot.Name)
				case "robot_get_queue":
					log.Printf("[%s] GET Robot : ID:%d, Nname:%s. Coordinates: X=%d, Y=%d, Z=%d", queueName, robot.ID, robot.Name, robot.XCord, robot.YCord, robot.ZCord)
				case "robot_updatecord_queue":
					log.Printf("[%s] %s", queueName, msg)
				case "robot_updatename_queue":
					log.Printf("[%s] %s", queueName, msg)
				case "robot_delete_queue":
					log.Printf("[%s] %s", queueName, msg)
				default:
					log.Printf("[%s] RobotInfo: ID=%d, Name=%s, XCord=%d, YCord=%d, ZCord=%d", queueName, robot.ID, robot.Name, robot.XCord, robot.YCord, robot.ZCord)
				}
			}
		}(binding.QueueName, msg)
	}
	forever := make(chan bool)
	<-forever
}

func AddQueue(robot Robot, queueName string) {
	log.Printf("[%s] ADD Robot: ID=%d, Name=%s", queueName, robot.ID, robot.Name)
}

func GetQueue(robot Robot, queueName string) {
	log.Printf("[%s] GET Robot : ID:%d, Nname:%s. Coordinates: X=%d, Y=%d, Z=%d", queueName, robot.ID, robot.Name, robot.XCord, robot.YCord, robot.ZCord)
}

func UpdateCordQueue(msg string, queueName string) {
	log.Printf("[%s] %s", queueName, msg)
}

func UpdateNameQueue(msg string, queueName string) {
	log.Printf("[%s] %s", queueName, msg)
}

func DeleteQueue(msg string, queueName string) {
	log.Printf("[%s] %s", queueName, msg)
}
