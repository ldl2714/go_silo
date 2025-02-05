package models

type TrendModel struct {
	ID      string `bson:"id"`
	SpecVol string `bson:"specVol"`
	Vol     string `bson:"vol"`
	Date    string `bson:"date"`
}
