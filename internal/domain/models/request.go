package models

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	ID      uuid.UUID
	Name    string
	Phone   string
	Email   string
	CarInfo string
	Date    time.Time
}
