package process

import (
	"bytes"
	"fmt"
	"github.com/fish2016/legogen/config"
	"github.com/fish2016/legogen/utils"
	"log"
	"path/filepath"
	"text/template"
)

type TemplateTargetPath struct {
	TargetType      string
	TargetTypeSnake string
	PrjName         string
}

// 处理无输入源代码的模板，这里意思是，不依赖基础go代码生成（不需要解析go代码）
func ProcessWithNoInputCode(typ string, prjName string) {

	// 遍历所有模板配置
	for _, templateConfig := range config.Config.Template {
		var buf bytes.Buffer

		templatePath := templateConfig.Templ
		outFilename := templateConfig.Out

		var pathBuf bytes.Buffer
		pathTmpl, err := template.New("path").Parse(outFilename)
		if err != nil {
			log.Fatalf("Path Template Parsing Error: %s", err)
		}

		files := []string{
			templatePath,
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatalf("Template Parsing Error: %s", err)
		}

		data := TemplateTargetPath{
			TargetType:      typ,
			TargetTypeSnake: utils.CamelToSnake(typ),
			PrjName:         prjName,
		}
		err = tmpl.Execute(&buf, data)
		if err != nil {
			log.Fatalf("Template execution failed: %s\n", err)
		}

		err = pathTmpl.ExecuteTemplate(&pathBuf, "path", data)
		if err != nil {
			log.Fatalf("Path Template Parsing Error: %s", err)
		}

		outFilename = pathBuf.String()

		func() {
			file := utils.OpenFile(filepath.Dir(outFilename), filepath.Base(outFilename))
			defer file.Close()

			fmt.Fprint(file, string(utils.FormatBuffer(buf, outFilename)))

			outFilenameAbs, _ := filepath.Abs(outFilename)
			log.Printf("generate file: %s, (%s%s)\n", outFilenameAbs, filepath.Dir(outFilename), filepath.Base(outFilename))
		}()

		// 清空buffer为下一个模板做准备
		buf.Reset()
		pathBuf.Reset()
	}

	fmt.Println("generate code success")

}

//func init() {
//	generator.RegisterProcess("no_input_code", processWithNoInputCode)
//}
