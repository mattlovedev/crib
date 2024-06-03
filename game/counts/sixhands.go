package counts

type sixHandChoice struct {
	Hand     uint32
	Crib     uint16
	Avg      float32
	Min      uint8
	Median   uint8
	Max      uint8
	Mode     uint8
	ModeP    float32
	BelowAvg uint8
	AboveAvg uint8
	StdDev   float32
}
