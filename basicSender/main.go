package main

import (
    "log"
   amqp"github.com/rabbitmq/amqp091-go"
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

}
