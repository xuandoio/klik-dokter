package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string `env:"APP_ENV"`
	Server struct {
		Host string `env:"APP_HOST"`
		Port string `env:"APP_PORT"`
	}
	Database struct {
		Host     string `env:"DB_HOST"`
		Port     string `env:"DB_PORT"`
		Name     string `env:"DB_NAME"`
		Username string `env:"DB_USERNAME"`
		Password string `env:"DB_PASSWORD"`
	}
	JWTSecret string `env:"JWT_SECRET"`
}

// ReadENVFile /**
func ReadENVFile(c *Config, configPath string) error {
	if len(configPath) > 0 {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}

	Decode(c)
	return nil
}

// ValidateConfigPath /**
func ValidateConfigPath(configPath string) error {
	s, err := os.Stat(configPath)
	if err != nil {
		return err
	}

	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", configPath)
	}

	return nil
}

// ParseConfigFlags /**
func ParseConfigFlags() (configPath string, err error) {

	flag.StringVar(&configPath, "config", ".env", "Path to the config file")

	flag.Parse()
	// parse the flags
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// return the config path
	return configPath, nil
}

// Decode /**
func Decode(i interface{}) {
	var t reflect.Type
	var v reflect.Value
	var ok bool

	if v, ok = i.(reflect.Value); !ok {
		v = reflect.ValueOf(i).Elem()
	}

	t = v.Type()

	for i := 0; i < v.NumField(); i++ {
		if v.CanInterface() {
			if v.Field(i).Kind() != reflect.Struct {
				if eTag, ok := t.Field(i).Tag.Lookup("env"); ok {
					var eRVal interface{}
					eString, exist := os.LookupEnv(eTag) // get string value from environment
					// try to cast to real value
					if t.Field(i).Type.String() == "int" {
						if v, err := strconv.Atoi(eString); err == nil {
							eRVal = v
						} else {
							if exist {
								panic(eTag + " must be int type")
							} else {
								eRVal = 0
							}
						}
					} else {
						eRVal = eString
					}

					eVal := reflect.ValueOf(eRVal)
					if eVal.IsValid() {
						v.Field(i).Set(eVal.Convert(t.Field(i).Type))
					}
				}
			} else if v.Field(i).Kind() == reflect.Struct {
				Decode(v.Field(i))
			} else {
				panic("can not reflect value:" + v.Field(i).Kind().String())
			}
		}
	}
}
