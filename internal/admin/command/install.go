package command

import (
	"database/sql"
	"fmt"
	"gonet/pkg"
	"gonet/pkg/database"
	"gonet/pkg/utils"
	"gonet/pkg/utils/yaml"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"gonet/pkg/utils/ini"

	"github.com/gin-gonic/gin"
)

type Install struct {
	MinGoVersion string
}

func (i Install) Index(c *gin.Context) {
	installLockFile := pkg.INSTALL_PATH + "install.lock"
	if _, err := os.Stat(installLockFile); !os.IsNotExist(err) {
		c.String(http.StatusOK, fmt.Sprintf("The system has been installed. If you need to reinstall, please remove %s first", "install.lock"))
		return
	}
	if c.Request.Method == http.MethodPost {
		mysqlHostname := c.DefaultPostForm("mysqlHostname", "127.0.0.1")
		mysqlHostport := c.DefaultPostForm("mysqlHostport", "3306")
		// 处理主机名和端口
		hostArr := strings.Split(mysqlHostname, ":")
		if len(hostArr) > 1 {
			mysqlHostname = hostArr[0]
			mysqlHostport = hostArr[1]
		}
		mysqlUsername := c.DefaultPostForm("mysqlUsername", "root")
		mysqlPassword := c.DefaultPostForm("mysqlPassword", "")
		mysqlDatabase := c.DefaultPostForm("mysqlDatabase", "")
		mysqlPrefix := c.DefaultPostForm("mysqlPrefix", "fa_")
		adminUsername := c.DefaultPostForm("adminUsername", "admin")
		adminPassword := c.DefaultPostForm("adminPassword", "")
		adminPasswordConfirmation := c.DefaultPostForm("adminPasswordConfirmation", "")
		adminEmail := c.DefaultPostForm("adminEmail", "admin@admin.com")
		siteName := c.DefaultPostForm("siteName", "My Website")

		// 验证密码匹配
		if adminPassword != adminPasswordConfirmation {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "The two passwords you entered did not match",
			})
			return
		}
		adminName, err := i.installation(mysqlHostname, mysqlHostport, mysqlDatabase, mysqlUsername, mysqlPassword, mysqlPrefix, adminUsername, adminPassword, adminEmail, siteName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Install Successed",
			"url":  nil,
			"data": gin.H{
				"adminName": adminName,
			},
		})
		return
	}
	errInfo := ""
	c.HTML(http.StatusOK, "install.html", gin.H{
		"errInfo": errInfo,
	})
}

/**
 * 执行安装
 */
