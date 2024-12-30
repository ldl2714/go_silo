package models

type PidModel struct {
	ID         string `bson:"id"`
	PID_MAN    bool   `bson:"pid_man"`
	PID_HALT   bool   `bson:"pid_halt"`
	PID_D_ON_X bool   `bson:"pid_d_on_x"`
	PID_D      bool   `bson:"pid_d"`
	PID_I      bool   `bson:"pid_i"`
	PID_P      bool   `bson:"pid_p"`
	PID_QMAX   bool   `bson:"pid_qmax"`
	PID_QMIN   bool   `bson:"pid_qmin"`

	PID_SP   float32 `bson:"pid_sp"`
	PID_PV   float32 `bson:"pid_pv"`
	PID_BIAS float32 `bson:"pid_bias"`
	PID_ERR  float32 `bson:"pid_err"`
	PID_GAIN float32 `bson:"pid_gain"`
	PID_Y    float32 `bson:"pid_y"`
	PID_YMAX float32 `bson:"pid_ymax"`
	PID_YMIN float32 `bson:"pid_ymin"`
	PID_YMAN float32 `bson:"pid_yman"`

	PID_TD     int32 `bson:"pid_td"`
	PID_TD_LAG int32 `bson:"pid_td_lag"`
	PID_TI     int32 `bson:"pid_ti"`
}
