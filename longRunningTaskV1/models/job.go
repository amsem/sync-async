package models

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
    ID uuid.UUID `json:"uuid"`
    Type string `json:"type"`
    ExtraData interface{} `json:"extra_data"`
}

type Log struct {
    ClientTime time.Time `json:"client_time"`
}

type CallBack struct {
    CallBackURL string `json:"callback_url"`
}

type Mail struct {
    EmailAdress string `json:"email_adress"`
}
