package domain

import (
	"encoding/json"
	"strings"
)

// RawRes is a single Resource element query from table `res`
type RawRes struct {
	ID   string          `json:"id,res_id"`
	Name string          `json:"name"`
	Info json.RawMessage `json:"info"`
}

// Res is a single Resource element with tag list
type Res struct {
	RawRes
	Tags Tags `json:"tags"`
}

// Tag is a Pair of Key-Value element query from table `tag`
type Tag struct {
	K string `json:"key"`
	V string `json:"val"`
}

// Tags is a List of TagPair element
type Tags = []Tag

// NewRes ...
func NewRes(id, name string, info json.RawMessage) Res {
	return Res{RawRes{id, name, info}, Tags{}}
}

// AddTags ...
func (r *Res) AddTags(tags ...Tag) {
	r.Tags = append(r.Tags, tags...)
}

// GetRes ...
func (r *Res) GetRes() Res {
	return *r
}

// GetTags ...
func (r *Res) GetTags() Tags {
	return r.Tags
}

// GetTag ...
func (r *Res) GetTag(k string) string {
	for _, v := range r.Tags {
		if v.K == k {
			return v.V
		}
	}
	return ""
}

// GetTagList ...
func (r *Res) GetTagList(k string) []string {
	l := []string{}
	for _, v := range r.Tags {
		if v.K == k {
			l = append(l, v.V)
		}
	}
	return l
}

// NewTagFS is New Tag from string
func NewTagFS(tag string) Tag {
	ts := strings.Split(tag, ":")
	if len(ts) == 2 {
		return Tag{K: ts[0], V: ts[1]}
	}
	return Tag{K: tag, V: tag}
}

// JoinTags ...
func JoinTags(tags ...Tag) Tags {
	return tags
}

// NewTag ...
func NewTag(key, val string) Tag {
	return Tag{K: key, V: val}
}
