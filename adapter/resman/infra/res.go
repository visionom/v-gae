package db

import (
	"github.com/visionom/v-gae/adapter/mysql"
	"github.com/visionom/v-gae/adapter/mysql/domain"
)

type ResRdList []*ResRd

type ResRd struct {
	domain.ActiveRecord
}

func QueryRes(sql string, args ...interface{}) (ResRdList,
	error) {
	rds, err := mysql.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	resRds := make([]*ResRd, len(rds))
	for i, rd := range rds {
		resRds[i] = &ResRd{*rd}
	}
	return resRds, err
}

func QueryFirstRes(sql string, args ...interface{}) (*ResRd,
	error) {
	rd, err := mysql.QueryFirst(sql, args...)
	if err != nil {
		return nil, err
	}
	return &ResRd{*rd}, err

}

func (rd *ResRd) GetID() (string, error) {
	return rd.GetString("id")
}

func (rd *ResRd) GetName() (string, error) {
	return rd.GetString("name")
}

func (rd *ResRd) GetInfo() (string, error) {
	return rd.GetString("info")
}

func (rd *ResRd) SetID(v string) error {
	return rd.SetString("id", v)
}

func (rd *ResRd) SetName(v string) error {
	return rd.SetString("name", v)
}

func (rd *ResRd) SetInfo(v string) error {
	return rd.SetString("info", v)
}
