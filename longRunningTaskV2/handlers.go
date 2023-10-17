package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

    redis "github.com/redis/go-redis/v9"
	"github.com/amsem/longRunningTaskV2/models"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type JobServer struct {
    Queue amqp.Queue
    Channel *amqp.Channel
    Conn *amqp.Connection
    redisClient *redis.Client
}


func (s *JobServer) Publish(jsonBody []byte) error  {
    message := amqp.Publishing{
        ContentType: "application/json",
        Body: jsonBody,
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
    err := s.Channel.PublishWithContext(
        ctx,
        "",
        queueName,
        false,
        false,
        message,
    )
    handleError(err, "Error while generating JobID")
    return err
}

func (s *JobServer) statusHandler(w http.ResponseWriter, r *http.Request)  {
    queryParams := r.URL.Query()
    uuid := queryParams.Get("uuid")
    w.Header().Set("Content-Type", "application/json")
    ctx := context.Background()
    jobStatus := s.redisClient.Get(ctx, uuid)
    status := map[string]string{"uuid": uuid, "status": jobStatus.Val()}
    response, err := json.Marshal(status)
    handleError(err,"Cannot create response for client")
    w.Write(response)
}
func (s *JobServer) asyncDBHandler(w http.ResponseWriter, r *http.Request)  {
    jobID, err := uuid.NewRandom()
    queryParams := r.URL.Query()
    unixTime, err := strconv.ParseInt(queryParams.Get("client_time"), 10, 64)
    clientTime := time.Unix(unixTime, 0)
    handleError(err, "Error while converting client time")

    jsonBody, err := json.Marshal(models.Job{ID: jobID, 
        Type:   "A",
        ExtraData: models.Log{ClientTime: clientTime},
    }) 
    handleError(err, "JSON body creation failed")

    if s.Publish(jsonBody) == nil {
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonBody)
    }else {
        w.WriteHeader(http.StatusInternalServerError)
    }
}
func (s *JobServer) asyncMailHandler(w http.ResponseWriter, r *http.Request)  {
    jobID, err := uuid.NewRandom()
    email := "aminesaidsel@gmail.com"
    handleError(err, "Error while converting client time")
    jsonBody, err := json.Marshal(models.Job{ID: jobID,
        Type:"B",
        ExtraData: models.Mail{EmailAdress: email}, 
    })
     handleError(err, "JSON body creation failed")

    if s.Publish(jsonBody) == nil {
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonBody)
    }else {
        w.WriteHeader(http.StatusInternalServerError)
    }
    
}

func (s *JobServer) asyncCallBackHandler(w http.ResponseWriter, r *http.Request)  {
    jobID, err := uuid.NewRandom()
    handleError(err, "Error while converting client time")
    jsonBody, err := json.Marshal(models.Job{ID: jobID,
        Type:"C",
        ExtraData: models.CallBack{CallBackURL:"amsem.com" },
    })
    handleError(err, "JSON body creation failed")

    if s.Publish(jsonBody) == nil {
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonBody)
    }else {
        w.WriteHeader(http.StatusInternalServerError)
    }
    
}

