# goFastAdmin

基于 Gin 框架的 Go Web 快速开发框架。

## 推荐 IDE 配置

[GoLand](https://www.jetbrains.com.cn/go/)  + [TypeScript Vue Plugin (Volar)](https://marketplace.visualstudio.com/items?itemName=Vue.vscode-typescript-vue-plugin)。

## 自定义配置

详见 [Vite 配置参考](https://vitejs.dev/config/)。

## 生成 API 文档

```sh
go install github.com/swaggo/swag/cmd/swag@latest

swag init -o ./src/docs
```

### 本地项目关联
```sh
在项目所在的目录
go work init ./gt
go work use ./utils
```

### 清理依赖缓存

```sh
go clean -modcache

rm go.sum

go get github.com/gta-spec/utils@latest

go mod tidy
```

### 开发环境运行（热重载）

```sh
go generate

go run .
```

### 生产环境编译

```sh
go build -tags release -o myapp
go build -tags debug -o myapp
go build -tags test -o myapp
```

## 技术栈

- **Web 框架**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **配置管理**: [Viper](https://github.com/spf13/viper)
- **缓存**: [Redis](https://github.com/go-redis/redis/v8)
- **权限控制**: [Casbin](https://casbin.org/)
- **国际化**: [go-i18n](https://github.com/nicksnyder/go-i18n)
- **API 文档**: [Swag](https://github.com/swaggo/swag)

## 环境要求

- Go 1.26.3 或更高版本
- MySQL 5.7+ 或 8.0+
- Redis 6.0+
