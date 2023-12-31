// Code generated by go-enum, DO NOT EDIT.
package multiple

import (
	"github.com/globalsign/mgo/bson"
)

func (cookie_enum Cookie) GetBSON() (interface{}, error) {
	err := cookie_enum.Validate() 
	if err != nil {
		return nil, err
	}

	return cookie_enum.String(), nil
}

func (cookie_enum *Cookie) SetBSON(raw bson.Raw) error {
	var str string
	
	err := raw.Unmarshal(&str)
	if err != nil {
		return err
	}

	enum, err := CookieFromString(str)
	if err != nil {
		return err
	}

	*cookie_enum = *enum
	return nil
}
