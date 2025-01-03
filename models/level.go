package models

// LevelModel level表模型
type LevelModel struct {
	ID             int32   `bson:"id"`
	MaterialLevel1 float64 `bson:"MaterialLevel1"`
	MaterialLevel2 float64 `bson:"MaterialLevel2"`
	Disk1          bool    `bson:"disk1"`
	Disk2          bool    `bson:"disk2"`
	Disk3          bool    `bson:"disk3"`
	Disk4          bool    `bson:"disk4"`
}
