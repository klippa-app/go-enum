// Code generated by go-enum, DO NOT EDIT.
package singlefile

import (
	"database/sql/driver"
	"fmt"
)

func (biscuit_enum Biscuit) Value() (driver.Value, error) {
	return biscuit_enum.String(), biscuit_enum.Validate()
}

func (biscuit_enum *Biscuit) Scan(val any) error {
	var str string

	switch v := val.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("unsupported type %T", v)
	}

	enum, err := BiscuitFromString(str)
	if err != nil {
		return err
	}

	*biscuit_enum = *enum
	return nil
}