package model

import (
	"github.com/mjiulee/lego"
	"time"
)

//go:generate legogen -type=Example

func init() {
	lego.AddBeanToSynList(new(DeliveryCron))
}

type Example struct {
	Id    int64     `xorm:"pk autoincr"`      // ["id","hide","hidden"]
	IfDel int       `xorm:"TINYINT(1)"`       // ["deleted","hide","hidden"]
	Cdate time.Time `xorm:"DateTime created"` // ["create","show","datetime"]
	Udate time.Time `xorm:"DateTime updated"` // ["update","hide","datetime"]
	// add your custom field here
}
