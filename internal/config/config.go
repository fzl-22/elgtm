package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	SCM    SCM    `mapstructure:"scm"`
	LLM    LLM    `mapstructure:"llm"`
	Review Review `mapstructure:"review"`
	System System `mapstructure:"system"`
}

type SCMPlatform string

const (
	PlatformGitHub SCMPlatform = "github"
	PlatformGitLab SCMPlatform = "gitlab"
)

type SCM struct {
	Platform    SCMPlatform `mapstructure:"platform"`
	Token       string      `mapstructure:"token"`
	Owner       string      `mapstructure:"owner"`
	Repo        string      `mapstructure:"repo"`
	PRNumber    int         `mapstructure:"pr_number"`
	MaxDiffSize int64       `mapstructure:"max_diff_size"`
}

type LLMProvider string

const ProviderGemini LLMProvider = "gemini"

type LLM struct {
	Provider    LLMProvider `mapstructure:"provider"`
	Model       string      `mapstructure:"model"`
	APIKey      string      `mapstructure:"api_key"`
	Temperature float32     `mapstructure:"temperature"`
	MaxTokens   int         `mapstructure:"max_tokens"`
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

func NewConfig() (*Config, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("scm.max_diff_size", 2097152)

	v.SetDefault("llm.temperature", 0.2)
	v.SetDefault("llm.max_tokens", 4096)

	v.SetDefault("review.prompt_type", "general")
	v.SetDefault("review.prompt_dir", ".reviewer")
	v.SetDefault("review.language", "en")

	v.SetDefault("system.log_level", "info")
	v.SetDefault("system.timeout", 300)

	cfg := &Config{}

	BindEnvs(v, cfg)

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return cfg, nil
}

func BindEnvs(v *viper.Viper, iface any, parts ...string) {
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
			BindEnvs(v, fieldV.Interface(), append(parts, tv)...)
		} else {
			v.BindEnv(keyPath)
		}
	}
}
