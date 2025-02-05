package models

type AlertModel struct {
	ID          string  `bson:"id"`
	DiffAlert   int32   `bson:"diffAlert"`
	DiffSeconds int32   `bson:"diffSeconds"`
	Diff        float64 `bson:"diff"`
}

type SendAlertModel struct {
	ID        string `bson:"id"`
	DiffAlert int32  `bson:"diffAlert"`
}
