/* -----------------
* brief
* 1. this is code gen by tools
* 2. for add/update/delete/query the data save in sql database
 */

package {{.StructName}}




import (
    "github.com/mjiulee/lego"
    "github.com/go-xorm/xorm"
    //"github.com/mjiulee/lego/utils"
)

/* 根据id获取记录
 * -----------------
 */
func (t *{{.StructName}} ) RecodeById(id int64) *{{.StructName}}  {
	item := new({{.StructName}} )
	lego.GetDBEngine().ID(id).Get(item)
	if item.Id <= 0 {
		return nil
	}
	return item
}

/* 添加
 * -----------------
 */
func (t *{{.StructName}} ) AddRecode(item2add *{{.StructName}} ) bool {
	item2add.Id = lego.UUID()
	_, err := lego.GetDBEngine().Insert(item2add)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}

func (t *{{.StructName}} ) AddRecodeWithsession(session*xorm.Session,item2add *{{.StructName}} ) bool {
	item2add.Id = lego.UUID()
	_, err := session.Insert(item2add)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}

/* 删除(硬删除)*/
func (t *{{.StructName}} ) DelRecodeById(id int64) bool {
	item := new({{.StructName}} )
	_, err := lego.GetDBEngine().ID(id).Delete(item)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}

func (t *{{.StructName}} ) DelRecodeByIdWithSession(session*xorm.Session,id int64) bool {
	item := new({{.StructName}} )
	_, err := session.ID(id).Delete(item)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}

/* 更新
 * -----------------*/
func (t *{{.StructName}} ) UpdateRecode(rc *{{.StructName}} ) bool {
	_, err := lego.GetDBEngine().Id(rc.Id).Update(rc)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}

func (t *{{.StructName}} ) UpdateRecodeWithSession(session *xorm.Session,rc *{{.StructName}} ) bool {
	_, err := session.ID(rc.Id).Update(rc)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}