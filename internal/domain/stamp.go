package domain

type Stamp struct {
	DD   DD
	FRMT string
}

type DD struct {
	RE    string   `xml:"RE"`
	TD    uint8    `xml:"TD"`
	F     int64    `xml:"F"`
	FE    string   `xml:"FE"`
	RR    string   `xml:"RR"`
	RSR   string   `xml:"RSR"`
	MNT   uint64   `xml:"MNT"`
	IT1   string   `xml:"IT1"`
	CAF   StampCAF `xml:"CAF"`
	TSTED string   `xml:"TSTED"`
}

// StampCAF represents the CAF structure specifically for stamp XML serialization
type StampCAF struct {
	Version string    `xml:"version,attr"`
	DA      StampDA   `xml:"DA"`
	FRMA    StampFRMA `xml:"FRMA"`
}

type StampDA struct {
	RE    string     `xml:"RE"`
	RS    string     `xml:"RS"`
	TD    uint8      `xml:"TD"`
	RNG   StampRNG   `xml:"RNG"`
	FA    string     `xml:"FA"`
	RSAPK StampRSAPK `xml:"RSAPK"`
	IDK   string     `xml:"IDK"`
}

type StampRNG struct {
	D int64 `xml:"D"`
	H int64 `xml:"H"`
}

type StampRSAPK struct {
	M string `xml:"M"`
	E string `xml:"E"`
}

type StampFRMA struct {
	Algorithm string `xml:"algoritmo,attr"`
	Value     string `xml:",chardata"`
}
