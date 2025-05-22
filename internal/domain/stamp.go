package domain

type Stamp struct {
	DD   DD
	FRMT string
}

type DD struct {
	RE    string
	TD    uint8
	F     int64
	FE    string
	RR    string
	RSR   string
	MNT   uint64
	IT1   string
	CAF   string
	TSTED string
}
