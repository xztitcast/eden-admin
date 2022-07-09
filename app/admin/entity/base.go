package entity

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

type BaseEntity struct {
	Created JSONTime `gorm:"comment:创建时间" json:"created" form:"created"`
	//Updated time.Time `gorm:"comment:更新时间" json:"updated"`
}

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
)

type JSONTime time.Time

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.Parse(DefaultTimeFormat, string(data))
	if err != nil {
		return err
	}
	*t = JSONTime(now)
	return nil
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DefaultTimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, DefaultTimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t JSONTime) String() string {
	return time.Time(t).Format(DefaultTimeFormat)
}

type Status bool

func (s *Status) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	status, err := strconv.ParseBool(string(data))
	fmt.Println(string(data))
	fmt.Println(status)
	if err != nil {
		return err
	}
	*s = Status(status)
	return nil
}

func (s Status) MarshalJSON() ([]byte, error) {
	var u int32
	if s {
		u = 1
	} else {
		u = 0
	}
	buff := bytes.NewBuffer([]byte{})
	binary.Write(buff, binary.BigEndian, u)
	return buff.Bytes(), nil
}
