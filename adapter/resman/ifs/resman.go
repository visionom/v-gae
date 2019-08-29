package ifs

import "github.com/visionom/v-gae/adapter/resman/domain"

// ComRepo is Resource Repo interface
type ComRepo interface {
	NewRes(reses []domain.Res) (err error)

	DelRes(ids []string) (err error)

	ModRes(res domain.Res) (err error)

	AddTag(resID string, tags domain.Tags) (err error)

	RmTag(resID string, tags domain.Tags) (c int, err error)

	ChangeTags(resID string, oldTags domain.Tags, newTags domain.Tags) (err error)

	ChangeTag(resID string, key, old, new string) (err error)

	FindRes(tags domain.Tags, page, size int) (reses []domain.Res, err error)

	FindResByID(ids []string) (reses []domain.Res, err error)

	Count(tags domain.Tags) (c int, err error)

	CountGroupBy(tags domain.Tags, groupKey string) (cMap map[string]int,
		err error)
}
