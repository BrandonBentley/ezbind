package ezbind

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func BindStruct(cfg *viper.Viper, v any) {
	bindStructWithPrefix(cfg, "", v)
}

func bindStructWithPrefix(cfg *viper.Viper, prefix string, v any) {
	t := reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		mapstructure := field.Tag.Get("mapstructure")
		envbind := field.Tag.Get("envbind")
		if envbind == "false" {
			continue
		}
		if mapstructure == "" {
			continue
		}

		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct {
			bindStructWithPrefix(cfg, join(prefix, mapstructure), reflect.New(fieldType).Elem().Interface())
		} else {
			cfg.BindEnv(joinKeyToEnv(cfg, prefix, mapstructure))
		}
	}
}

func join(strs ...string) string {
	if len(strs) > 0 && strs[0] == "" {
		strs = strs[1:]
	}
	return strings.Join(strs, ".")
}

func joinKeyToEnv(cfg *viper.Viper, keys ...string) (string, string) {
	key := join(keys...)

	envVar := strings.ReplaceAll(key, ".", "_")
	if cfg.GetEnvPrefix() != "" {
		envVar = cfg.GetEnvPrefix() + "_" + envVar
	}

	return key, strings.ToUpper(envVar)
}
