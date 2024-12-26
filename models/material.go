package models

type MaterialModel struct {
	ID           int32   `bson:"id"`
	MaterialName string  `bson:"materialName"`
	MaxWater     float64 `bson:"maxWater"`
	MinWater     float64 `bson:"minWater"`
	Water        float64 `bson:"water"`
	MaxRatio     int32   `bson:"maxRatio"`
	MinRatio     int32   `bson:"minRatio"`
}
