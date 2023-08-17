package options

import (
	"github.com/klippa-app/go-enum/internal/util"
)

type Option string

const (
	DefaultOption Option = "default"
	InvalidOption Option = "invalid"
)

var validOptions = []Option{
	DefaultOption,
	InvalidOption,
}

func (o Option) isValid() bool {
	return util.Contains(validOptions, o)
}
