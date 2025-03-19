package ezbind

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestBindStruct(t *testing.T) {
	os.Clearenv()
	testConfig := struct {
		F1 string `mapstructure:"f1" envbind:"true"`
		F2 struct {
			F3 *struct {
				F4 int `mapstructure:"f4"`
			} `mapstructure:"f3"`
		} `mapstructure:"f2"`
		F5 struct {
			F6 float64 `mapstructure:"f6"`
		} `mapstructure:"f5"`
		F7 struct {
			F8 bool `mapstructure:"f7" envbind:"true"`
		} `mapstructure:"f7" envbind:"false"`
	}{}

	os.Setenv("F1", "string1")
	os.Setenv("F2_F3_F4", "1")
	os.Setenv("F5_F6", "3.14")
	os.Setenv("F7_F8", "true")

	cfg := viper.New()
	cfg.SetConfigType("json")

	BindStruct(cfg, &testConfig)

	assert.NoError(t, cfg.ReadConfig(strings.NewReader("{}")))

	assert.NoError(t, cfg.Unmarshal(&testConfig))

	assert.Equal(t, "string1", testConfig.F1)
	assert.Equal(t, 1, testConfig.F2.F3.F4)
	assert.Equal(t, float64(3.14), testConfig.F5.F6)
	assert.False(t, testConfig.F7.F8)
}

func TestBindStructWithPrefix_NoMatch(t *testing.T) {
	os.Clearenv()
	testConfig := struct {
		F1 string `mapstructure:"f1" envbind:"true"`
		F2 struct {
			F3 struct {
				F4 int `mapstructure:"f4"`
			} `mapstructure:"f3"`
		} `mapstructure:"f2"`
		F5 struct {
			F6 float64 `mapstructure:"f6"`
		} `mapstructure:"f5"`
		F7 struct {
			F8 bool `mapstructure:"f7" envbind:"true"`
		} `mapstructure:"f7" envbind:"false"`
	}{}

	os.Setenv("F1", "string1")
	os.Setenv("F2_F3_F4", "1")
	os.Setenv("F5_F6", "3.14")
	os.Setenv("F7_F8", "true")

	cfg := viper.New()
	cfg.SetConfigType("json")

	bindStructWithPrefix(cfg, "NotFound", &testConfig)

	assert.NoError(t, cfg.ReadConfig(strings.NewReader("{}")))

	assert.NoError(t, cfg.Unmarshal(&testConfig))

	assert.NotEqual(t, "string1", testConfig.F1)
	assert.NotEqual(t, 1, testConfig.F2.F3.F4)
	assert.NotEqual(t, float64(3.14), testConfig.F5.F6)
	assert.False(t, testConfig.F7.F8)
}

func TestBindStructWithPrefix_Match(t *testing.T) {
	os.Clearenv()
	testConfig := struct {
		F1 string `mapstructure:"f1" envbind:"true"`
		F2 struct {
			F3 struct {
				F4 int `mapstructure:"f4"`
			} `mapstructure:"f3"`
		} `mapstructure:"f2"`
		F5 struct {
			F6 float64 `mapstructure:"f6"`
		} `mapstructure:"f5"`
		F7 struct {
			F8 string `mapstructure:"f8" envbind:"true"`
		} `mapstructure:"f7" envbind:"true"`
	}{}

	prefix := "SomePrefix"

	upperPrefix := strings.ToUpper(prefix) + "_"
	os.Setenv("F1", "string1")
	os.Setenv(upperPrefix+"F2_F3_F4", "1")
	os.Setenv(upperPrefix+"F5_F6", "3.14")

	cfg := viper.New()
	cfg.SetEnvPrefix(prefix)

	cfg.SetConfigType("json")

	BindStruct(cfg, testConfig)

	expectedJsonValue := "somethingReallyCool"
	assert.NoError(t, cfg.ReadConfig(strings.NewReader(`{"f7": {"f8": "`+expectedJsonValue+`"} }`)))

	assert.NoError(t, cfg.Unmarshal(&testConfig))

	assert.NotEqual(t, "string1", testConfig.F1)
	assert.Equal(t, 1, testConfig.F2.F3.F4)
	assert.Equal(t, float64(3.14), testConfig.F5.F6)
	assert.Equal(t, expectedJsonValue, testConfig.F7.F8)
}

func TestBindStructWithPrefix_NonMapStucture(t *testing.T) {
	os.Clearenv()
	testConfig := struct {
		F1 string `mapstructure:"f1" envbind:"true"`
		F2 struct {
			F3 struct {
				F4 int `mapstructure:"f4"`
			} `mapstructure:"f3"`
		} `mapstructure:"f2"`
		F5 struct {
			F6 float64 `mapstructure:"f6"`
		}
		F7 struct {
			F8 string `mapstructure:"f8" envbind:"true"`
		} `mapstructure:"f7" envbind:"true"`
	}{}

	os.Setenv("F1", "string1")
	os.Setenv("F2_F3_F4", "1")
	os.Setenv("F5_F6", "3.14")

	cfg := viper.New()

	cfg.SetConfigType("json")

	BindStruct(cfg, testConfig)

	expectedJsonValue := "somethingReallyCool"
	assert.NoError(t, cfg.ReadConfig(strings.NewReader(`{"f7": {"f8": "`+expectedJsonValue+`"} }`)))

	assert.NoError(t, cfg.Unmarshal(&testConfig))

	assert.Equal(t, "string1", testConfig.F1)
	assert.Equal(t, 1, testConfig.F2.F3.F4)
	assert.NotEqual(t, float64(3.14), testConfig.F5.F6)
	assert.Equal(t, expectedJsonValue, testConfig.F7.F8)
}

func TestBindStructWithPrefix_GracefulHandlingNonStruct(t *testing.T) {
	os.Clearenv()
	testConfig := struct {
		F1 string `mapstructure:"f1" envbind:"true"`
		F2 struct {
			F3 struct {
				F4 int `mapstructure:"f4"`
			} `mapstructure:"f3"`
		} `mapstructure:"f2"`
		F5 struct {
			F6 float64 `mapstructure:"f6"`
		}
		F7 struct {
			F8 string `mapstructure:"f8" envbind:"true"`
		} `mapstructure:"f7" envbind:"true"`
	}{}

	os.Setenv("F1", "string1")
	os.Setenv("F2_F3_F4", "1")
	os.Setenv("F5_F6", "3.14")

	cfg := viper.New()

	cfg.SetConfigType("json")

	BindStruct(cfg, "NotAStruct")

	expectedJsonValue := "somethingReallyCool"
	assert.NoError(t, cfg.ReadConfig(strings.NewReader(`{"f7": {"f8": "`+expectedJsonValue+`"} }`)))

	assert.NoError(t, cfg.Unmarshal(&testConfig))
}
