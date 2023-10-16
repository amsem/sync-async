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
    testQueue, err := channel.QueueDeclare(
        "test",
        false,
        false,
        false,
        false,
        nil,
    )
    handleError(err, "Queue Creation failed")
    message, err := channel.Consume(
        testQueue.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    handleError(err, "Failed to register a consumer")
    go func() {
        for m := range message {
            log.Printf("Received a message from the queue : %s", m.Body)
        }
    }()
    log.Println("Worker has started")
    wait := make(chan bool)
    <-wait
}
