package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// 参考网址：https://www.jianshu.com/p/84499381a7da
// yaml.v2只能读取yaml格式的文件，而Viper可以处理多种格式的配置
// 对应yarm文件中的字段

type Config struct {
	App       *App       `yaml:"app"`
	MySQL     *Mysql     `yaml:"mysql"`
	Snowflake *SnowFlake `yaml:"snowflake"`
}

type App struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	Video string `yaml:"video"`
	Img   string `yaml:"img"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type SnowFlake struct { // 雪花算法
	MachineID int64 `yaml:"machineID"`
}

var ProjectConfig *Config // 声明成全局变量

// 使用Viper解析yarm文件
func InitConfigByViper() {
	viper.SetConfigType("yaml")   // 配置文件的格式
	viper.AddConfigPath(".")      // 第一个搜索路径
	viper.AddConfigPath("../")    // 第二个搜索路径
	viper.SetConfigName("config") // 配置文件的名称（可无后缀）
	// viper.SetConfigFile("../config.yaml")        // 配置文件的路径
	if err := viper.ReadInConfig(); err != nil { // 先读取配置文件，看是否正常
		fmt.Println(err.Error())
	}
	if err := viper.Unmarshal(&ProjectConfig); err != nil { // 解析到结构体
		fmt.Println(err.Error())
	}
	// 查看一下解析得到的内容
	// fmt.Printf("config: %#v\n", ProjectConfig)
}

// 即使被import多次，init函数只执行一次
func init() {
	fmt.Println("config.go ... init")
	InitConfigByViper()
	fmt.Println("config.go ... 初始化成功！")
}
