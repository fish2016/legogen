package config

import (
	"github.com/jinzhu/configor"
	"legogen/utils"
	pathlib "path"
)

type GenCodeScriptLoader struct {
	Config GenCodeScriptType
}

type GenCodeScriptType struct {
	PrjName string           `yaml:"prjname" required:"true"`
	Mode    GenMode          `yaml:"" default:""`
	Scripts []GenCodeCmdItem `yaml:"scripts" required:"true"`
}

type GenCodeCmdItem struct {
	Cmd  string `required:"true"`
	Conf string `required:"true"`
}

// 返回值bool表示配置文件是否存在
func (t *GenCodeScriptLoader) LoadConfig(dir string) bool {

	if dir == "" {
		pwd := utils.GetPwd()
		dir = pwd
	}

	//fmt.Println("LoadConfig", dir)
	_path := pathlib.Join(dir, ".gencode.yaml")
	exist := utils.CheckExist(_path)
	if !exist {
		return false
	}
	err := configor.Load(&t.Config, _path)
	if err != nil {
		panic(err)
	}

	return true
}

func (t *GenCodeScriptLoader) GetCmdItem(cmd string) *GenCodeCmdItem {
	for i, v := range t.Config.Scripts {
		if v.Cmd == cmd {
			return &t.Config.Scripts[i]
		}
	}

	return nil
}
