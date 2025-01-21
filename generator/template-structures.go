package generator

import (
	"fmt"
	"strings"
)

type TemplateCommon struct {
	BasePackage         string
	BasePackageImport   string
	BasePackageName     string
	EndpointPackage     string
	EndpointPackageName string
	EndpointPrefix      string
	InterfaceName       string
	InterfaceNameLcase  string
	StructName          string
	StructNameLcase     string
}

type TemplateParam struct {
	PublicName string
	Name       string
	Type       string
	IsContext  bool
}

func createTemplateParam(p Param) TemplateParam {
	return TemplateParam{
		Type: p.Typ.String(),
	}
}

type TemplateMethod struct {
	TemplateCommon
	HasContextParam        bool
	ContextParamName       string
	HasErrorResult         bool
	ErrorResultName        string
	HasMoreThanOneResult   bool
	LocalName              string
	MethodName             string
	MethodNameLcase        string
	MethodArguments        string
	MethodResults          string
	MethodResultNamesStr   string
	MethodArgumentNamesStr string
	MethodArgumentNames    []string
	MethodResultNames      []string
	Params                 []TemplateParam
	Results                []TemplateParam
}

func publicVariableName(str string) string {
	firstLetter := string([]byte{str[0]})
	rest := ""
	if len(str) > 1 {
		rest = str[1:]
	}

	return strings.ToUpper(firstLetter) + rest
}

func privateVariableName(str string) string {
	firstLetter := string([]byte{str[0]})
	rest := ""
	if len(str) > 1 {
		rest = str[1:]
	}

	return strings.ToLower(firstLetter) + rest
}

func createTemplateMethods(basePackage, endpointPackage *Import, interf Interface, methods []Method, reseveredNames []string) []TemplateMethod {
	results := make([]TemplateMethod, 0, len(methods))
	for _, meth := range methods {
		var names []string
		names = append(names, reseveredNames...)
		names = append(names, meth.UsedNames()...)
		var params []TemplateParam
		var methodsResults []TemplateParam

		var paramNames []string
		for _, p := range meth.params {
			// skip the context name, as this is primarily used to retrieve
			// values after transport.
			if !meth.hasContextParam || meth.contextParamName != p.Names[0] {
				paramNames = append(paramNames, p.Names...)
			}
			for _, n := range p.Names {
				param := TemplateParam{
					PublicName: publicVariableName(n),
					Name:       n,
					Type:       p.Typ.String(),
				}
				if param.Type == "context.Context" {
					param.IsContext = true
				}

				params = append(params, param)
			}
		}

		var resultNames []string
		for _, r := range meth.results {
			resultNames = append(resultNames, r.Names...)
			for _, n := range r.Names {
				methodsResults = append(methodsResults, TemplateParam{
					PublicName: publicVariableName(n),
					Name:       n,
					Type:       r.Typ.String(),
				})
			}
		}

		contextParamName := "_ctx"
		if meth.hasContextParam {
			contextParamName = meth.contextParamName
		}

		errorResultName := "err"
		if meth.hasErrResult {
			errorResultName = meth.errorResultName
		}

		lcaseName := DetermineLocalName(strings.ToLower(interf.name), reseveredNames)
		results = append(results, TemplateMethod{
			TemplateCommon: TemplateCommon{
				BasePackage:         basePackage.path,
				BasePackageImport:   basePackage.ImportSpec(),
				BasePackageName:     basePackage.name,
				EndpointPackage:     endpointPackage.path,
				EndpointPackageName: endpointPackage.name,
				EndpointPrefix:      fmt.Sprintf("/%s", strings.ToLower(interf.name)),
				InterfaceName:       interf.name,
				InterfaceNameLcase:  privateVariableName(interf.name),
			},
			HasContextParam:        meth.hasContextParam,
			ContextParamName:       contextParamName,
			HasErrorResult:         meth.hasErrResult,
			ErrorResultName:        errorResultName,
			HasMoreThanOneResult:   meth.moreThanOneResult,
			MethodName:             meth.name,
			MethodNameLcase:        privateVariableName(meth.name),
			LocalName:              lcaseName,
			MethodArguments:        meth.methodArguments(),
			MethodResults:          meth.methodResults(),
			MethodArgumentNamesStr: meth.methodArgumentNames(),
			MethodResultNamesStr:   meth.methodResultNames(),
			MethodArgumentNames:    paramNames,
			MethodResultNames:      resultNames,
			Params:                 params,
			Results:                methodsResults,
		})
	}
	return results
}

type TemplateBase struct {
	TemplateCommon
	Imports            []string //
	ImportsWithoutTime []string
	ExtraImports       []string //表示，interface中内嵌的成员变量（不包含函数的参数）引用到的import
	Methods            []TemplateMethod
	ExtraInterfaces    []TemplateParam
}

func CreateTemplateBase(basePackage, endpointPackage *Import, i Interface, s Struct, oimps []*Import) TemplateBase {
	// imps := filteredImports(i, oimps)
	imps := oimps

	names := make([]string, 0, len(imps))
	for _, i := range imps {
		names = append(names, i.name)
	}

	var impSpecs []string
	var impSpecsWithoutTime []string
	var extraImpSpecs []string
	for _, i := range imps {
		if !i.isParam && i.isEmbeded {
			extraImpSpecs = append(extraImpSpecs, i.ImportSpec())
			// skip non-params for these imports
			continue
		}

		if i.isParam {
			impSpecs = append(impSpecs, i.ImportSpec())
			if i.path != "time" {
				impSpecsWithoutTime = append(impSpecsWithoutTime, i.ImportSpec())
			}
		}
	}

	var extraInterfaces []TemplateParam
	for _, t := range i.types {
		var publicNamePieces = strings.Split(t.String(), ".")
		if len(publicNamePieces) < 1 {
			panic("This type is empty?!")
		}

		var publicName = publicNamePieces[len(publicNamePieces)-1]

		extraInterfaces = append(extraInterfaces, TemplateParam{
			PublicName: publicName,
			Name:       publicName,
			Type:       t.String(),
		})
	}

	return TemplateBase{
		TemplateCommon: TemplateCommon{
			BasePackage:         basePackage.path,
			BasePackageImport:   basePackage.ImportSpec(),
			BasePackageName:     basePackage.name,
			EndpointPackage:     endpointPackage.path,
			EndpointPackageName: endpointPackage.name,
			EndpointPrefix:      fmt.Sprintf("/%s", strings.ToLower(i.name)),
			InterfaceName:       i.name,
			//InterfaceNameLcase:  privateVariableName(i.name), fixme后续修改输入参数i可能是nil
			StructName:      s.name,
			StructNameLcase: privateVariableName(s.name),
		},
		Imports:            impSpecs,
		ImportsWithoutTime: impSpecsWithoutTime,
		ExtraImports:       extraImpSpecs,
		Methods:            createTemplateMethods(basePackage, endpointPackage, i, i.methods, names),
		ExtraInterfaces:    extraInterfaces,
	}
}

func filteredExtraImports(i Interface, imps []Import) []Import {
	var res []Import
	var tmp []string
	for _, imp := range imps {
		for _, t := range i.types {
			if strings.HasPrefix(t.String(), fmt.Sprintf("%s.", imp.name)) {
				if !sliceContains(tmp, imp.ImportSpec()) {
					res = append(res, imp)
					tmp = append(tmp, imp.ImportSpec())
				}
			}
		}
	}
	return res
}
