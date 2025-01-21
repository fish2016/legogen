package model

const K_TMPL_MODEL_CODE = `
/* -----------------
* brief 
* 1. this is code gen by tools 
* 2. array like [a,b,c] after each field, define the ui html tag in view, 
		a for title, 
		b for show or not ,  
		c for element html tag , currently support : check、radio、select、input、textarea、upload
* 3. XXX
*/

package %s

import "github.com/mjiulee/lego"

func init() {
	lego.AddBeanToSynList(new(%s))
}

/* table: "tb_%s"
 * -----------------
*/
type %s struct {
	Id        int64  ###xorm:"pk"###                  // ["id","hide","hidden"]
	IfDel     int ###xorm:"default 0"### // ["deleted","hide","hidden"]
	Cdate     string ###xorm:"DateTime created"###    // ["create","show","datetime"]
	Udate     string ###xorm:"DateTime updated"###    // ["update","show","datetime"]
    // add your custom field here
}`
