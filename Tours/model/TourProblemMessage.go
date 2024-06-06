package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type TourProblemMessage struct {
	ID           int64     `json:"id"`
	SenderId     int64     `json:"senderId"`
	RecipientId  int64     `json:"recipientId"`
	CreationTime time.Time `json:"creationTime"`
	Description  string    `json:"description"`
	IsRead       bool      `json:"isRead"`
}

func (d TourProblemMessage) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *TourProblemMessage) Scan(value interface{}) error {
	if value == nil {
		*d = TourProblemMessage{}
		return nil
	}

	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("Scan source is not []byte")
	}

	return json.Unmarshal(bytes, d)
}

type TourProblemMessages []TourProblemMessage

func (d TourProblemMessages) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *TourProblemMessages) Scan(value interface{}) error {
	if value == nil {
		*d = TourProblemMessages{}
		return nil
	}

	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("Scan source is not []byte")
	}

	return json.Unmarshal(bytes, d)
}
