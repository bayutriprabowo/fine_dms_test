package controller

import (
	"net/http"
	"strconv"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	router  *gin.RouterGroup
	usecase usecase.UserUsecase
}

func (c *UserController) GetAllUser(ctx *gin.Context) {
	res := c.usecase.GetAllUser()

	ctx.JSON(http.StatusOK, res)
}

func (c *UserController) AddUser(ctx *gin.Context) {
	var user model.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid input")
		return
	}

	res := c.usecase.AddUser(&user)
	ctx.JSON(http.StatusCreated, res)
}

func (c *UserController) EditUser(ctx *gin.Context) {
	var user model.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid input")
		return
	}

	res := c.usecase.EditUser(&user)
	ctx.JSON(http.StatusOK, res)
}

func (c *UserController) RemoveUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid user ID")
		return
	}

	res := c.usecase.RemoveUser(id)
	ctx.JSON(http.StatusOK, res)
}

func NewUserController(r *gin.RouterGroup, u usecase.UserUsecase) *UserController {
	controller := UserController{
		router:  r,
		usecase: u,
	}

	return &controller
}
