package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	Port          string          `mapstructure:"PORT"`
	Environment   string          `mapstructure:"ENVIRONMENT"`
	AppName       string          `mapstructure:"APP_NAME"`
	DebugMode     bool            `mapstructure:"DEBUG"`
	PostgresDSN   string          `mapstructure:"DB_POSTGRES_DSN"`
	JWTAccessKey  *rsa.PrivateKey `mapstructure:"JWT_ACCESS_KEY"`
	JWTRefreshKey *rsa.PrivateKey `mapstructure:"JWT_REFRESH_KEY"`
}

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()

		viper.SetTypeByDefaultValue(true)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("error reading config file: %s", err)
		}

		decCfg := &mapstructure.DecoderConfig{
			Result:  &config,
			TagName: "mapstructure",
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				stringToBoolHook,
				stringToPrivateKeyHook,
			),
		}
		decoder, err := mapstructure.NewDecoder(decCfg)
		if err != nil {
			log.Fatalf("error creating decoder: %s", err)
		}

		if err := decoder.Decode(viper.AllSettings()); err != nil {
			log.Fatalf("error decoding config: %s", err)
		}
	})
	return config
}

func stringToBoolHook(
	from reflect.Type, to reflect.Type, data any,
) (any, error) {
	if from.Kind() == reflect.String && to.Kind() == reflect.Bool {
		switch strings.ToLower(data.(string)) {
		case "true", "1", "yes", "on":
			return true, nil
		case "false", "0", "no", "off":
			return false, nil
		default:
			return nil, fmt.Errorf("invalid boolean value: %s", data)
		}
	}
	return data, nil
}

func stringToPrivateKeyHook(
	from reflect.Type, to reflect.Type, data any,
) (any, error) {
	if from.Kind() == reflect.String && to == reflect.TypeOf((*rsa.PrivateKey)(nil)) {
		return loadPrivateKeyFromPEM(data.(string))
	}
	return data, nil
}

func loadPrivateKeyFromPEM(pemStr string) (*rsa.PrivateKey, error) {
	pemStr = strings.ReplaceAll(pemStr, "\\n", "\n")

	decoded, err := base64.StdEncoding.DecodeString(pemStr)
	if err != nil {
		decoded = []byte(pemStr)
	}

	block, _ := pem.Decode(decoded)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM: no block found")
	}

	keyIfc, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8 key: %w", err)
	}
	rsaKey, ok := keyIfc.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parsed key is not RSA")
	}
	return rsaKey, nil
}
