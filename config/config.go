package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	LogLevel string `yaml:"logLevel"`
	Db       struct {
		Driver string `yaml:"driver"`
		Mysql  string `yaml:"mysql"`
	}
	Mode string `yaml:"mode"`
	Auth struct {
		TokenExpire bool `yaml:"tokenExpire"`
	}
	BareMetal struct {
		SadTalker struct {
			PythonPath  string   `yaml:"pythonPath"`
			ProjectPath string   `yaml:"projectPath"`
			ExtraArgs   []string `yaml:"extraArgs"`
		}
		Roop struct {
			PythonPath  string   `yaml:"pythonPath"`
			ProjectPath string   `yaml:"projectPath"`
			ExtraArgs   []string `yaml:"extraArgs"`
		}
		OpenVoice struct {
			PythonPath  string   `yaml:"pythonPath"`
			ProjectPath string   `yaml:"projectPath"`
			ExtraArgs   []string `yaml:"extraArgs"`
		}
	}
	CloudNative struct {
		SadTalker struct {
			JobNamespace string `yaml:"jobNamespace"`
		}
	}
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return
}
