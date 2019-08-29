package usecase

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/visionom/v-gae/adapter/mysql"
	"github.com/visionom/v-gae/adapter/resman/domain"
	"github.com/visionom/v-gae/adapter/resman/errors"
	db "github.com/visionom/v-gae/adapter/resman/infra"
	"github.com/visionom/v-gae/adapter/resman/interfaces"
)

type ComRepoImpl struct{}

func NewComRepo() interfaces.ComRepo {
	return &ComRepoImpl{}
}

func (*ComRepoImpl) NewRes(reses []domain.Res) error {
	for _, res := range reses {
		if res.ID == "" || res.Name == "" {
			fmt.Println("warning: id or name is empty")
			//TODO: add warning
			continue
		}
		if res.Info == nil {
			res.Info = json.RawMessage{'{', '}'}
		}

		if len(res.Tags) > 0 {
			valStr := ""
			var params []interface{}
			for i, tag := range res.Tags {
				if i > 0 {
					valStr += ","
				}
				valStr += "(?, ?, ?)"
				params = append(params, res.ID, tag.K, tag.V)
			}
			tx, err := mysql.BeginTx()
			if err != nil {
				return errors.Wrapf(err, "Fail to Start a transaction")
			}

			_, err = tx.Update("insert into res (id, name, info) value (?, ?, ?)", res.ID, res.Name, res.Info)
			if err != nil {
				return errors.Wrapf(err,
					"Fail to create resource")
			}

			_, err = tx.Update("insert into tags (res, k, v) values "+valStr,
				params...)
			if err != nil {
				return errors.Wrapf(err, "Fail to create resource when add the tags")
			}

			err = tx.Commit()
			if err != nil {
				return errors.Wrapf(err, "Fail to create resource when commit modify")
			}
		} else {
			_, err := mysql.Update("insert into res (id, name, "+
				"info) value (?, ?, ?)", res.ID, res.Name, res.Info)
			if err != nil {
				return errors.Wrapf(err, "Fail to create resource")
			}
		}
	}
	return nil
}

func (*ComRepoImpl) DelRes(ids []string) error {
	for _, id := range ids {
		if id == "" {
			fmt.Println("warning: id or name is empty")
			//TODO: add warning
			continue
		}
		tx, err := mysql.BeginTx()
		if err != nil {
			return errors.Wrapf(err, "Fail to Start a transaction")
		}

		_, err = tx.Update("delete from res where id=?", id)
		if err != nil {
			return errors.Wrapf(err,
				"Fail to delete resource")
		}

		_, err = tx.Update("delete from tags where res=?", id)
		if err != nil {
			return errors.Wrapf(err,
				"Fail to delete resource when delete the tags")
		}

		err = tx.Commit()
		if err != nil {
			return errors.Wrapf(err,
				"Fail to delete resource when commit modify")
		}
	}
	return nil
}

func (*ComRepoImpl) ModRes(res domain.Res) error {
	if res.ID == "" || res.Name == "" {
		return errors.New("Fail to modify res, params is empty")
	}
	if res.Info == nil {
		res.Info = json.RawMessage{'{', '}'}
	}

	if len(res.Tags) > 0 {
		valStr := ""
		var params []interface{}
		for i, tag := range res.Tags {
			if i > 0 {
				valStr += ","
			}
			valStr += "(?, ?, ?)"
			params = append(params, res.ID, tag.K, tag.V)
		}

		tx, err := mysql.BeginTx()
		if err != nil {
			return errors.Wrapf(err, "Fail to Start a transaction")
		}

		_, err = tx.Update("update res set name=?, info=? where id=?",
			res.Name,
			res.Info,
			res.ID)
		if err != nil {
			return errors.Wrapf(err,
				"Fail to modify resource")
		}

		_, err = tx.Update("delete from tags where res=?", res.ID)
		if err != nil {
			return errors.Wrapf(err,
				"Fail to modify resource when delete the tags")
		}

		_, err = tx.Update("insert ignore into tags (res, k, v) values "+valStr,
			params...)
		if err != nil {
			return errors.Wrapf(err,
				"Fail to modify resource when add the tags")
		}

		err = tx.Commit()
		if err != nil {
			return errors.Wrapf(err,
				"Fail to modify resource when commit modify")
		}
	} else {
		_, err := mysql.Update("update res set name=?, info=? where id=?", res.Name,
			res.Info,
			res.ID)
		return errors.Wrapf(err, "Fail to modify res")
	}
	return nil
}

