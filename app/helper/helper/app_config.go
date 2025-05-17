package helper

import (
	"gin_base/app/appconfig"
	"github.com/spf13/viper"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

var (
	appConfig     appconfig.Config
	appConfigOnce sync.Once
)

// 初始化应用各种配置
func initAppConfig() {
	appConfigOnce.Do(func() {
		v := viper.New()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.AutomaticEnv()

		configPath := "app/appconfig"
		filepath.Walk(configPath, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() && filepath.Ext(path) == ".yaml" {
				subV := viper.New()
				subV.SetConfigFile(path)
				subV.SetConfigType("yaml")
				subV.ReadInConfig()
				v.MergeConfigMap(subV.AllSettings())
			}
			return nil
		})

		v.Unmarshal(&appConfig)
	})
}

// 获取应用的各种配置
func GetAppConfig() appconfig.Config {
	initAppConfig()
	return appConfig
}
