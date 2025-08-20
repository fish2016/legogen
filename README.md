# LegoGen golang代码生成器
LegoGen 是一个基于 Golang 的 AST 和 text/template 实现的代码生成器。它通过解析 Go 文件中的 struct 或 interface 类型定义，结合预设的代码模板，自动生成相应的代码。此生成器的核心特点是，只需使用 Go 文件中的类型信息，便可高效生成所需代码。

# 快速开始
## 源代码编译
``````
git clone xxx
go mod tidy
go build

将生成执行文件legogen放置在环境变量Path路径下

``````

## go install安装方式
``````
go install github.com/fish2016/legogen@v0.2.4

``````

## 执行生成代码命令
``````
方法1: 使用go generate .
cd ./examples/gen_examplete_model
go generate .


``````

``````
方法2：
cd ./examples/gen_examplete_model
legogen -type=Example
``````

# 代码准备
## go文件中定义struct或interface

生成器会使用ast包解析go定义，将定义信息输出到一个对象里（这个对象会在渲染代码模板时用来生成代码）

会内部输出以下字段：

StructName: go结构名，下例中为Example     

StructNameLcase: go结构名(snake风格)，下例中为example

注意，这里加入go:generate指令，支持使用go generate . 命令
``````
//go:generate legogen -type=Example

type Example struct {
	Id    int64     `xorm:"pk autoincr"`      // ["id","hide","hidden"]
	IfDel int       `xorm:"TINYINT(1)"`       // ["deleted","hide","hidden"]
	Cdate time.Time `xorm:"DateTime created"` // ["create","show","datetime"]
	Udate time.Time `xorm:"DateTime updated"` // ["update","hide","datetime"]
	// add your custom field here
}

``````


## 代码模板例子：
可以看到，代码模板可以使用上一步输出的字段StructName，生成时会替换StructName的值
代码模板使用go模板语法
``````
func (t *{{.StructName}} ) RecodeById(id int64) *{{.StructName}}  {
	item := new({{.StructName}} )
	lego.GetDBEngine().ID(id).Get(item)
	if item.Id <= 0 {
		return nil
	}
	return item
}

func (t *{{.StructName}} ) AddRecode(item2add *{{.StructName}} ) bool {
	item2add.Id = lego.UUID()
	_, err := lego.GetDBEngine().Insert(item2add)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}
``````

## 代码生成配置问题
该文件定义模板文件的路径，以及输出代码文件的路径

路径里可以使用变量值，下例中使用"{{.StructNameLcase}}_sql_gen.go"作为文件名，生成example_sql_gen.go文件
``````
pkgname: example

template:
    - {templ: "gencode_tmpl/template_sql.go", out: "{{.StructNameLcase}}_sql_gen.go"}

``````

# 生成方法
``````
cd ./examples/gen_example_model
go generate .  //这里会触发调用example_model.go里面的指令 go:generate legogen -type=Example

生成代码会在example_sql_gen.go
``````