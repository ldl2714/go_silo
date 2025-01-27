package models

type BeltModel struct {
	ID         string  `bson:"id"`
	MaterialId int32   `bson:"materialId"`
	Parent     int32   `bson:"parent"`
	MaxRatio   int32   `bson:"maxRatio"`
	MinRatio   int32   `bson:"minRatio"`
	Ratio      int32   `bson:"ratio"`
	SpecVol    float64 `bson:"specVol"`
	Vol        float64 `bson:"vol"`
	Diff       float64 `bson:"diff"`
	Rate       float64 `bson:"rate"`
	WetAcc     float64 `bson:"wetAcc"`
	DryAcc     float64 `bson:"dryAcc"`
	Running    int32   `bson:"running"`
	//-----------------------------
	MaterialName string  `bson:"materialName,omitempty"`
	MaxWater     float64 `bson:"maxWater,omitempty"`
	MinWater     float64 `bson:"minWater,omitempty"`
	Water        float64 `bson:"water,omitempty"`
}
