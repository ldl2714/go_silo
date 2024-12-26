package models

import "time"

type StaticModel struct {
	ID           int32     `bson:"id"`
	SetFlowRate1 int32     `bson:"setFlowRate1"`
	SetFlowRate2 int32     `bson:"setFlowRate2"`
	Shift        string    `bson:"shift"`
	Team         string    `bson:"team"`
	Time         time.Time `bson:"time"`
}
