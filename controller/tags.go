package controller

import (
	"net/http"

	"enigmacamp.com/fine_dms/usecase"
	"github.com/gin-gonic/gin"
)

type TagsController struct {
	tagsUsecase usecase.TagsUsecase
}

func NewTagsController(rg *gin.RouterGroup, u usecase.TagsUsecase) {
	uc := TagsController{u}

	rg.GET("/tags", uc.GetAll)
}

func (self *TagsController) GetAll(ctx *gin.Context) {
	res, err := self.tagsUsecase.GetAll()
	if err != nil {
		if err == usecase.ErrUsecaseNoData {
			FailedJSONResponse(ctx, http.StatusNotFound, "no data")
		} else {
			FailedJSONResponse(ctx, http.StatusInternalServerError,
				"internal server error")
		}

		return
	}

	SuccessJSONResponse(ctx, http.StatusOK, "", res)
}
