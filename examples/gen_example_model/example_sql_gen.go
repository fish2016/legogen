/* -----------------
* brief
* 1. this is code gen by tools
* 2. for add/update/delete/query the data save in sql database
 */

package Example

import (
	"github.com/go-xorm/xorm"
	"github.com/mjiulee/lego"
	//"github.com/mjiulee/lego/utils"
)

/* 根据id获取记录
 * -----------------
 */
func (t *Example) RecodeById(id int64) *Example {
	item := new(Example)
	lego.GetDBEngine().ID(id).Get(item)
	if item.Id <= 0 {
		return nil
	}
	return item
}

/* 添加
 * -----------------
 */
func (t *Example) AddRecode(item2add *Example) bool {
	item2add.Id = lego.UUID()
	_, err := lego.GetDBEngine().Insert(item2add)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}

func (t *Example) AddRecodeWithsession(session *xorm.Session, item2add *Example) bool {
	item2add.Id = lego.UUID()
	_, err := session.Insert(item2add)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}

/* 删除(硬删除)*/
func (t *Example) DelRecodeById(id int64) bool {
	item := new(Example)
	_, err := lego.GetDBEngine().ID(id).Delete(item)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}

func (t *Example) DelRecodeByIdWithSession(session *xorm.Session, id int64) bool {
	item := new(Example)
	_, err := session.ID(id).Delete(item)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}

/* 更新
 * -----------------*/
func (t *Example) UpdateRecode(rc *Example) bool {
	_, err := lego.GetDBEngine().Id(rc.Id).Update(rc)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}

func (t *Example) UpdateRecodeWithSession(session *xorm.Session, rc *Example) bool {
	_, err := session.ID(rc.Id).Update(rc)
	if err != nil {
		lego.LogError(err)
		return false
	}
	return true
}
