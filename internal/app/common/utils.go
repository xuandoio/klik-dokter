package common

import (
	"encoding/json"
	"github.com/go-pg/pg/v10/types"
	"strings"
	"time"
)

type (
	Date     struct{ time.Time }
	DateTime struct{ time.Time }
)

// UnmarshalJSON /**
func (d *Date) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.TrimSpace(strings.Trim(strInput, `"`))
	if len(strInput) == 0 {
		return nil
	}
	newTime, err := time.Parse("2006-01-02", strInput)
	if err != nil {
		return err
	}

	*d = Date{newTime}
	return nil
}

// MarshalJSON /**
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.FormatISO())
}

// FormatISO /**
func (d Date) FormatISO() string {
	return d.Format("2006-01-02")
}

// Scan /**
func (d *Date) Scan(b interface{}) error {
	if b == nil {
		d.Time = time.Time{}
		return nil
	}
	newD, err := types.ParseTime(b.([]byte))
	if err != nil {
		return err
	}
	d.Time = newD
	return nil
}

// UnmarshalJSON /**
func (d *DateTime) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	if len(strInput) == 0 {
		return nil
	}
	strInput = strings.Trim(strInput, `"`)
	newTime, err := time.Parse("2006-01-02 15:04:05", strInput)

	if err != nil {
		return err
	}

	*d = DateTime{newTime}
	return nil
}

// MarshalJSON /**
func (d DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.FormatISO())
}

// FormatISO /**
func (d DateTime) FormatISO() string {
	return d.Format("2006-01-02 15:04:05")
}

// Scan /**
func (d *DateTime) Scan(b interface{}) error {
	if b == nil {
		d.Time = time.Time{}
		return nil
	}
	newD, err := types.ParseTime(b.([]byte))
	if err != nil {
		return err
	}
	d.Time = newD
	return nil
}
