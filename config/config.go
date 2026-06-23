package config

import (
	"os"

	"github.com/gta-spec/utils"
)

var (
	AppName = utils.FirstNonNil(os.Getenv("APP_NAME"), "skeleton")
	AppEnv  = utils.FirstNonNil(os.Getenv("APP_ENV"), "dev")
)
