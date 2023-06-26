package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)
import amqp "github.com/rabbitmq/amqp091-go"

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func getEnvVarString(key string) string {
	err := godotenv.Load(".env")
	failOnError(err, "Could not load .env file.")
	return os.Getenv(key)
}

func initializeServer(rmqConnection *amqp.Connection) {
	var logServer LogServer
	var err error

	logServer.amqpChannel, err = rmqConnection.Channel()
	failOnError(err, "Failed to open a channel")

	defer logServer.amqpChannel.Close()

	err = logServer.amqpChannel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	logServer.clients = map[string]string{}

	return
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
