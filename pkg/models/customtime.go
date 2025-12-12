package models

import (
	"strings"
	"time"
)

type CustomTime struct {
    time.Time
}

const customLayout = "2006-01-02 15:04:05" // The format of your input string

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
    s := strings.Trim(string(b), "\"")
    if s == "null" {
        ct.Time = time.Time{}
        return nil
    }
    t, err := time.Parse(customLayout, s)
    if err != nil {
        return err
    }
    ct.Time = t
    return nil
}

