package config

import (
	"encoding/json"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
)

var (
	AppConfig        *AppConfiguration
	AppConfigDynamic *AppConfigurationDynamic
	appConfigMutex   sync.Mutex
)

func init() {
	mylog.Log.Info("init begin: config")

	var err error
	// 加载App配置数据
	if AppConfig, err = LoadAppConfig(); err != nil {
		mylog.Log.Error("LoadAppConfig err=", err)
	}

	// 加载APP动态配置数据
	if AppConfigDynamic, err = LoadAppConfigDynamic(); err != nil {
		mylog.Log.Error("LoadAppConfigDynamic err=", err)
	}

	mylog.Log.Info("init end  : config")
}

// App配置数据
type AppConfiguration struct {
}

// App配置数据(可动态更新)
type AppConfigurationDynamic struct {
	IPWhiteList     []string `json:"ipWhiteList"`
	IPBlackList     []string `json:"ipBlackList"`
	IPAdminList     []string `json:"ipAdminList"`
	RateLimitConfig struct {
		RequestsPerSecond float64 `json:"requests_per_second"`
		BucketSize        int     `json:"bucket_size"`
	} `json:"rateLimitConfig"`
	CodeSandboxType string `json:"codeSandboxType"` //docker、native、dockerAndNative
}

// 加载App配置数据
func LoadAppConfig() (*AppConfiguration, error) {
	filePath := "conf/appconfig.json"
	config := &AppConfiguration{}

	// 判断配置文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		mylog.Log.Warn("App配置文件不存在")
		return nil, err
	}

	// 打开项目配置文件
	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	// 解码配置文件内容到结构体
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(config); err != nil {
		return nil, err
	}
	mylog.Log.Info("App配置加载成功")
	return config, nil
}

// 加载App配置数据(可动态更新)
func LoadAppConfigDynamic() (*AppConfigurationDynamic, error) {
	filePath := "conf/appdynamicconfig.json"
	config := &AppConfigurationDynamic{}

	// 判断配置文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		mylog.Log.Warn("App动态配置文件不存在")
		return nil, err
	}

	// 打开项目配置文件
	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	// 解码配置文件内容到结构体
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(config); err != nil {
		return nil, err
	}
	mylog.Log.Info("App动态配置加载成功")
	return config, nil
}

// 定时任务，轮询更新动态配置
func LoadAppDynamicConfigCycle() {
	filePath := "conf/appdynamicconfig.json"
	ticker := time.NewTicker(2 * time.Second) // 每3秒检查一次配置文件
	defer ticker.Stop()

	var lastModTime time.Time
	var lastConfig *AppConfigurationDynamic // 保存配置数据

	for range ticker.C {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			mylog.Log.Errorf("Error reading config file: %v", err)
			continue
		}

		if fileInfo.ModTime() != lastModTime {
			lastModTime = fileInfo.ModTime()

			newConfig, err := LoadAppConfigDynamic()
			if err != nil {
				mylog.Log.Errorf("Error loading config: %v", err)
				// todo 更新加载App配置数据失败，需报警
				continue
			}

			// 检查新配置与旧配置是否相同，避免不必要的重新加载
			appConfigMutex.Lock()
			if !reflect.DeepEqual(lastConfig, newConfig) {
				lastConfig = newConfig
				// 在这里使用最新的配置数据进行处理
				mylog.Log.Infof("Loaded new config: %+v", newConfig)
				AppConfigDynamic = newConfig
			}
			appConfigMutex.Unlock()
		}
	}
}
