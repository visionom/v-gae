package db

import (
	"github.com/visionom/v-gae/adapter/mysql"
	"github.com/visionom/v-gae/adapter/mysql/domain"
)

type TagRdList []*TagRd

type TagRd struct {
	domain.ActiveRecord
}

func QueryTag(sql string, args ...interface{}) (TagRdList,
	error) {
	rds, err := mysql.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	resRds := make([]*TagRd, len(rds))
	for i, rd := range rds {
		resRds[i] = &TagRd{*rd}
	}
	return resRds, err
}

func QueryFirstTag(sql string, args ...interface{}) (*TagRd,
	error) {
	rd, err := mysql.QueryFirst(sql, args...)
	if err != nil {
		return nil, err
	}
	return &TagRd{*rd}, err

}

func (rd *TagRd) GetID() (string, error) {
	return rd.GetString("id")
}

func (rd *TagRd) GetTag() (string, error) {
	return rd.GetString("tag")
}

func (rd *TagRd) GetRes() (string, error) {
	return rd.GetString("res")
}

func (rd *TagRd) SetID(v string) error {
	return rd.SetString("id", v)
}

func (rd *TagRd) SetTag(v string) error {
	return rd.SetString("tag", v)
}

func (rd *TagRd) SetRes(v string) error {
	return rd.SetString("res", v)
}
