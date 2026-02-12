package config

import (
	"log"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Config struct {
	AI     AI     `mapstructure:"ai"`
	Git    Git    `mapstructure:"git"`
	Review Review `mapstructure:"review"`
	System System `mapstructure:"system"`
}

type AI struct {
	Provider    string  `mapstructure:"provider"`
	Model       string  `mapstructure:"model"`
	APIKey      string  `mapstructure:"api_key"`
	Temperature float64 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
}

type Git struct {
	Platform      string `mapstructure:"platform"`
	Token         string `mapstructure:"token"`
	RepoOwner     string `mapstructure:"repo_owner"`
	RepoName      string `mapstructure:"repo_name"`
	PullRequestID int    `mapstructure:"pr_id"`
}

type Review struct {
	PromptType string `mapstructure:"prompt_type"`
	PromptDir  string `mapstructure:"prompt_dir"`
	Language   string `mapstructure:"language"`
}

type System struct {
	LogLevel string `mapstructure:"log_level"`
	Timeout  int    `mapstructure:"timeout"`
}

func New() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	bindEnvs(Config{})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables.")
		} else {
			log.Printf("Error reading config file: %v. Error type: %T\n", err, err)
			return nil, err
		}
	}

	var config Config
	err := viper.Unmarshal(&config, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			StringToSliceHookFunc(),
		),
	))
	if err != nil {
		log.Fatalf("Unable to decode config: %v", err)
		return nil, err
	}

	return &config, nil
}

func StringToSliceHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Slice {
			return data, nil
		}

		raw := data.(string)
		if raw == "" {
			return []string{}, nil
		}

		parts := strings.Split(raw, ",")
		result := make([]string, 0, len(parts))

		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				result = append(result, trimmed)
			}
		}

		return result, nil
	}
}

func bindEnvs(iface any, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)

	if ift.Kind() == reflect.Pointer {
		ifv = ifv.Elem()
		ift = ift.Elem()
	}

	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}

		keyPath := strings.Join(append(parts, tv), ".")

		if v.Kind() == reflect.Struct {
			bindEnvs(v.Interface(), append(parts, tv)...)
		} else {
			viper.BindEnv(keyPath)
		}
	}
}
