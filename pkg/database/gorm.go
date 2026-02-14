package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	instance *gorm.DB
	once     sync.Once
	prefix   string
)

func Gorm(configs ...*Database) *gorm.DB {
	once.Do(func() {
		var err error
		if len(configs) < 1 {
			log.Fatalf("database initialization failed: %v", errors.New("configuration is required"))
		}
		c := configs[0]

		loggerConfig := logger.Config{
			SlowThreshold:             time.Millisecond * 100, // 慢速 SQL 阈值 (100毫秒)
			LogLevel:                  logger.Silent,          // Log level
			IgnoreRecordNotFoundError: true,                   // 忽略record not found
			ParameterizedQueries:      false,                  // 开发阶段启用 可以快速定位 SQL 语句与参数不匹配的问题,生产环境中建议关闭避免敏感数据暴露
			Colorful:                  false,
		}

		if gin.IsDebugging() {
			loggerConfig.LogLevel = logger.Info
			loggerConfig.ParameterizedQueries = true
		}

		instance, err = gorm.Open(mysql.Open(c.GetDsn()), &gorm.Config{
			QueryFields:            true, // select 字段 而不是 select *
			SkipDefaultTransaction: true, // 禁用全局事务
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   c.Prefix, // 为所有表添加前缀
				SingularTable: true,     // 禁用复数表名
			},
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到控制台
				loggerConfig,
			),
		})
		prefix = c.Prefix
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		sqlDB, err := instance.DB()
		if err == nil {
			if c.MaxOpenConns != 0 {
				sqlDB.SetMaxOpenConns(c.MaxOpenConns)
			}
			if c.MaxIdleConns != 0 {
				sqlDB.SetMaxIdleConns(c.MaxIdleConns)
			}
			if c.ConnMaxLifetime == 0 {
				c.ConnMaxLifetime = getLifetime(sqlDB)
			}
			sqlDB.SetConnMaxLifetime(time.Second * time.Duration(c.ConnMaxLifetime))
		}
	})
	return instance
}

// getLifetime
// [mysql] 2025/03/11 17:36:01 packets.go:149 write tcp 127.0.0.1:50154->127.0.0.1:3306: wsasend: An established connection was aborted by the software in your host machine.
// [mysql] connection.go:173: bad connection 和 invalid connection
// maxLifeTime必须要比mysql服务器设置的wait_timeout小，否则会导致golang侧连接池依然保留已被mysql服务器关闭了的连接。
// SHOW VARIABLES LIKE 'wait_timeout'; 查看wait_timeout
// 一般建议wait_timeout/2
func getLifetime(db *sql.DB) int {
	variable := struct {
		variableName string
		waitTimeout  int
	}{}
	err := db.QueryRow(fmt.Sprintf("SHOW VARIABLES LIKE '%s'", "wait_timeout")).Scan(&variable.variableName, &variable.waitTimeout)
	if err == nil {
		return variable.waitTimeout / 2
	}
	return 14400
}

// GetTablePrefix 获取当前数据库连接的表前缀
func GetTablePrefix() string {
	return prefix
}