func (*ComRepoImpl) AddTag(resID string, tags domain.Tags) (err error) {
	if resID == "" || len(tags) == 0 {
		return nil
	}
	valStr := ""
	var params []interface{}
	for i, tag := range tags {
		if i > 0 {
			valStr += ","
		}
		valStr += "(?, ?, ?)"
		params = append(params, resID, tag.K, tag.V)
	}

	_, err = mysql.Update("insert ignore into tags (res, k, v) values "+valStr,
		params...)
	return errors.Wrapf(err, "Fail to add tags")
}

func (*ComRepoImpl) RmTag(resID string, tags domain.Tags) (c int, err error) {
	if resID == "" || len(tags) == 0 {
		return 0, nil
	}
	valStr := "1 = 0 "
	var params []interface{}
	for _, tag := range tags {
		valStr += ""
		valStr += "or (res=? and k=? and v=? ) "
		params = append(params, resID, tag.K, tag.V)
	}
	rc, err := mysql.Update("delete from tags where "+valStr,
		params...)
	return int(rc), errors.Wrapf(err, "Fail to delete tags")
}

func (*ComRepoImpl) ChangeTags(resID string, oldTags domain.Tags, newTags domain.Tags) (err error) {

	if resID == "" || len(oldTags) == 0 || len(newTags) == 0 {
		return nil
	}

	oldValStr := "1 = 0 "
	var oldParams []interface{}
	for _, tag := range oldTags {
		oldValStr += ""
		oldValStr += "or (res=? and k=? and v=?) "
		oldParams = append(oldParams, resID, tag.K, tag.V)
	}

	newValStr := ""
	var newParams []interface{}
	for i, tag := range newTags {
		if i > 0 {
			newValStr += ","
		}
		newValStr += "(?, ?, ?)"
		newParams = append(newParams, resID, tag.K, tag.V)
	}

	tx, err := mysql.BeginTx()
	if err != nil {
		return errors.Wrapf(err, "Fail to Start a transaction")
	}
	c, err := tx.Update("delete from tags where "+oldValStr,
		oldParams...)
	if c != int64(len(oldTags)) {
		txErr := tx.RollBack()
		if txErr != nil {
			return errors.Wrapf(txErr, "Fail to change tags when delete tags, "+
				"not enough tags changed and rollback fail")
		}
		return errors.New("Fail to change tags when delete tags, " +
			"not enough tags changed")
	}

	_, err = tx.Update("insert into tags (res,k,v) values "+newValStr,
		newParams...)
	if err != nil {
		return errors.Wrapf(err,
			"Fail to change tags when add the tags")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrapf(err,
			"Fail to change tags when commit modify")
	}
	return nil
}

func (*ComRepoImpl) ChangeTag(resID string, key, old, new string) (err error) {
	if resID == "" || key == "" || old == "" || new == "" {
		return nil
	}

	_, err = mysql.Update("update tags set v=? where res=? and k=? and v=?",
		new, resID, key, old)
	return errors.Wrapf(err,
		"Fail to change tag when commit modify")
}

func (*ComRepoImpl) FindRes(tags domain.Tags, page, size int) (reses []domain.Res, err error) {
	resIDs := []string{}
	if size <= 0 {
		resIDs, err = findResIDsNoLimit(tags)
		if err != nil {
			return
		}
	} else {
		resIDs, err = findResIDsWithLimit(tags, page, size)
		if err != nil {
			return
		}
	}
	return findResById(resIDs)
}

func findResIDsWithLimit(tags domain.Tags, page int, size int) ([]string, error) {
	return findResIDs(tags, true, page, size)
}

func findResIDsNoLimit(tags domain.Tags) ([]string, error) {
	return findResIDs(tags, false, 0, 0)
}

func findResIDs(tags domain.Tags, limit bool, page int, size int) (resIDs []string,
	err error) {
	if len(tags) == 0 {
		return
	}

	var where string
	params := make([]interface{}, len(tags)*2)
	for i, tag := range tags {
		if i == 0 {
			where += "and (k=? and v=?) "
		} else {
			where += "or (k=? and v=?) "
		}
		params[2*i+0] = tag.K
		params[2*i+1] = tag.V
	}
	sql := fmt.Sprintf("select res "+
		"from tags "+
		"where 1=1 %s "+
		"group by res having count(res)=? order by res desc",
		where)
	params = append(params, len(tags))
	if limit {
		sql += " limit ?, ? "
		offset, rawCount := PageToLimit(page, size)
		params = append(params, offset, rawCount)
	}
	rds, err := db.QueryTag(sql, params...)
	if err != nil {
		return []string{},
			errors.Wrapf(err,
				"Fail to find res")
	}

	resIDs, err = getResIDs(rds)
	if err != nil {
		return []string{},
			errors.Wrapf(err,
				"Fail to find res when parse resIDs")
	}
	return
}

