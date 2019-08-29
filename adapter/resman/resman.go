package resman

import (
	"github.com/visionom/v-gae/adapter/resman/interfaces"
	"github.com/visionom/v-gae/adapter/resman/usecase"
)

func Init() interfaces.ComRepo {
	return usecase.NewComRepo()
}
