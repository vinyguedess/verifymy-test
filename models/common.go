package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

const (
	DateFormat = "2006-01-02"
)

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	stringedTime := strings.Trim(string(b), "\"")
	parsedTime, err := time.Parse(DateFormat, stringedTime)
	if err != nil {
		return err
	}

	*d = Date(parsedTime)
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	timeObj := time.Time(*d)
	return []byte(fmt.Sprintf(`"%s"`, timeObj.Format(DateFormat))), nil
}

func (t *Date) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return t.UnmarshalText(string(v))
	case string:
		return t.UnmarshalText(v)
	case time.Time:
		*t = Date(v)
	case nil:
		*t = Date{}
	default:
		return fmt.Errorf("cannot sql.Scan() MyTime from: %#v", v)
	}
	return nil
}

func (t Date) Value() (driver.Value, error) {
	return driver.Value(time.Time(t).Format(DateFormat)), nil
}

func (t *Date) UnmarshalText(value string) error {
	dd, err := time.Parse(DateFormat, value)
	if err != nil {
		return err
	}
	*t = Date(dd)
	return nil
}

func (Date) GormDataType() string {
	return "DATE"
}

type SecretValue string

func (sv *SecretValue) MarshalJSON() ([]byte, error) {
	return []byte(`null`), nil
}

func (sv *SecretValue) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return sv.UnmarshalText(string(v))
	case string:
		return sv.UnmarshalText(v)
	case nil:
		*sv = ""
	default:
		return fmt.Errorf("cannot sql.Scan() MyTime from: %#v", v)
	}
	return nil
}

func (sv *SecretValue) Value() (driver.Value, error) {
	return driver.Value(sv), nil
}

func (sv *SecretValue) UnmarshalText(value string) error {
	*sv = SecretValue(value)
	return nil
}

func (*SecretValue) GormDataType() string {
	return "VARCHAR"
}
