package test

import (
	"bytes"
	"fmt"
	"math/rand"
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
			bts, err := os.ReadFile(ConfEnvPath + "/" + f0.Name())
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

func TestRandInts(t *testing.T) {
	list := []string{"111", "222", "111", "222", "111", "222", "111", "222"}
	fmt.Println(rand.Int())
	for i := 0; i < 10; i++ {
		fmt.Printf("%v\n", rand.Int()%len(list))
	}
}
