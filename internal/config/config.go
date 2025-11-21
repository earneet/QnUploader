package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	// 七牛云配置
	QiniuAccessKey string `mapstructure:"qiniu_access_key"`
	QiniuSecretKey string `mapstructure:"qiniu_secret_key"`
	QiniuBucket    string `mapstructure:"qiniu_bucket"`
	QiniuDomain    string `mapstructure:"qiniu_domain"`

	// 快捷键配置
	HotkeyKeys  []int `mapstructure:"hotkey_keys"`
	HotkeyCtrl  bool  `mapstructure:"hotkey_ctrl"`
	HotkeyShift bool  `mapstructure:"hotkey_shift"`
	HotkeyAlt   bool  `mapstructure:"hotkey_alt"`

	// UI配置
	AutoCopyURL  bool `mapstructure:"auto_copy_url"`
	ShowProgress bool `mapstructure:"show_progress"`
}

func Load() (*Config, error) {
	// 加载.env文件（如果存在）
	_ = godotenv.Load()

	// 设置配置文件路径
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，使用环境变量和默认值
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// 绑定环境变量
	bindEnvVars()

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// getConfigDir 获取配置目录
func getConfigDir() (string, error) {
	// 优先使用用户配置目录
	if configDir := os.Getenv("QINIU_UPLOADER_CONFIG_DIR"); configDir != "" {
		return configDir, nil
	}

	// 使用系统默认配置目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "qu")

	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

// setDefaults 设置默认配置
func setDefaults() {
	// 七牛云配置默认值
	viper.SetDefault("qiniu_access_key", "")
	viper.SetDefault("qiniu_secret_key", "")
	viper.SetDefault("qiniu_bucket", "")
	viper.SetDefault("qiniu_domain", "")

	// 快捷键配置默认值 (Ctrl+Shift+U)
	viper.SetDefault("hotkey_keys", []int{85}) // U键
	viper.SetDefault("hotkey_ctrl", true)
	viper.SetDefault("hotkey_shift", true)
	viper.SetDefault("hotkey_alt", false)

	// UI配置默认值
	viper.SetDefault("auto_copy_url", true)
	viper.SetDefault("show_progress", true)
}

// bindEnvVars 绑定环境变量
func bindEnvVars() {
	viper.BindEnv("qiniu_access_key", "QINIU_ACCESS_KEY")
	viper.BindEnv("qiniu_secret_key", "QINIU_SECRET_KEY")
	viper.BindEnv("qiniu_bucket", "QINIU_BUCKET")
	viper.BindEnv("qiniu_domain", "QINIU_DOMAIN")
}

// Save 保存配置
func Save(cfg *Config) error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	// 设置配置值
	viper.Set("qiniu_access_key", cfg.QiniuAccessKey)
	viper.Set("qiniu_secret_key", cfg.QiniuSecretKey)
	viper.Set("qiniu_bucket", cfg.QiniuBucket)
	viper.Set("qiniu_domain", cfg.QiniuDomain)
	viper.Set("hotkey_keys", cfg.HotkeyKeys)
	viper.Set("hotkey_ctrl", cfg.HotkeyCtrl)
	viper.Set("hotkey_shift", cfg.HotkeyShift)
	viper.Set("hotkey_alt", cfg.HotkeyAlt)
	viper.Set("auto_copy_url", cfg.AutoCopyURL)
	viper.Set("show_progress", cfg.ShowProgress)

	// 保存到文件
	configFile := filepath.Join(configDir, "config.yaml")
	return viper.WriteConfigAs(configFile)
}