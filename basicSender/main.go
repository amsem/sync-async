package main

import (
	"log"
	"context"
    "time"

	amqp "github.com/rabbitmq/amqp091-go"
)




func handleError(err error, msg string)  {
    if err != nil {
        log.Fatalf("%s : %s", msg, err)
    }
}

func main()  {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    handleError(err, "Dialing failed to  RabbitMQ broker")
    defer conn.Close()

    channel, e := conn.Channel()
    handleError(e, "Fetching channel failed ")
    defer channel.Close()
    testQueue, err := channel.QueueDeclare(
        "test",
        false,
        false,
        false,
        false,
        nil,
    )
    handleError(err, "Queue Creation failed")
    serverTime := time.Now()
    message := amqp.Publishing{
        ContentType: "text/plain",
        Body: []byte(serverTime.String()),
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    err = channel.PublishWithContext(
        ctx,
        "",
        testQueue.Name,
        false,
        false,
        message,
    )
    handleError(err, "Failed publishing the msg")
    log.Println("SuccesFully published a message to the queue")
}
