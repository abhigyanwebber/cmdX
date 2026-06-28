package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// GlobalConfig holds user preferences that persist across sessions
type GlobalConfig struct {
	DefaultTheme string `mapstructure:"default_theme"`
	RenderMode   string `mapstructure:"render_mode"`
	ChafaPath    string `mapstructure:"chafa_path"`
	AssetsDir    string `mapstructure:"assets_dir"`
	ThemesDir    string `mapstructure:"themes_dir"`
	PluginsDir   string `mapstructure:"plugins_dir"`
	AutoInject   bool   `mapstructure:"auto_inject"`
	ShowBanner   bool   `mapstructure:"show_banner"`
}

// DefaultGlobalConfig returns sensible defaults
func DefaultGlobalConfig() GlobalConfig {
	return GlobalConfig{
		DefaultTheme: "default",
		RenderMode:   "auto",
		ChafaPath:    "chafa",
		AssetsDir:    "",
		ThemesDir:    "",
		PluginsDir:   "",
		AutoInject:   false,
		ShowBanner:   true,
	}
}

// ConfigDir returns the cmdX config directory (~/.cmdx)
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	return filepath.Join(home, ".cmdx"), nil
}

// InitGlobalConfig sets up viper with defaults and config file
func InitGlobalConfig() error {
	configDir, err := ConfigDir()
	if err != nil {
		return err
	}

	// set config file location
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	// environment variable overrides
	viper.SetEnvPrefix("CMDX")
	viper.AutomaticEnv()

	// defaults
	defaults := DefaultGlobalConfig()
	viper.SetDefault("default_theme", defaults.DefaultTheme)
	viper.SetDefault("render_mode", defaults.RenderMode)
	viper.SetDefault("chafa_path", defaults.ChafaPath)
	viper.SetDefault("assets_dir", defaults.AssetsDir)
	viper.SetDefault("themes_dir", defaults.ThemesDir)
	viper.SetDefault("plugins_dir", defaults.PluginsDir)
	viper.SetDefault("auto_inject", defaults.AutoInject)
	viper.SetDefault("show_banner", defaults.ShowBanner)

	// read config file — ignore if not found
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config: %w", err)
		}
	}

	return nil
}

// GetGlobalConfig returns the current global config
func GetGlobalConfig() (GlobalConfig, error) {
	var cfg GlobalConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return DefaultGlobalConfig(), fmt.Errorf("could not parse config: %w", err)
	}
	return cfg, nil
}

// SaveGlobalConfig writes the current config to ~/.cmdx/config.yaml
func SaveGlobalConfig(cfg GlobalConfig) error {
	configDir, err := ConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	viper.Set("default_theme", cfg.DefaultTheme)
	viper.Set("render_mode", cfg.RenderMode)
	viper.Set("chafa_path", cfg.ChafaPath)
	viper.Set("assets_dir", cfg.AssetsDir)
	viper.Set("themes_dir", cfg.ThemesDir)
	viper.Set("plugins_dir", cfg.PluginsDir)
	viper.Set("auto_inject", cfg.AutoInject)
	viper.Set("show_banner", cfg.ShowBanner)

	configPath := filepath.Join(configDir, "config.yaml")
	return viper.WriteConfigAs(configPath)
}
