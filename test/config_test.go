package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfigByViper(t *testing.T) {
	var (
		ConfEnvPath  = "conf/dev"
		ViperConfMap = map[string]*viper.Viper{}
	)
	f, err := os.Open(ConfEnvPath + "/")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fileList, err := f.Readdir(1024)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, f0 := range fileList {
		if !f0.IsDir() {
			bts, err := ioutil.ReadFile(ConfEnvPath + "/" + f0.Name())
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			v := viper.New()
			v.SetConfigType("yaml")
			v.ReadConfig(bytes.NewBuffer(bts))
			pathArr := strings.Split(f0.Name(), ".")
			if ViperConfMap == nil {
				ViperConfMap = make(map[string]*viper.Viper)
			}
			ViperConfMap[pathArr[0]] = v
		}
	}

	fmt.Printf("\n%+v\n", ViperConfMap)
	fmt.Printf("\n%+v\n", ViperConfMap["base"].GetString("server_pord"))

	t.Log("ok")
}
