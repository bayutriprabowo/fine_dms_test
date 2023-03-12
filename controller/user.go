package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"enigmacamp.com/fine_dms/middleware"
	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/model/dto"
	"enigmacamp.com/fine_dms/usecase"
	"enigmacamp.com/fine_dms/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

func (self *UserController) HandleLogin(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiResponseFailed(err.Error()))
		return
	}

	userId, err := self.userUsecase.AuthenticateUser(loginData.Username, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewApiResponseFailed(err.Error()))
		return
	}

	secretKey, err := bcrypt.GenerateFromPassword([]byte(loginData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiResponseFailed(err.Error()))
		return
	}

	exp := time.Hour
	token, err := utils.GenerateToken(secretKey, userId, exp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiResponseFailed(err.Error()))
		return
	}

	responseData := gin.H{"token": token}
	c.JSON(http.StatusOK, dto.NewApiResponseSuccess("Login success", responseData))
}
