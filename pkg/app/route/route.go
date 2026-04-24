package route

import (
	"fmt"
	"gota/pkg"
	"gota/pkg/utils"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var controllers []Controller

type Controller struct {
	File         string
	Name         string
	Alias        string
	Method       []string
	NoNeedLogin  []string
	NoNeedRight  []string
	HandlersFunc []gin.HandlerFunc
	Actions      []Action
}

type Action struct {
	Name        string
	Path        []string
	Method      []string
	HandlerFunc gin.HandlerFunc
}

type IInitialize interface {
	Initialize() gin.HandlerFunc
}

type IBeforeAction interface {
	BeforeAction() gin.HandlerFunc
}

func Register(Struct any) {
	t := reflect.TypeOf(Struct)

	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
	case reflect.Struct:
	default:
		return
	}

	_, file, _, _ := runtime.Caller(1)

	relPath, _ := filepath.Rel(pkg.RootPath(), file)

	controller := Controller{
		File: relPath,
		Name: t.Name(),
		// PUT‌：用于‌全量替换‌资源。客户端需发送整个资源的最新状态，服务器会用该请求体完全覆盖原有资源
		// PATCH‌：用于‌部分更新‌资源。客户端仅需发送需要修改的字段或一组修改指令，服务器根据这些指令对现有资源进行局部调整
		Method:       []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
		HandlersFunc: []gin.HandlerFunc{},
	}

	props := utils.GetStructProperty(Struct, "Alias", "Method", "NoNeedLogin", "NoNeedRight")

	// 处理别名
	if aliasAny, ok := props["Alias"]; ok {
		if alias, _ := aliasAny.(string); alias != "" {
			controller.Alias = alias
		}
	}

	if methodAny, ok := props["Method"]; ok {
		if method, _ := methodAny.([]string); len(method) > 0 {
			controller.Method = method
		}
	}

	if noNeedLoginAny, ok := props["NoNeedLogin"]; ok {
		controller.NoNeedLogin, _ = noNeedLoginAny.([]string)
	}

	if noNeedRightAny, ok := props["NoNeedRight"]; ok {
		controller.NoNeedRight, _ = noNeedRightAny.([]string)
	}

	if handler, ok := Struct.(IInitialize); ok {
		controller.HandlersFunc = append(controller.HandlersFunc, handler.Initialize())
	}

	if handler, ok := Struct.(IBeforeAction); ok {
		controller.HandlersFunc = append(controller.HandlersFunc, handler.BeforeAction())
	}

	for name, method := range utils.GetStructMethods(Struct) {
		action := Action{
			Name:   name,
			Path:   []string{name},
			Method: controller.Method,
		}
		switch fun := method.(type) {
		case func(*gin.Context):
			action.HandlerFunc = fun
		case func() (gin.HandlerFunc, []string, []string):
			action.HandlerFunc, action.Path, action.Method = fun()
		default:
			continue
		}
		controller.Actions = append(controller.Actions, action)
	}

	controllers = append(controllers, controller)
}

func Build(e *gin.Engine, moduleGroup func(name string) (*gin.RouterGroup, string)) {
	for _, controller := range controllers {
		var cRouterGroup *gin.RouterGroup
		var chains []gin.HandlerFunc
		mRouterGroup, modulename := moduleGroup(controller.File)

		if controller.Alias != "" {
			cRouterGroup = mRouterGroup.Group(utils.CaseSnake(controller.Alias))
		} else {
			cRouterGroup = mRouterGroup.Group(utils.CaseSnake(controller.Name))
		}

		chains = append(chains, controller.HandlersFunc...)

		for _, action := range controller.Actions {
			for _, method := range action.Method {
				for _, path := range action.Path {
					path = filepath.ToSlash(utils.CaseSnake(path))
					clonedChains := []gin.HandlerFunc{func(c *gin.Context) {
						c.Set("startTime", time.Now().UnixMilli())
						c.Set("modulename", utils.CaseSnake(modulename))
						c.Set("controllername", utils.CaseSnake(controller.Name))
						c.Set("actionname", utils.CaseSnake(action.Name))
						c.Set("noNeedLogin", controller.NoNeedLogin)
						c.Set("noNeedRight", controller.NoNeedRight)
						c.Set("url", fmt.Sprintf("%s/%s/%s", c.GetString("modulename"), c.GetString("controllername"), c.GetString("actionname")))
					}}
					clonedChains = append(clonedChains, chains...)
					clonedChains = append(clonedChains, action.HandlerFunc)

					switch {
					case strings.HasPrefix(path, "/"):
						e.Handle(method, path, clonedChains...)
					case strings.HasPrefix(path, "."):
						e.Handle(method, filepath.ToSlash(filepath.Clean(filepath.Join(cRouterGroup.BasePath(), utils.CaseSnake(path)))), clonedChains...)
					default:
						cRouterGroup.Handle(method, path, clonedChains...)
					}
				}
			}
		}
	}
}
