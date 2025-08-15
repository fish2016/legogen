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
	pathlib "path"
	"strings"

	_ "legogen/process"
)

const VERSION = "0.2"

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
	fmt.Fprintf(os.Stderr, "lego code generator %s\n", VERSION)
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s [flags] -type T [directory] -project=[prj-name] --dirconf=[dirconf] \n", "legocli")
}

type GenCmdParam struct {
	TypeName  string //处理目标结构体名称
	ConfPath  string //模板配置文件路径
	InputPath string //输入代码文件路径
	PrjName   string
	Mode      config.GenMode
}

var genCmdParam GenCmdParam

func main() {
	flag.Var(userDefineVal, "kv", "kv usage: kv MoudleName=User, KeyName2=Vaule2")
	flag.Usage = Usage
	flag.Parse()
	binaryName = os.Args[0]

	//判断当前工作目录是否存在.gencmd.yml配置文件，有则进入预定义命令模式（所谓预定义命令模式，是这个配置文件里面配置了预定义的命令）
	//检测命令行的 -type参数是否为空，空则不执行
	scripLoader := config.GenCodeScriptLoader{}
	exist := scripLoader.LoadConfig("")
	if exist {
		if len(os.Args) < 3 {
			fmt.Println("no command")
			os.Exit(2)
		}
		cmd := os.Args[1]
		typName := os.Args[2]

		cmdItem := scripLoader.GetCmdItem(cmd)
		if cmdItem == nil {
			fmt.Println("not support command:", cmd)
			os.Exit(2)
		}

		genCmdParam.TypeName = typName
		genCmdParam.ConfPath = cmdItem.Conf
		genCmdParam.PrjName = scripLoader.Config.PrjName

	} else {
		//没有配置文件，则进入普通模式
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

		//if *noInputCode == true {
		//	config.LoadConfig(*dirOfConf)
		//
		//	if config.Config.PrjName == "" {
		//		config.Config.PrjName = *prjName
		//	}
		//
		//	for _, typeName := range types {
		//		process.ProcessWithNoInputCode(typeName, *prjName)
		//	}
		//
		//	return
		//}

		genCmdParam.TypeName = types[0]
		genCmdParam.ConfPath = *dirOfConf
		genCmdParam.PrjName = config.Config.PrjName
		if *noInputCode == true {
			genCmdParam.Mode = config.GenMode_NoInputCode
		}

	}

	config.LoadConfig(genCmdParam.ConfPath)
	var (
		gen generator.Generator
	)

	if config.Config.PrjName == "" {
		config.Config.PrjName = *prjName
	}

	if config.Config.PrjName == "" {
		config.Config.PrjName = genCmdParam.PrjName
	}

	genCmdParam.Mode = config.Config.Mode
	genCmdParam.InputPath = pathlib.Join(config.Config.In, fmt.Sprintf("%s.go", utils.CamelToSnake(genCmdParam.TypeName))) //fixme 这里暂时这样以第一个配置为全局input处理,目前不能没一行一个input，目前是全局一个input
	if genCmdParam.Mode == config.GenMode_NoInputCode {
		process.ProcessWithNoInputCode(genCmdParam.TypeName, *prjName)
	} else {
		if genCmdParam.InputPath == "" {
			panic("no input path")
		}
		if utils.IsDirectory(genCmdParam.InputPath) {
			gen.ParsePackageDir(genCmdParam.InputPath)
		} else {
			gen.ParsePackageFiles([]string{genCmdParam.InputPath})
		}
		middlewaresToGenerate := ""
		gen.Generate(genCmdParam.TypeName, summarize, &middlewaresToGenerate, userDefineVal)
	}

	return
}
