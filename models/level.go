package models

// LevelModel level表模型
type LevelModel struct {
	ID             int32   `bson:"id"`
	MaterialLevel1 float64 `bson:"MaterialLevel1"`
	MaterialLevel2 float64 `bson:"MaterialLevel2"`
}
