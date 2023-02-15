package util

import "os"

const (
	EnvDebug = "debug"
	EnvDev   = "dev"
	EnvProd  = "Prod"
)

func GetRunEnv() string {
	e := os.Getenv("GORUNEVN")
	return e
}
