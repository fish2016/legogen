package config

import (
	"github.com/fish2016/legogen/utils"
	"github.com/jinzhu/configor"
	pathlib "path"
	"path/filepath"
)

type GenMode string

const (
	GenMode_Default     GenMode = ""
	GenMode_NoInputCode GenMode = "no-input-code"
)

type GenTemplateConfig struct {
	Templ string  `required:"true"`
	Out   string  `required:"true"`
	Mode  GenMode `default:""`
}

type GenConfig struct {
	AppName string `default:"legogen"`

	//用于import "{{.PrjName}}/internal/models"语句中，生成代码工程路径的root路径，表示import应用当前工程
	PrjName string  `yaml:"-"` //这里先不从配置文件读取，从命令行参数读取；暂时借用这个结构的字段，这个GenConfig结构会注入渲染
	Mode    GenMode `default:""`

	In string ``

	Template []GenTemplateConfig
}

var Config = GenConfig{}

func LoadConfig(dir string) {

	if dir == "" {
		pwd := utils.GetPwd()
		dir = pwd
	}

	var _path string
	if utils.IsDirectory(dir) {
		_path = pathlib.Join(dir, "gencode_config.yml")
	} else {
		_path = dir
	}
	//fmt.Println("LoadConfig", dir)

	err := configor.Load(&Config, _path)
	if err != nil {
		panic(err)
	}

	if Config.Template == nil {
		panic("no template config")
	}
	//configJson, _ := json.MarshalIndent(&Config, "", "    ")
	//fmt.Printf("load config: %s\n", configJson)
	var _dir string
	if !utils.IsDirectory(dir) {
		_dir = pathlib.Dir(dir) //取目录部分
	} else {
		_dir = dir
	}
	for idx, _ := range Config.Template {
		configtemplate := &Config.Template[idx]

		//这里意图是，配置文件里面是相对配置文件目录的路径
		configtemplate.Templ = pathlib.Join(_dir, configtemplate.Templ)
		configtemplate.Templ = filepath.ToSlash(configtemplate.Templ)

		configtemplate.Out = pathlib.Join(_dir, configtemplate.Out)
		configtemplate.Out = filepath.ToSlash(configtemplate.Out)

	}

	Config.In = pathlib.Join(_dir, Config.In)
	Config.In = filepath.ToSlash(Config.In)

	return
}
