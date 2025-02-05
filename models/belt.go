package models

type BeltModel struct {
	ID          string  `bson:"id"`
	MaterialId  int32   `bson:"materialId"`
	Parent      int32   `bson:"parent"`
	MaxRatio    int32   `bson:"maxRatio"`
	MinRatio    int32   `bson:"minRatio"`
	Ratio       int32   `bson:"ratio"`
	SpecVol     float64 `bson:"specVol"`
	Vol         float64 `bson:"vol"`
	Diff        float64 `bson:"diff"`
	Rate        float64 `bson:"rate"`
	WetAcc      float64 `bson:"wetAcc"`
	DryAcc      float64 `bson:"dryAcc"`
	Running     int32   `bson:"running"`
	Seconds     int32   `bson:"seconds"`
	Numerator   float64 `bson:"numerator"`   //rate 分子
	Denominator float64 `bson:"denominator"` //rate 分母
	//-----------------------------
	MaterialName string  `bson:"materialName,omitempty"`
	MaxWater     float64 `bson:"maxWater,omitempty"`
	MinWater     float64 `bson:"minWater,omitempty"`
	Water        float64 `bson:"water,omitempty"`
}

// BeltModel 定义 belt 表的模型
type BeltModel_PV_SP struct {
	ID      string  `bson:"id"`
	Vol     float64 `bson:"vol"`
	SpecVol float64 `bson:"specVol"`
}
type BeltModel_Other struct {
	ID           string  `bson:"id"`
	Parent       int32   `bson:"parent"`
	Ratio        int32   `bson:"ratio"`
	SpecVol      float64 `bson:"specVol"`
	Vol          float64 `bson:"vol"`
	Diff         float64 `bson:"diff"`
	Rate         float64 `bson:"rate"`
	WetAcc       float64 `bson:"wetAcc"`
	DryAcc       float64 `bson:"dryAcc"`
	MaterialName string  `bson:"materialName,omitempty"`
	Water        float64 `bson:"water,omitempty"`
}
