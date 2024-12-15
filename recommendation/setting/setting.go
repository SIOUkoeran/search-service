package setting

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Configuration struct {
	Server   ServerSettingS   `mapstructure:"Server"`
	App      AppSettingS      `mapstructure:"App"`
	Log      LogSettingS      `mapstructure:"Log"`
	Database DatabaseSettingS `mapstructure:"Database"`
}

type ServerSettingS struct {
	RunMode      string        `mapstructure:"RunMode"`
	HTTPPort     string        `mapstructure:"HTTPPort"`
	ReadTimeout  time.Duration `mapstructure:"ReadTimeout"`
	WriteTimeout time.Duration `mapstructure:"WriteTimeout"`
}

type AppSettingS struct {
	ServerShutdownTimeout time.Duration `mapstructure:"ServerShutdownTimeout"`
}

type LogSettingS struct {
	LogSavePath string `mapstructure:"LogSavePath"`
	LogFileName string `mapstructure:"LogFileName"`
	MaxSize     int    `mapstructure:"MaxSize"`
	MaxBackups  int    `mapstructure:"MaxBackups"`
	Compress    bool   `mapstructure:"Compress"`
	Level       string `mapstructure:"Level"`
}

type DatabaseSettingS struct {
	Address  string `mapstructure:"address"`
	IdleHost int    `mapstructure:"idleHost"`
}

func LoadConfig(cfg *Configuration) error {
	viper.AddConfigPath("configs/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return err
	}

	return nil
}
