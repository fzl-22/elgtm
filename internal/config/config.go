package config

import (
	"fmt"
	"reflect"
	"strings"

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
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	bindEnvs(v, Config{})

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}

func bindEnvs(v *viper.Viper, iface any, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)

	if ift.Kind() == reflect.Pointer {
		ifv = ifv.Elem()
		ift = ift.Elem()
	}

	for i := 0; i < ift.NumField(); i++ {
		fieldV := ifv.Field(i)
		fieldT := ift.Field(i)
		tv, ok := fieldT.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}

		keyPath := strings.Join(append(parts, tv), ".")

		if fieldV.Kind() == reflect.Struct {
			bindEnvs(v, fieldV.Interface(), append(parts, tv)...)
		} else {
			v.BindEnv(keyPath)
		}
	}
}
