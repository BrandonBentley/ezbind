# ezbind
 
ezbind is a simple module for creating environment bindings for [spf13/viper](https://github.com/spf13/viper) struct configurations.


## Installation
``` bash
go get github.com/BrandonBentley/ezbind
```

## Usage
Usage is simply calling the `ezbind.BindStruct` function with the [viper](https://github.com/spf13/viper) and instance of the destination struct and it will create the bindings converting the key paths into enviroment variables by replacing all `.` with `_` and uppercasing all of it.

Example: `some.key.path` -> `SOME_KEY_PATH`

### Excluding Fields
To exclude fields from binding to environment variables, add `envbind:"false"` to the struct field tags that you want to exclude. If the `envbind` tag is not present or anything besides `false`, the behavior is defaulted to `true`

Example Struct:
``` go 
type MyStruct struct {
    IWantThisOne string `mapstructure:"i_want_this_one"`
    NoEnv string `mapstructure:"i_want_this_one" envbind:"false"` // no env binding
    NoneOfTheseFields struct {
        Field1 string `mapstructure:"i_want_this_one"` // no env binding
    }`mapstructure:"none_of_these_fields" envbind:"false"`
}
```


## Full Example w/ Output
``` go
package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/BrandonBentley/ezbind"
	"github.com/spf13/viper"
)

// Emulating a file set to be read
var jsonConfigFile = strings.NewReader(`{"server":{"http":{"port":80}}}`)

func init() {
	// Assume these are set before running
	os.Setenv("AUTH_TOKEN", "ABCDEFG")
}

type Config struct {
	Server Server `mapstructure:"server"`
	Auth   Auth   `mapstructure:"auth"`
}

type Server struct {
	Http Http `mapstructure:"http"`
}

type Http struct {
	Port int `mapstructure:"port"`
}

type Auth struct {
	Token string `mapstructure:"token"`
}

func main() {
	cfg := viper.New()
	cfg.SetConfigType("json")

	var config Config

	ezbind.BindStruct(cfg, config)

	cfg.ReadConfig(jsonConfigFile)

	cfg.Unmarshal(&config)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	enc.Encode(config)
}
```

#### Output 
``` bash
> go run example/example.go
{
  "Server": {
    "Http": {
      "Port": 80
    }
  },
  "Auth": {
    "Token": "ABCDEFG"
  }
}
```