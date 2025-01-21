package generator

import (
	"fmt"
	"strings"
)

type CmdMapValue struct {
	kv map[string]string
}

func NewCmdMapValue() *CmdMapValue {
	t := new(CmdMapValue)
	t.kv = make(map[string]string)
	return t
}

func (t *CmdMapValue) Set(val string) error {
	//*s = KeymapValue(strings.Split(val, ";"))
	list := strings.Split(val, ",")
	for _, v := range list {
		pair := strings.Split(v, "=")
		if len(pair) != 2 {
			panic("the para of kv must be key/value format, such as kv=key1=val2")
		}
		t.kv[pair[0]] = pair[1]
	}
	return nil
}

// flag为slice的默认值default is me,和return返回值没有关系
func (t *CmdMapValue) String() string {
	//*KeymapValue = KeymapValue(strings.Split("default is me", ","))
	return "It's none of my business"
}

func (t *CmdMapValue) Print() {
	fmt.Println("CmdMapValue:")
	for k, v := range t.kv {
		fmt.Printf("\tk:%s, v:%s\n", k, v)
	}
}

func (t *CmdMapValue) Get() map[string]string {
	return t.kv
}
