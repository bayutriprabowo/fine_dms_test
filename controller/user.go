package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	secret      *config.Secret
}

func NewUserController(router *gin.RouterGroup, u usecase.UserUsecase,
	cfg *config.Secret) {

	uc := UserController{u, cfg}
	authMiddleware := middleware.ValidateToken(cfg.Key)

	router.POST("/login", uc.HandleLogin)
	router.POST("/user", uc.Add)

	router.GET("/user", authMiddleware, uc.GetAll)
	router.GET("/profile", authMiddleware, uc.GetById)
	router.PUT("/user", authMiddleware, uc.Edit)
	router.DELETE("/user", authMiddleware, uc.Delete)
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
	id, err := self.getUserId(ctx)
	if err != nil {
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

	if err := ctx.BindJSON(&user); err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, "Invalid Input")
		return
	}

	id, err := self.getUserId(ctx)
	if err != nil {
		return
	}

	user.ID = id

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
	id, err := self.getUserId(ctx)
	if err != nil {
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

	userID, err := self.userUsecase.AuthenticateUser(loginRequest.Username,
		loginRequest.Password)
	if err != nil {
		FailedJSONResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := utils.GenerateToken(self.secret.Key, userID, self.secret.Exp)
	if err != nil {
		FailedJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// TODO: A proper expiration time (unix epoch)
	responseData := gin.H{"token": token, "expired": self.secret.Exp}
	SuccessJSONResponse(ctx, http.StatusOK, "Login success", responseData)
}

// private
func (self *UserController) getUserId(ctx *gin.Context) (int, error) {
	user_id, ok := ctx.Get("user_id")
	if !ok {
		FailedJSONResponse(ctx, http.StatusBadRequest, "invalid user id")
		return -1, errors.New("invalid user id")
	}

	id, err := strconv.Atoi(user_id.(string))
	if err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, "invalid user id")
		return -1, errors.New("invalid user id")
	}

	return id, nil
}
