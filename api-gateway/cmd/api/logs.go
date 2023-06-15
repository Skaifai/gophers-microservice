package main

//import (
//	"context"
//	"encoding/json"
//	amqp "github.com/rabbitmq/amqp091-go"
//	"io"
//	"os"
//	"runtime/debug"
//	"sync"
//	"time"
//)
//
//type Level int8
//
//const (
//	LevelInfo Level = iota
//	LevelError
//	LevelFatal
//	LevelOff
//)
//
//func (l Level) String() string {
//	switch l {
//	case LevelInfo:
//		return "INFO"
//	case LevelError:
//		return "ERROR"
//	case LevelFatal:
//		return "FATAL"
//	default:
//		return ""
//	}
//}
//
//type loggerService struct {
//	out           io.Writer
//	minLevel      Level
//	mu            sync.Mutex
//	rmqConnection *amqp.Connection
//}
//
//func (l *loggerService) New(out io.Writer, minLevel Level, rmqConn *amqp.Connection) *loggerService {
//	return &loggerService{
//		out:           out,
//		minLevel:      minLevel,
//		rmqConnection: rmqConn,
//	}
//}
//
//func (l *loggerService) SendInfo(message string, properties map[string]string) {
//	err := l.printLog(LevelInfo, message, properties)
//	if err != nil {
//
//	}
//}
//func (l *loggerService) SendError(err error, properties map[string]string) {
//	l.printLog(LevelError, err.Error(), properties)
//}
//func (l *loggerService) SendFatal(err error, properties map[string]string) {
//	l.printLog(LevelFatal, err.Error(), properties)
//	os.Exit(1)
//}
//
//func (ls *loggerService) sendLog(message string) error {
//	conn, err := amqp.Dial(rmqDSN)
//	failOnError(err, "Failed to connect to RabbitMQ")
//	defer conn.Close()
//
//	ch, err := conn.Channel()
//	failOnError(err, "Failed to open a channel")
//	defer ch.Close()
//
//	err = ch.ExchangeDeclare(
//		"logs",
//		"fanout",
//		true,
//		false,
//		false,
//		false,
//		nil)
//	if err != nil {
//		return err
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	body := message
//	err = ch.PublishWithContext(ctx,
//		"logs", // exchange
//		"",     // routing key
//		false,  // mandatory
//		false,  // immediate
//		amqp.Publishing{
//			ContentType: "text/plain",
//			Body:        []byte(body),
//		})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (l *loggerService) printLog(level Level, message string, properties map[string]string) error {
//	if level < l.minLevel {
//		return nil
//	}
//	aux := struct {
//		Level      string            `json:"level"`
//		Time       string            `json:"time"`
//		Message    string            `json:"message"`
//		Properties map[string]string `json:"properties,omitempty"`
//		Trace      string            `json:"trace,omitempty"`
//	}{
//		Level:      level.String(),
//		Time:       time.Now().UTC().Format(time.RFC3339),
//		Message:    message,
//		Properties: properties,
//	}
//	if level >= LevelError {
//		aux.Trace = string(debug.Stack())
//	}
//	var line []byte
//	line, err := json.Marshal(aux)
//	if err != nil {
//		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
//	}
//
//	l.mu.Lock()
//	defer l.mu.Unlock()
//
//	err = l.sendLog(string(line))
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