func (i Install) installation(mysqlHostname, mysqlHostport, mysqlDatabase, mysqlUsername, mysqlPassword, mysqlPrefix, adminUsername, adminPassword, adminEmail, siteName string) (string, error) {
	err := i.checkEnv()
	if err != nil {
		return "", err
	}

	// 修改后台入口
	adminName := utils.RandomAlpha(10)

	// ==============================初始化数据库-start===================================
	if mysqlDatabase == "" {
		return "", fmt.Errorf("please input correct database")
	}

	usernameRegex := regexp.MustCompile(`^\w{3,12}$`)
	if !usernameRegex.MatchString(adminUsername) {
		return "", fmt.Errorf("please input correct username")
	}

	passwordRegex := regexp.MustCompile(`^[\S]{6,16}$`)
	if !passwordRegex.MatchString(adminPassword) {
		return "", fmt.Errorf("please input correct password")
	}

	weakPasswordArr := []string{"123456", "12345678", "123456789", "654321", "111111", "000000", "password", "qwerty", "abc123", "1qaz2wsx"}
	for _, weakPass := range weakPasswordArr {
		if adminPassword == weakPass {
			return "", fmt.Errorf("password is too weak")
		}
	}

	if siteName == "" || regexp.MustCompile("(?i)fastadmin").MatchString(siteName) {
		return "", fmt.Errorf("please input correct website")
	}

	sqlFile := pkg.INSTALL_PATH + "fastadmin.sql"

	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		return "", err
	}

	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		return "", err
	}

	sqlStr := strings.ReplaceAll(string(sqlBytes), "`fa_", "`"+mysqlPrefix)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", mysqlUsername, mysqlPassword, mysqlHostname, mysqlHostport))
	if err != nil {
		return "", err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", mysqlDatabase))
	if err != nil {
		return "", fmt.Errorf("failed to create database: %v", err)
	}

	_, err = db.Exec(fmt.Sprintf("USE `%s`", mysqlDatabase))
	if err != nil {
		return "", fmt.Errorf("failed to connect to database: %v", err)
	}

	// 查询一次SQL,判断连接是否正常
	err = db.Ping()
	if err != nil {
		return "", fmt.Errorf("database connection test failed: %v", err)
	}

	// 执行SQL语句
	err = database.Exec(db, sqlStr)
	if err != nil {
		return "", fmt.Errorf("failed to execute SQL: %v", err)
	}

	// ==============================初始化数据库-end===================================

	// ==========================初始化环境变量文件-start================================

	if env, err := ini.Load(".env.sample"); err == nil {
		env.SetSectionOrderMap("database", map[string]string{
			"hostname": mysqlHostname,
			"hostport": mysqlHostport,
			"username": mysqlUsername,
			"password": mysqlPassword,
			"database": mysqlDatabase,
			"prefix":   mysqlPrefix,
		}, []string{"hostname", "hostport", "username", "password", "database", "prefix"})
		env.Save(".env")
	}

	// ==========================初始化环境变量文件-end================================

	// ============================变更配置文件-start=================================

	if y, err := yaml.Load(pkg.CONF_PATH + "ini.yaml"); err == nil {
		y.Set("token.key", utils.RandomAlnum(32))
		y.Save()
	}

	if y, err := yaml.Load(pkg.APP_PATH + "extra/site.yaml"); err == nil {
		y.Set("name", siteName)
		y.Save()
	}

	// =============================变更配置文件-end==================================

	// ============================变更默认管理员-start================================

	avatar := "/assets/img/avatar.png"
	// 变更默认管理员密码
	if adminPassword == "" {
		adminPassword = utils.RandomAlnum(8)
	}
	if adminEmail == "" {
		adminEmail = "admin@admin.com"
	}
	newSalt := utils.Substr(utils.Md5(utils.Uniqid("")), 0, 6)
	newPassword := utils.Md5(utils.Md5(adminPassword) + newSalt)
	data := map[string]string{
		"username": adminUsername,
		"email":    adminEmail,
		"avatar":   avatar,
		"password": newPassword,
		"salt":     newSalt,
	}
	var condition []string
	var args []any
	for key, val := range data {
		condition = append(condition, key+"=?")
		args = append(args, val)
	}
	query := fmt.Sprintf("UPDATE `%s` SET %s WHERE username='admin'", mysqlPrefix+"admin", strings.Join(condition, ","))
	_, _ = db.Exec(query, args...)

	// 变更前台默认用户的密码,随机生成
	newSalt = utils.Substr(utils.Md5(utils.Uniqid("")), 0, 6)
	newPassword = utils.Md5(utils.Md5(utils.RandomAlnum(8)) + newSalt)
	query = fmt.Sprintf("UPDATE `%s` SET avatar=?,password=?,salt=? WHERE username='admin'", mysqlPrefix+"user")
	_, _ = db.Exec(query, avatar, newPassword, newSalt)

	// ============================变更默认管理员-end================================

	query = fmt.Sprintf("UPDATE `%s` SET value=? WHERE name='name'", mysqlPrefix+"config")
	_, _ = db.Exec(query, siteName)

	// 创建install.lock文件 app install方法会检测到
	installLockFile := pkg.INSTALL_PATH + "install.lock"
	err = os.WriteFile(installLockFile, []byte(adminName), 0644)
	if err != nil {
		return "", fmt.Errorf("the current permissions are insufficient to write the file %s", installLockFile)
	}

	return adminName, nil
}

// 检测环境
func (i Install) checkEnv() error {

	if versionCompare(runtime.Version(), i.MinGoVersion) == -1 {
		return fmt.Errorf("the current GO %s is too low, please use GO %s or higher", runtime.Version(), i.MinGoVersion)
	}

	return nil
}

// 比较Go版本的辅助函数
func versionCompare(version1, version2 string) int {
	// 移除版本号前的 "go" 前缀（如果存在）
	v1 := strings.TrimPrefix(version1, "go")
	v2 := strings.TrimPrefix(version2, "go")

	// 按照 "." 分割版本号
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	// 获取较长的版本号长度
	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	// 逐个比较各部分
	for i := 0; i < maxLen; i++ {
		var num1, num2 int

		// 如果当前部分存在，则转换为整数，否则为0
		if i < len(parts1) {
			num1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			num2, _ = strconv.Atoi(parts2[i])
		}

		// 比较数字
		if num1 < num2 {
			return -1
		} else if num1 > num2 {
			return 1
		}
	}

	// 版本号相等
	return 0
}
