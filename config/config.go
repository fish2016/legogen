package config

import (
	"github.com/jinzhu/configor"
	"legogen/utils"
	pathlib "path"
	"path/filepath"
)

type GenTemplateConfig struct {
	Templ string `required:"true"`
	Out   string `required:"true"`
}

type GenConfig struct {
	AppName string `default:"legogen"`

	//用于import "{{.PrjName}}/internal/models"语句中，生成代码工程路径的root路径，表示import应用当前工程
	PrjName string `yaml:"-"` //这里先不从配置文件读取，从命令行参数读取；暂时借用这个结构的字段，这个GenConfig结构会注入渲染

	Template []GenTemplateConfig
}

var Config = GenConfig{}

func LoadConfig(dir string) {

	if dir == "" {
		pwd := utils.GetPwd()
		dir = pwd
	}

	//fmt.Println("LoadConfig", dir)
	_path := pathlib.Join(dir, "gencode_config.yml")
	err := configor.Load(&Config, _path)
	if err != nil {
		panic(err)
	}

	if Config.Template == nil {
		panic("no template config")
	}
	//configJson, _ := json.MarshalIndent(&Config, "", "    ")
	//fmt.Printf("load config: %s\n", configJson)

	for idx, _ := range Config.Template {
		configtemplate := &Config.Template[idx]
		configtemplate.Templ = pathlib.Join(dir, configtemplate.Templ)
		configtemplate.Templ = filepath.ToSlash(configtemplate.Templ)

		configtemplate.Out = pathlib.Join(dir, configtemplate.Out)
		configtemplate.Out = filepath.ToSlash(configtemplate.Out)

	}

	return
}
