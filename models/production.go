package models

type ProductionModel struct {
	ID           string  `bson:"id"`
	MaterialName string  `bson:"materialName"`
	Water        float64 `bson:"water"`
	Ratio        int32   `bson:"ratio"`
	SpecVol      float64 `bson:"specVol"`
	Vol          float64 `bson:"vol"`
	Diff         float64 `bson:"diff"`
	Rate         float64 `bson:"rate"`
	WetAcc       float64 `bson:"wetAcc"`
	DryAcc       float64 `bson:"dryAcc"`
	Running      int32   `bson:"running"`
	Shift        string  `bson:"shift"`
	Team         string  `bson:"team"`
	Date         string  `bson:"date"`
	Parent       int32   `bson:"parent"`
}
