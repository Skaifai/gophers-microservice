package main

import (
	"flag"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"logger-service/internal/util"
	"os"
	"strconv"
)

type config struct {
	rmqPort     int
	rmqUsername string
	rmqPassword string
}

type LogServer struct {
	clients     map[string]string
	amqpChannel *amqp.Channel
}

const fileNameLayout = "02-01-2006-15H-04M-05S"
const logHeaderLayout = "02.01.2006 15:04:05"

var rmqDSN string

func main() {
	var cfg config
	rabbitMQPort, err := strconv.Atoi(getEnvVarString("PORT"))
	failOnError(err, "Could not convert string to int")
	flag.IntVar(&cfg.rmqPort, "rabbitMQPort", rabbitMQPort, "The message broker port")
	flag.StringVar(&cfg.rmqUsername, "rabbitMQUsername", getEnvVarString("RMQ_USERNAME"), "The message broker username")
	flag.StringVar(&cfg.rmqPassword, "rabbitMQPassword", getEnvVarString("RMQ_PASSWORD"), "The message broker password")

	flag.Parse()

	rmqDSN = fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.rmqUsername, cfg.rmqPassword, "rabbitmq", cfg.rmqPort)
	conn, err := amqp.Dial(rmqDSN)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	initializeServer(conn)

	var f *os.File

	var input string

	for input != "exit" {
		fmt.Print("Enter command: ")
		input = util.ReadAndCleanString()
		switch input {
		case "create":
			if f != nil {
				f.Close()
			}
			f = CreateLogFile()
			fmt.Println("A new log file has been created at " + f.Name())
		case "dump":
			if f != nil {
				ReceiveMessages(f)
			} else {
				fmt.Println("Assign the file first.")
			}
		case "exit":
			break
		default:
			fmt.Println("Incorrect input.")
		}
	}

	defer f.Close()
}