func findResById(resIDs []string) (reses []domain.Res, err error) {
	if len(resIDs) == 0 {
		return
	}
	where := ""
	params := make([]interface{}, len(resIDs))
	for i, id := range resIDs {
		if i > 0 {
			where += ", "
		}
		where += "?"
		params[i] = id
	}

	sql := fmt.Sprintf("select "+
		"res.id id, "+
		"res.name name, "+
		"res.info info, "+
		"group_concat(tags.k, ':', tags.v) tags "+
		"from res "+
		"left join tags on tags.res = res.id "+
		"where id in (%s) "+
		"group by id order by res.id desc", where)

	resRds, err := db.QueryRes(sql, params...)
	if err != nil {
		return []domain.Res{},
			errors.Wrapf(err,
				"Fail to find res when parse resIDs")
	}
	reses, err = getResList(resRds)
	err = errors.Wrapf(err, "Fail to find res when parse res")
	return
}

func getResList(resRdLists db.ResRdList) ([]domain.Res, error) {
	reses := make([]domain.Res, len(resRdLists))
	var err error
	for i, rd := range resRdLists {
		reses[i].ID, err = rd.GetID()
		if err != nil {
			return []domain.Res{}, err
		}

		reses[i].Name, err = rd.GetName()
		if err != nil {
			return []domain.Res{}, err
		}

		info, err := rd.GetInfo()
		if err != nil {
			return []domain.Res{}, err
		}
		reses[i].Info = json.RawMessage(info)

		tags, err := rd.GetString("tags")
		if err != nil {
			return []domain.Res{}, err
		}
		reses[i].Tags = domain.Tags{}
		for _, tag := range strings.Split(tags, ",") {
			reses[i].Tags = append(reses[i].Tags,
				newTagFS(tag))
		}
	}
	return reses, nil
}
func newTagFS(tag string) domain.Tag {
	ts := strings.Split(tag, ":")
	if len(ts) == 2 {
		return domain.Tag{K: ts[0], V: ts[1]}
	}
	return domain.Tag{K: tag, V: tag}
}

func getResIDs(rds db.TagRdList) ([]string, error) {
	resIDs := make([]string, len(rds))
	var err error
	for i, rd := range rds {
		resIDs[i], err = rd.GetRes()
		if err != nil {
			return []string{}, err
		}
	}
	return resIDs, err
}

func (*ComRepoImpl) FindResByID(ids []string) (reses []domain.Res, err error) {
	return findResById(ids)
}

func (*ComRepoImpl) Count(tags domain.Tags) (c int, err error) {
	if len(tags) == 0 {
		return
	}

	var where string
	params := make([]interface{}, len(tags)*2)
	for i, tag := range tags {
		if i == 0 {
			where += "and (k=? and v=?) "
		} else {
			where += "or (k=? and v=?) "
		}
		params[2*i+0] = tag.K
		params[2*i+1] = tag.V
	}

	sql := fmt.Sprintf("select count(a.res) c from (select res "+
		"from tags "+
		"where 1=1 %s "+
		"group by res having count(res)=? ) as a",
		where)

	params = append(params, len(tags))
	rd, err := mysql.QueryFirst(sql, params...)
	if err != nil {
		return -1, errors.Wrapf(err, "Fail to count res")
	}
	c, err = rd.GetInt("c")
	if err != nil {
		return -1, errors.Wrapf(err, "Fail to count res when get count")
	}
	return
}

func (*ComRepoImpl) CountGroupBy(tags domain.Tags,
	groupKey string) (cMap map[string]int, err error) {
	resIDs, err := findResIDsNoLimit(tags)
	cMap = make(map[string]int)
	if len(resIDs) <= 0 {
		return
	}
	var where string
	params := make([]interface{}, len(resIDs))
	for i, res := range resIDs {
		if i == 0 {
			where += "? "
		} else {
			where += ", ?"
		}
		params[i] = res
	}
	sql := fmt.Sprintf("select v, count(v) c from tags where res in ("+
		"%s) and k=? group by v", where)
	params = append(params, groupKey)
	rds, err := mysql.Query(sql, params...)
	if err != nil {
		err = errors.Wrap(err, "Fail to count")
	}

	for _, rd := range rds {
		c, err := rd.GetInt("c")
		if err != nil {
			err = errors.Wrap(err, "Fail to count when get count number")
		}
		v, err := rd.GetString("v")
		if err != nil {
			err = errors.Wrap(err, "Fail to count when get v value")
		}
		cMap[v] = c
	}
	return
}

func PageToLimit(page, size int) (int, int) {
	_size := 20
	_page := 1
	if size > 0 {
		_size = size
	}
	if page > 1 {
		_page = page
	}
	return (_page - 1) * _size, _size
}
