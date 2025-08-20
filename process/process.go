package process

import (
	"bytes"
	"fmt"
	"github.com/fish2016/legogen/config"
	"github.com/fish2016/legogen/generator"
	"github.com/fish2016/legogen/utils"
	"log"
	"path"
	"path/filepath"
	"text/template"
)

func processDefault(g *generator.Generator, f *generator.File, mapValue *generator.CmdMapValue) {
	//gopath := os.Getenv("GOPATH")
	convertedPath := filepath.ToSlash(f.Pkg.Dir)
	endpointPackage := generator.CreateImportWithPath(path.Join(convertedPath, "endpoint"))
	basePackage := generator.CreateImportWithPath(convertedPath)

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

		for _, s := range f.Structs {
			data := generator.CreateTemplateBase(config.Config, basePackage, endpointPackage, generator.Interface{}, s, f.Imports)
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
	}

	fmt.Println("generate code success")

}

func init() {
	generator.RegisterProcess("default", processDefault)
}
