package model

import (
	"encoding/json"
	"io"
)

type TransportationType int
type TourDurations []TourDuration

type TourDuration struct {
	TimeInSeconds  uint `json:"TimeInSeconds" bson:"timeInSeconds"`
	Transportation int  `json:"Transportation" bson:"transportation"`
}

const (
	Walking TransportationType = iota
	Bicycle
	Car
)

func (p *TourDurations) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
func (p *TourDuration) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

// func (d TourDuration) Value() (driver.Value, error) {
// 	return json.Marshal(d)
// }

// func (d *TourDuration) Scan(value interface{}) error {
// 	if value == nil {
// 		*d = TourDuration{}
// 		return nil
// 	}

// 	bytes, ok := value.([]byte)

// 	if !ok {
// 		return fmt.Errorf("Scan source is not []byte")
// 	}

// 	return json.Unmarshal(bytes, d)
// }

//

// func (td TourDurations) Value() (driver.Value, error) {
// 	return json.Marshal(td)
// }

// func (td *TourDurations) Scan(value interface{}) error {
// 	if value == nil {
// 		*td = TourDurations{}
// 		return nil
// 	}

// 	bytes, ok := value.([]byte)

// 	if !ok {
// 		return fmt.Errorf("Scan source is not []byte")
// 	}

// 	return json.Unmarshal(bytes, td)
// }
