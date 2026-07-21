package tanswer

import (
	"os"
	"strconv"
	"time"

	"github.com/chaitin/chaitin-cli/config"
	"github.com/spf13/cobra"
)

const defaultTimeout = 30 * time.Second

type RuntimeConfig struct {
	URL      string        `yaml:"url"`
	APIKey   string        `yaml:"api_key"`
	Timeout  time.Duration `yaml:"timeout"`
	Insecure bool          `yaml:"insecure"`
	Output   string        `yaml:"output"`
}

func DefaultRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		Timeout: defaultTimeout,
		Output:  "json",
	}
}

func ApplyRuntimeConfig(cmd *cobra.Command, cfg config.Raw) {
	runtime := DefaultRuntimeConfig()
	if productCfg, err := config.DecodeProduct[RuntimeConfig](cfg, "tanswer"); err == nil {
		runtime = mergeRuntimeConfig(runtime, productCfg)
	}

	runtime = mergeEnvConfig(runtime)
	applyFlagString(cmd, "url", runtime.URL)
	applyFlagString(cmd, "api-key", runtime.APIKey)
	applyFlagDuration(cmd, "timeout", durationOrDefault(runtime.Timeout, defaultTimeout))
	applyFlagBool(cmd, "insecure", runtime.Insecure)
	applyFlagString(cmd, "output", runtime.Output)
}

func mergeRuntimeConfig(base, next RuntimeConfig) RuntimeConfig {
	if next.URL != "" {
		base.URL = next.URL
	}
	if next.APIKey != "" {
		base.APIKey = next.APIKey
	}
	if next.Timeout > 0 {
		base.Timeout = next.Timeout
	}
	if next.Insecure {
		base.Insecure = true
	}
	if next.Output != "" {
		base.Output = next.Output
	}
	return base
}

func mergeEnvConfig(cfg RuntimeConfig) RuntimeConfig {
	if v := firstEnv("TANSWER_URL", "TA_ANSWER_ADDR"); v != "" {
		cfg.URL = v
	}
	if v := firstEnv("TANSWER_API_KEY", "TA_ANSWER_TOKEN"); v != "" {
		cfg.APIKey = v
	}
	if v := firstEnv("TA_ANSWER_TIMEOUT"); v != "" {
		if timeout, err := time.ParseDuration(v); err == nil {
			cfg.Timeout = timeout
		}
	}
	if v := firstEnv("TA_ANSWER_INSECURE_SKIP_VERIFY"); v != "" {
		if insecure, err := strconv.ParseBool(v); err == nil {
			cfg.Insecure = insecure
		}
	}
	return cfg
}

func firstEnv(names ...string) string {
	for _, name := range names {
		if value := os.Getenv(name); value != "" {
			return value
		}
	}
	return ""
}

func applyFlagString(cmd *cobra.Command, name, value string) {
	flag := cmd.Flags().Lookup(name)
	if flag == nil || flag.Changed || value == "" {
		return
	}
	_ = cmd.Flags().Set(name, value)
}

func applyFlagDuration(cmd *cobra.Command, name string, value time.Duration) {
	flag := cmd.Flags().Lookup(name)
	if flag == nil || flag.Changed || value <= 0 {
		return
	}
	_ = cmd.Flags().Set(name, value.String())
}

func applyFlagBool(cmd *cobra.Command, name string, value bool) {
	flag := cmd.Flags().Lookup(name)
	if flag == nil || flag.Changed || !value {
		return
	}
	_ = cmd.Flags().Set(name, strconv.FormatBool(value))
}
