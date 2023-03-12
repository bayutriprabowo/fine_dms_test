package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"enigmacamp.com/fine_dms/config"
	"enigmacamp.com/fine_dms/middleware"
	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/model/dto"
	"enigmacamp.com/fine_dms/usecase"
	"enigmacamp.com/fine_dms/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase usecase.UserUsecase
}

func NewUserController(router *gin.Engine, u usecase.UserUsecase, secret []byte) {
	uc := UserController{u}

	authMiddleware := middleware.ValidateToken(secret)

	router.POST("/login", uc.HandleLogin)
	router.POST("/user", uc.Add)

	router.GET("/user", authMiddleware, uc.GetAll)
	router.GET("/profile", authMiddleware, uc.GetById)
	router.PUT("/user/:id", authMiddleware, uc.Edit)
	router.DELETE("/user/:id", authMiddleware, uc.Delete)
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

func (self *UserController) HandleLogin(ctx *gin.Context) {
	var loginRequest dto.ApiloginRequest
	if err := ctx.BindJSON(&loginRequest); err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	secretKey := []byte(config.NewAppConfig().Secret.Key)
	userID, err := self.userUsecase.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		FailedJSONResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	exp := config.NewAppConfig().Secret.Exp
	token, err := utils.GenerateToken(secretKey, userID, exp)
	if err != nil {
		FailedJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	expired := time.Now().Add(exp).Unix()
	responseData := gin.H{"token": token, "expired": expired - time.Now().Unix()}
	SuccessJSONResponse(ctx, http.StatusOK, "Login success", responseData)
}
