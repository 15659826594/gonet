package src

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gta-spec/utils"
)

const (
	supportMinGoVer = uint64(26)
	DS              = string(filepath.Separator)
	AppPath         = "app/"
	ConfPath        = "src/config/"
	InstallPath     = "app/admin/command/install/"
)

var (
	EnvGinMode = gin.DebugMode
	MainPath   string
	GoVersion  string
	RootPath   string
)

func init() {
	RootPath, _ = os.Getwd()
	info, ok := debug.ReadBuildInfo()
	if ok {
		if info.Main.Path != "" {
			MainPath = info.Main.Path
		} else {
			_, filename, _, _ := runtime.Caller(0)
			MainPath = readModuleName(utils.GoModFilepath(filepath.Dir(filename)))
		}
		GoVersion = info.GoVersion
		if v, e := getMinVer(GoVersion); e == nil && v < supportMinGoVer {
			log.Fatal(fmt.Sprintf("[ERROR] Now Gin requires Go 1.%d+.", supportMinGoVer))
		}
	}
}

func getMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

func readModuleName(goModPath string) string {
	file, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}

	return ""
}
