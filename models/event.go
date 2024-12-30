package models

import "time"

type EventModel struct {
	Shift string    `bson:"shift"`
	Date  time.Time `bson:"date"`
	Event string    `bson:"event"`
}
