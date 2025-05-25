package db_helper

import (
	"fmt"
	"gin_base/app/helper/exception_helper"
	"gin_base/app/helper/helper"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"path/filepath"
	"sync"
)

var dbsRWMutex sync.RWMutex
var dbs = make(map[string]*gorm.DB)

// 初始化数据库
func initDb(connectName string) *gorm.DB {
	dbsRWMutex.RLock() //读锁
	db, exists := dbs[connectName]
	dbsRWMutex.RUnlock() //立即解开读锁
	if exists {
		return db
	}
	//不存在，需要初始化
	dbsRWMutex.Lock()         //写锁
	defer dbsRWMutex.Unlock() //初始化完毕自动解开写锁
	//防止在获取写锁已经被其他协程初始化
	if db, exists := dbs[connectName]; exists {
		return db
	}

	var dialector gorm.Dialector
	switch helper.GetAppConfig().Database[connectName].Driver {
	case "mysql": //暂时只支持mysql
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			helper.GetAppConfig().Database[connectName].Username,
			helper.GetAppConfig().Database[connectName].Password,
			helper.GetAppConfig().Database[connectName].Host,
			helper.GetAppConfig().Database[connectName].Port,
			helper.GetAppConfig().Database[connectName].Name,
		) // 连接数据库
		dialector = mysql.Open(dsn)
	case "sqlite":
		// SQLite数据库路径
		dbPath := helper.GetAppConfig().Database[connectName].Name
		if !filepath.IsAbs(dbPath) {
			// 如果不是绝对路径，则相对于应用程序目录
			appDir, _ := os.Getwd()
			dbPath = filepath.Join(appDir, dbPath)
		}

		// 确保数据库文件所在目录存在
		dbDir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			exception_helper.CommonException(fmt.Sprintf("无法创建SQLite数据库目录: %v", err))
		}

		dialector = sqlite.Open(dbPath)
	default:
		exception_helper.CommonException(fmt.Sprintf("不支持的数据库驱动: %s", helper.GetAppConfig().Database[connectName].Driver))
	}

	db, _ = gorm.Open(dialector, &gorm.Config{
		Logger: &DbLogger{logger.Default},
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   helper.GetAppConfig().Database[connectName].Table_Prefix,
			SingularTable: true,
		},
	})

	dbs[connectName] = db

	return db
}

// 获取Db对象
func Db(connectName ...string) *gorm.DB {
	if len(connectName) > 0 { //指定连接
		return initDb(connectName[0])
	} else { //默认数据库连接
		return initDb("default")
	}
}
