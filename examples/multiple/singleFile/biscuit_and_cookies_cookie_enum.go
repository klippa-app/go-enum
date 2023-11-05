// Code generated by go-enum, DO NOT EDIT.
package singlefile

import (
	"fmt"
)

func AllCookies() []Cookie {
	return []Cookie{
		ChocolateDigestive,
		ChocolateShortbread,
		ChocolateFinger,
		JaffaCake,
		ChocolateHobnob,
	}
}

func validCookies() []Cookie {
	return []Cookie{
		ChocolateDigestive,
		ChocolateShortbread,
		ChocolateFinger,
		JaffaCake,
		ChocolateHobnob,
	}
}

func ToCookie(value int) Cookie {
	cookie_enum := Cookie(value)
	switch cookie_enum {
	case ChocolateDigestive, ChocolateShortbread, ChocolateFinger, JaffaCake, ChocolateHobnob:
		return cookie_enum
	default:
		return ChocolateDigestive
	}
}

func (cookie_enum Cookie) String() string {
	switch cookie_enum {
	case ChocolateDigestive:
		return "ChocolateDigestive"
	case ChocolateShortbread:
		return "ChocolateShortbread"
	case ChocolateFinger:
		return "ChocolateFinger"
	case JaffaCake:
		return "JaffaCake"
	case ChocolateHobnob:
		return "ChocolateHobnob"
	default:
		return ChocolateDigestive.String()
	}
}

func CookieFromString(val string) (*Cookie, error) {
	valid := validCookies()	
	for i := range valid {
		if valid[i].String() == val {
			return &valid[i], nil
		}
	}	

	return nil, fmt.Errorf("%s is not a valid Cookie", val)
}

func (cookie_enum Cookie) Validate() error {
	_, err := CookieFromString(cookie_enum.String())
	return err
}