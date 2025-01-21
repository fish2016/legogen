package generator

import (
	"fmt"
	"go/ast"
	"path"
	// "path/filepath"
	"strings"
)

type Import struct {
	name      string
	path      string
	last      string
	isEmbeded bool //意思是，这个import的库有使用过，后面生成就需要加入，如果没有使用过，就不需要生成对应的import语句
	isParam   bool
}

func CreateImportWithPath(p string) *Import {
	last := path.Base(p)
	name := last
	if strings.Contains(last, "-") {
		lastPieces := strings.Split(last, "-")
		name = lastPieces[len(lastPieces)-1]
	}
	return &Import{
		name: name,
		path: p,
		last: last,
	}
}

func CreateImport(imp *ast.ImportSpec) *Import {
	var name string
	pth := strings.TrimPrefix(strings.TrimSuffix(imp.Path.Value, "\""), "\"")
	last := path.Base(pth)
	if n := imp.Name; n == nil {
		name = last
	} else {
		name = n.String()
	}

	if strings.Contains(name, "-") {
		namePieces := strings.Split(name, "-")
		name = namePieces[len(namePieces)-1]
	}

	return &Import{
		name: name,
		path: pth,
		last: last,
	}
}

func (i Import) ImportSpec() string {
	if i.name == i.last {
		return fmt.Sprintf("\"%s\"", i.path)
	}

	return fmt.Sprintf("%s \"%s\"", i.name, i.path)
}
