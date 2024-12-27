package models

type BeltModel struct {
	ID         string  `bson:"id"`
	MaterialId int32   `bson:"materialId"`
	Parent     int32   `bson:"parent"`
	MaxRatio   int32   `bson:"maxRatio"`
	MinRatio   int32   `bson:"minRatio"`
	Ratio      int32   `bson:"ratio"`
	SpecVol    float64 `bson:"specVol"`
	//-----------------------------
	MaterialName string  `bson:"materialName,omitempty"`
	MaxWater     float64 `bson:"maxWater,omitempty"`
	MinWater     float64 `bson:"minWater,omitempty"`
	Water        float64 `bson:"water,omitempty"`
}
