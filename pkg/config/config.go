package pkgconfig

import (
	"github.com/go-playground/validator"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func LoadConfig(prefix string, v interface{}) error {
	if err := envconfig.Process(prefix, v); err != nil {
		return errors.Wrap(err, "failed to load configuration")
	}
	if err := validator.New().Struct(v); err != nil {
		return errors.Wrap(err, "failed to validate configuration")
	}
	return nil
}
