package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"reflect"
)


func Parse(spec interface{}) {
	t := reflect.TypeOf(spec)
	if t.Kind() != reflect.Ptr || spec == nil {
		log.Panicf("failed reading environment variable(s) for config type '%s' since it is not a pointer", t.Name())
	}
	err := envconfig.Process("", spec)
	if err != nil {
		log.Panicf("failed reading environment variable(s) for config type '%s'. Error: %s", t.Elem().Name(), err)
	}
}