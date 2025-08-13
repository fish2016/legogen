package main

import (
	"flag"
	"fmt"
	"legogen/config"
	"legogen/generator"
	_ "legogen/logger"
	"legogen/process"
	"legogen/utils"
	"os"
	"strings"

	_ "legogen/process"
)

var (
	typeNames = flag.String("type", "", "comma-separated list of type names; must be set")
	prjName   = flag.String("project", "", "prj name; must be set")
	//middlewaresToGenerate = flag.String("middleware", "logging,instrumenting,transport,zipkin", "comma-seperated list of middlewares to process. Options: [logging,instrumenting,transport,zipkin]")

	summarize                            = flag.String("summarize", "", "Prints out the Summary of Found structures intead of generating code")
	binaryName                           = ""
	userDefineVal *generator.CmdMapValue = generator.NewCmdMapValue()

	dirOfConf = flag.String("dirconf", "", "the dir of config")

	noInputCode = flag.Bool("no-input-code", false, "prj name; option")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s [flags] -type T [directory] -project=[prj-name] --dirconf=[dirconf] \n", "legocli")
}

func main() {
	flag.Var(userDefineVal, "kv", "kv usage: kv MoudleName=User, KeyName2=Vaule2")
	flag.Usage = Usage
	flag.Parse()
	binaryName = os.Args[0]

	//检测命令行的 -type参数是否为空，空则不执行
	if len(*typeNames) == 0 {
		fmt.Println("-type must be set")
		flag.Usage()
		os.Exit(2)
	}

	if len(*prjName) == 0 {
		fmt.Println("-project must be set")
		flag.Usage()
		os.Exit(2)
	}

	//如果-type指定多个类型，则用逗号分割
	types := strings.Split(*typeNames, ",")

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	if *noInputCode == true {
		config.LoadConfig(*dirOfConf)

		if config.Config.PrjName == "" {
			config.Config.PrjName = *prjName
		}

		for _, typeName := range types {
			process.ProcessWithNoInputCode(typeName, *prjName)
		}

		return
	}
	var (
		dir string
		gen generator.Generator
	)

	if len(args) == 1 && utils.IsDirectory(args[0]) {
		dir = args[0]
		gen.ParsePackageDir(dir) //注意和下面的函数名区别一个后缀是Dir，下面的后缀是Files
		// parsePackageDir(args[0])
	} else {
		//dir = filepath.Dir(args[0])
		gen.ParsePackageFiles(args)
		// parsePackageFiles(args)
	}

	config.LoadConfig(*dirOfConf)

	if config.Config.PrjName == "" {
		config.Config.PrjName = *prjName
	}

	for _, typeName := range types {
		middlewaresToGenerate := ""
		gen.Generate(typeName, summarize, &middlewaresToGenerate, userDefineVal)
	}
	return
}
