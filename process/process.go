package process

import (
	"bytes"
	"fmt"
	"legogen/config"
	"legogen/generator"
	"legogen/utils"
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

	var buf bytes.Buffer

	templatePath := config.Config.Template[0].Templ
	outFilename := config.Config.Template[0].Out

	var pathBuf bytes.Buffer
	pathTmpl, err := template.New("path").Parse(outFilename)
	if err != nil {
		log.Fatalf("Path Template Parsing Error: %s", err)
	}

	outFilename = pathBuf.String()

	files := []string{
		templatePath,
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatalf("Template Parsing Error: %s", err)
	}

	for _, s := range f.Structs {
		data := generator.CreateTemplateBase(basePackage, endpointPackage, generator.Interface{}, s, f.Imports)
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
		}()

	}

	fmt.Println("generate code success")

}

func init() {
	generator.RegisterProcess("default", processDefault)
}
