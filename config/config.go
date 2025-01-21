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
	APPName string `default:"legogen"`

	Template []GenTemplateConfig
}

var Config = GenConfig{}

func LoadConfig(dir string) {

	if dir == "" {
		pwd := utils.GetPwd()
		dir = pwd
	}

	//fmt.Println("LoadConfig", dir)
	path := pathlib.Join(dir, "gencode_tmpl", "gencode_config.yml")
	err := configor.Load(&Config, path)
	if err != nil {
		panic(err)
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
