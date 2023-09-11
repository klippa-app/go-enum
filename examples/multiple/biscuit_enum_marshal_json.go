// Code generated by go-enum, DO NOT EDIT.
package multiple

import (
	"strconv"
	"strings"
)

func (biscuit_enum Biscuit) MarshalJSON() ([]byte, error) {
	err := biscuit_enum.Validate() 
	if err != nil {
		return nil, err
	}

	return []byte(strconv.Quote(biscuit_enum.String())), nil
}

func (biscuit_enum *Biscuit) UnmarshalJSON(val []byte) error {
	str := string(val)
	str = strings.Trim(str, "\"")

	enum, err := BiscuitFromString(str)
	if err != nil {
		return err
	}

	*biscuit_enum = *enum
	return nil
}