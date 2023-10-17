package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/amsem/longRunningTaskV2/models"
	amqp "github.com/rabbitmq/amqp091-go"
	redis "github.com/redis/go-redis/v9"
)

type Workers struct {
    conn *amqp.Connection
    redisClient redis.Client
}

func (w *Workers) dbWork(job models.Job)  {
    ctx := context.Background()
    result := job.ExtraData.(map[string]interface{})
    w.redisClient.Set(ctx, job.ID.String(), "STARTED", 0)
    log.Printf("Worker %s: exracting data ..., JOB: %s \n", job.Type, result)
    w.redisClient.Set(ctx, job.ID.String(), "IN PROGRESS", 0)
    time.Sleep(1 * time.Minute)
    log.Printf("Worker %s: saving data to database .... , JOB: %s \n",job.Type, job.ID)
    w.redisClient.Set(ctx, job.ID.String(), "DONE", 0)
}

func (w *Workers) emailWork(job models.Job)  {
    log.Printf("Worker %s: sending email ..., JOB %s \n", job.Type, job.ID)
    time.Sleep(2 * time.Second)
    log.Printf("Worker %s: sent email successfully, JOB : %s \n", job.Type, job.ID)
}

func (w *Workers) callbackWork(job models.Job)  {
    log.Printf("Worker %s: performing some long runniung process ...., JOB : %s \n", job.Type, job.ID)
    time.Sleep(10 * time.Second)
    log.Printf("Worker %s: posting the data back to the given callback ..., JOB: %s \n", job.Type, job.ID)
}

func (w *Workers) run()  {
    log.Println("Workers are booted up and running")
    channel, err := w.conn.Channel()
    handleError(err, "Failed to fetch channel")
    defer channel.Close()
    w.redisClient = *redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
    })
    jobQueue, err := channel.QueueDeclare(
        queueName,
        false,
        false,
        false,
        false,
        nil,
    )
    handleError(err, "Job queue fetch failed")

    messages, err := channel.Consume(
        jobQueue.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    handleError(err, "failed to consume message")
    go func() {
        for message := range messages {
            job := models.Job{}
            errc := json.Unmarshal(message.Body, &job)
            log.Printf("Workers received a message from the queue : %s", &job)
            handleError(errc, "Unable to load queue message")
            switch job.Type {
            case "A":
                w.dbWork(job)
            case "B":
                w.callbackWork(job)
            case "C":
                w.emailWork(job)
            }
        }
    }()
    defer w.conn.Close()
    wait := make(chan bool)
    <-wait
}
