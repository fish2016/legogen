package generator

import (
	"bytes"
	"fmt"
	. "github.com/fish2016/legogen/logger"
	"github.com/fish2016/legogen/utils"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"log"

	"strings"
)

const (
	SubCmdTempl       = "templ"
	SubCmdFlagLogging = "logging"
	SubCmdFlagUrl     = "url"
)

// Generator holds the state of the analysis.  Primarily used to buffer the
// output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated ouptu.
	pkg *Package
}

// Printf writes the given output to the internalized buffer.
func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) ParsePackageDir(directory string) {
	pkg, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		Logger.Fatalf("cannot process directory %s: %s", directory, err)
	}
	/*golang老版本的src模式已经废除
	d, e := os.Getwd()
	gopath := os.Getenv("GOPATH")

	if e != nil {
		log.Fatalf("Error Grabbing WD: %s\n", e)
	}

	prefix := filepath.Join(gopath, "src") + string([]rune{filepath.Separator})

	d, err = filepath.Rel(prefix, d)
	if err != nil {
		log.Fatalf("Unable to get a relative path: %s\n", err)
	}
	*/
	var names []string
	names = append(names, pkg.GoFiles...) //指定目录下的go文件

	names = utils.PrefixDirectory(directory, names)
	g.parsePackage(directory, names, nil)
}

func (g *Generator) ParsePackageFiles(names []string) {
	g.parsePackage(".", names, nil)
}

// parsePackage analyzes the signle package constructed from the named files.
// If text is non-nil, it is a string to be used instead of the content of the file,
// to be used for testing.  parsePackage exists if there is an error.
func (g *Generator) parsePackage(directory string, names []string, text interface{}) {
	var files []*File
	var astFiles []*ast.File
	g.pkg = new(Package)
	fs := token.NewFileSet()

	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}

		parsedFile, err := parser.ParseFile(fs, name, text, parser.ParseComments)
		//ast.Print(fs, parsedFile)
		for _, v := range parsedFile.Comments {
			str := v.Text()
			if strings.HasPrefix(str, SubCmdTempl) {
				lines := strings.Split(str, "\n")
				if len(lines) <= 0 {
					continue
				}
				var firstLine = lines[0]

				//typ := strings.TrimPrefix(firstLine, SubCmdExtra)
				lineTmp := strings.TrimPrefix(firstLine, SubCmdTempl)
				arrString := strings.Split(lineTmp, " ")
				typ := arrString[0]

				//check the validity of subcmd flag
				if (typ != SubCmdFlagLogging) && (typ != SubCmdFlagUrl) {
					Logger.Fatalf("%s: invalid subcmd flag: <%s>", directory, typ)
				}

				lineTmp = strings.TrimPrefix(lineTmp, typ)
				lineTmp = strings.TrimSpace(lineTmp)

				if len(lineTmp) > 0 {
					Extras[typ] += (lineTmp + "\n")
				}

				if len(lines) > 1 {
					//extras[typ] = strings.Join(lines[1:], "\n")
					Extras[typ] += strings.Join(lines[1:], "\n")
				}
			}
		}

		if err != nil {
			log.Fatalf("parsing package: %s: %s", name, err)
		}
		astFiles = append(astFiles, parsedFile)
		files = append(files, &File{
			file:     parsedFile,
			Pkg:      g.pkg,
			path:     directory,
			fileName: name,
		})
	}

	if len(astFiles) == 0 {
		log.Fatalf("%s: no buildable Go files", directory)
	}
	g.pkg.name = astFiles[0].Name.Name
	g.pkg.files = files
	g.pkg.Dir = directory
	g.pkg.check(fs, astFiles)
	//fmt.Println("parsePackage g.pkg.name:", g.pkg.name)
}

// generate does 'things'
func (g *Generator) Generate(typeName string, summarize *string, middlewaresToGenerate *string, mapVal *CmdMapValue) {
	// pre-process
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		if file.file != nil {
			ast.Inspect(file.file, file.genImportsAndTypes)
		}
	}

	if *summarize != "" {
		g.pkg.Summarize()
		return
	}

	var targetFile *File

	for _, file := range g.pkg.files {
		for _, i := range file.Interfaces {
			if i.name == typeName {
				targetFile = file
				break
			}
		}
		//检测struct列表
		for _, i := range file.Structs {
			if i.name == typeName {
				targetFile = file
				break
			}
		}
	}

	if targetFile == nil {
		Logger.Fatalf("Unable to fine the type specified: %s\n", typeName)
	}

	//这里再对targetFile检测一道，检查目的是，跳出interfaces和structs中对应-type=xxx的类型
	//命令行中配置-type=xxx对应的类型，才是代码生成的依据
	var foundInterface *Interface
	for i, v := range targetFile.Interfaces {
		if v.name == typeName {
			foundInterface = &targetFile.Interfaces[i]
			break
		}
	}

	if foundInterface == nil {
		targetFile.Interfaces = nil
	} else {
		targetFile.Interfaces = []Interface{
			*foundInterface,
		}
	}

	//检测struct列表
	var foundStruct *Struct
	for i, v := range targetFile.Structs {
		if v.name == typeName {
			foundStruct = &targetFile.Structs[i]
			break
		}
	}

	if foundStruct == nil {
		targetFile.Structs = nil
	} else {
		targetFile.Structs = []Struct{
			*foundStruct,
		}
	}
	// begin generation
	list := strings.Split(*middlewaresToGenerate, ",")
	//list = append(list, "endpoint")
	list = append(list, "default")
	for _, l := range list {
		if bindings[l] != nil {
			bindings[l](g, targetFile, mapVal)
		}
	}
}

// format returns gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		Logger.Printf("warning: internal error: invalid Go generated: %s", err)
		Logger.Printf("warning: compile the pacakge to analyze the error")
		return g.buf.Bytes()
	}
	return src
}
