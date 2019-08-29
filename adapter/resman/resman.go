package resman

import (
	"github.com/visionom/v-gae/adapter/resman/ifs"
	"github.com/visionom/v-gae/adapter/resman/usecase"
)

func Init() ifs.ComRepo {
	return usecase.NewComRepo()
}
