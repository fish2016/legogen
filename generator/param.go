package generator

import (
	"fmt"
	"go/ast"

	"strings"
)

type Param struct {
	Names []string
	Typ   Type
}

func createParam(field *ast.Field, reservedNames []string, suggestion string, file File) Param {
	p := Param{
		Names: make([]string, 0, len(field.Names)),
		Typ:   createType(field.Type, file.Pkg),
	}

	for _, n := range field.Names {
		p.Names = append(p.Names, n.Name)
	}

	// no name specified, let's create one...
	if len(p.Names) == 0 {
		n := DetermineLocalName(suggestion, reservedNames)
		p.Names = []string{n}
	}

	return p
}

func (p Param) ParamSpec() string {
	if len(p.Names) > 0 {
		return fmt.Sprintf("%s %s", strings.Join(p.Names, ", "), p.Typ)
	}
	return p.Typ.String()
}
