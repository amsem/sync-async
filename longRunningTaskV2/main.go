package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	redis "github.com/redis/go-redis/v9"
)


const queueName string = "jobQueue"
const hostString string = "127.0.0.1:8000"

func handleError(err error, msg string)  {
    if err != nil {
        log.Fatalf("%s : %s", msg, err)
    }
}

func getServer(name string) JobServer {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    handleError(err, "Failed to connect to RabbitMQ")

    channel, err := conn.Channel()
    handleError(err, "Fetching channel Failed")

    jobQueue, err := channel.QueueDeclare(
        name,
        false,
        false,
        false,
        false,
        nil,
    )
    handleError(err, "Job queue creation Failed")
    return JobServer{Conn: conn, Channel: channel, Queue: jobQueue}
}

func main()  {
    jobS := getServer(queueName)
    jobS.redisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
    })
    go func(conn *amqp.Connection) {
        workerProcess := Workers{
        conn: jobS.Conn,
    }
    workerProcess.run()
    }(jobS.Conn)
    router := mux.NewRouter()
    router.HandleFunc("/job/db", jobS.asyncDBHandler)
    router.HandleFunc("/job/callback", jobS.asyncCallBackHandler)
    router.HandleFunc("/job/mail", jobS.asyncMailHandler)
    router.HandleFunc("/job/status", jobS.statusHandler)

    Server := &http.Server{
        Handler: router,
        Addr: hostString,
        ReadTimeout: 15 * time.Second,
        WriteTimeout: 15 * time.Second,
    }
    log.Fatal(Server.ListenAndServe())
    defer jobS.Channel.Close()
    defer jobS.Conn.Close()
}


