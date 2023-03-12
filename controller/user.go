package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"enigmacamp.com/fine_dms/middleware"
	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/usecase"
	"enigmacamp.com/fine_dms/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	userUsecase usecase.UserUsecase
}

func NewUserController(rg *gin.RouterGroup, u usecase.UserUsecase, secret []byte) {
	uc := UserController{u}

	authMiddleware := middleware.ValidateToken(secret)

	rg.POST("/login", uc.HandleLogin)
	auth := rg.Group("/")
	auth.Use(authMiddleware)
	{
		auth.GET("/user", uc.GetAll)
		auth.GET("/user/:id", uc.GetById)
		auth.PUT("/user/:id", uc.Edit)
		auth.DELETE("/user/:id", uc.Delete)
	}

	rg.POST("/user", uc.Add)
	rg.GET("/profile", uc.HandleProfile)
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
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	userId, err := self.userUsecase.AuthenticateUser(username, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	secretKey, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	exp := time.Hour
	token, err := utils.GenerateToken(secretKey, userId, exp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (self *UserController) HandleProfile(c *gin.Context) {
	userId, err := self.userUsecase.AuthenticateUser("Bearer", c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	profile, err := self.userUsecase.GetById(int(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}
