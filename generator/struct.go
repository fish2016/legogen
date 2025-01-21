package generator

import (
	"go/ast"
)

// Value represents a declared constant.
type Struct struct {
	name  string // the name of the constant.
	types []Type

	pkg  *Package
	file File

	filedlist []Filed
}

type Filed struct {
	Name       string
	TypeName   string
	SelectorX  string
	ImportName string
	ImportPath string

	IsEmbedStruct bool //是否是内置结构
	CommentStr    []string
	//IsShowInView bool //在管理端界面是否显示这个字段
	//IsOnlyRead bool //管理界面edit上的字段不能被编辑，只能被显示时

	//types []Type
}

func createStruct(name string, astStru *ast.StructType, file File) Struct {
	stru := Struct{
		name:  name,
		types: nil,
	}
	stru.filedlist = make([]Filed, 0)
	//iface.Fields.List[0].Tag
	for _, field := range astStru.Fields.List {
		f := new(Filed)
		switch v := field.Type.(type) {
		case *ast.Ident:
			if len(field.Names) <= 0 {
				continue
			}

			f.Name = field.Names[0].Name

			typ := field.Type.(*ast.Ident)
			f.TypeName = typ.Name
			//stru.filedlist = append(stru.filedlist, *f)
		case *ast.SelectorExpr:
			//f := new(Filed)
			XVal, _ := v.X.(*ast.Ident)
			f.SelectorX = XVal.Name
			//SelVal, _:= v.Sel.( *ast.Ident)
			f.TypeName = v.Sel.Name
			if len(field.Names) > 0 {
				f.Name = field.Names[0].Name //如果类型为表达式，且变量名为非空，如： State lego.STATE_TYPE
			} else {
				//如果类型表达式格式形如
				/*
					type StruName struct {
					    lego.STATE_TYPE
					}
				*/
				f.Name = v.Sel.Name //字段的名字和类型都设置为一样的
			}

			f.IsEmbedStruct = true
			//typ := filed.Type.(*ast.Ident)
			//f.TypeName = typ.Name

		}

		if field.Comment != nil {
			f.CommentStr = make([]string, len(field.Comment.List))
			for i, v := range field.Comment.List {
				f.CommentStr[i] = v.Text
			}
		}

		stru.filedlist = append(stru.filedlist, *f)
	}
	return stru
}
