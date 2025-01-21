package generator

var bindings map[string]func(*Generator, *File, *CmdMapValue)
var Extras map[string]string

func init() {
	bindings = make(map[string]func(*Generator, *File, *CmdMapValue))
	Extras = make(map[string]string)
}
func RegisterProcess(argument string, fun func(*Generator, *File, *CmdMapValue)) {
	bindings[argument] = fun
}
