package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase usecase.UserUsecase
}

func NewUserController(rg *gin.RouterGroup, u usecase.UserUsecase) {
	uc := UserController{u}

	rg.GET("/user", uc.GetAll)
	rg.GET("/user/:id", uc.GetById)
	rg.POST("/user", uc.Add)
	rg.PUT("/user/:id", uc.Edit)
	rg.DELETE("/user/:id", uc.Delete)
}

func (self *UserController) GetAll(ctx *gin.Context) {
	res, err := self.userUsecase.GetAll()
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

func (self *UserController) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	res, err := self.userUsecase.GetById(id)
	if err != nil {
		if err == usecase.ErrUsecaseNoData {
			FailedJSONResponse(ctx, http.StatusNotFound, "no data")
		} else {
			FailedJSONResponse(ctx, http.StatusInternalServerError,
				"Internal server error")
		}

		return
	}

	SuccessJSONResponse(ctx, http.StatusOK, "", res)
}

func (self *UserController) Add(ctx *gin.Context) {
	var user model.User

	if err := ctx.BindJSON(&user); err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, "Invalid Input")
		return
	}

	if err := self.userUsecase.Add(&user); err != nil {
		if err == usecase.ErrUsecaseInternal {
			FailedJSONResponse(ctx, http.StatusInternalServerError,
				"Internal server error",
			)
		} else {
			FailedJSONResponse(ctx, http.StatusBadRequest,
				err.Error(),
			)
		}
		return
	}

	SuccessJSONResponse(ctx, http.StatusCreated,
		fmt.Sprintf("user = %s successfully added", user.Username),
		nil,
	)
}

func (self *UserController) Edit(ctx *gin.Context) {
	var user model.User

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, "Invalid user id")
		return
	}

	user.ID = id

	if err := ctx.BindJSON(&user); err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, "Invalid Input")
		return
	}

	if err := self.userUsecase.Edit(&user); err != nil {
		if err == usecase.ErrUsecaseInternal {
			FailedJSONResponse(ctx, http.StatusInternalServerError,
				"Internal server error")
		} else {
			FailedJSONResponse(ctx, http.StatusBadRequest,
				err.Error())
		}
		return
	}

	SuccessJSONResponse(ctx, http.StatusOK,
		fmt.Sprintf("user with id = %d has been updated", id),
		nil,
	)
}

func (self *UserController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, "Invalid user id")
		return
	}

	err = self.userUsecase.Del(id)
	if err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	SuccessJSONResponse(ctx, http.StatusOK,
		fmt.Sprintf("user with id = %d has been deleted", id),
		nil,
	)
}
