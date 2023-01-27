package middleware

import (
	"fmt"

	"github.com/yhmain/5th-simple-tiktok/config"
	"github.com/yhmain/5th-simple-tiktok/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var myDB *gorm.DB //全局变量，数据库连接

func init() {
	fmt.Println("conn.go ... init")
	params := config.GetConfig().MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		params.Username, params.Password, params.Host, params.Port, params.Database)
	var err error
	myDB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		PrepareStmt: true, //执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
	})
	if err != nil {
		fmt.Println("conn.go: 数据库连接失败", err)
		return
	}
	fmt.Println("conn.go: 数据库连接成功")
	// 创建表的同时进行表属性配置  数据库迁移，只增不减
	myDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.User{}, &model.Video{}, &model.Like{})
}

func GetMySQLClient() *gorm.DB {
	return myDB
}
