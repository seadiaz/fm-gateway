package datatypes

import (
	"encoding/xml"
	"fmt"
	"time"
)

// Date es un tipo personalizado para parsear fechas en formato YYYY-MM-DD.
type Date struct {
	time.Time
}

func (d *Date) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return fmt.Errorf("error parsing date: %w", err)
	}
	d.Time = t
	return nil
}
