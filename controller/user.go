package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/usecase"
	"github.com/dgrijalva/jwt-go"
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
	rg.POST("/login", uc.HandleLogin)
	rg.GET("/profile", validateToken(), uc.HandleProfile)

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

// -------------------------------------------memanggil jwt-----------------------------------------------

var signingKey = []byte("secret-key")

const tokenExpiration = time.Hour * 24

func generateToken(userId int64) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(tokenExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func validateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println(authHeader)
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return signingKey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := claims["userId"].(float64)
			c.Set("userId", int64(userId))
			c.Next() // ini keyword untuk menlanjutkan ke middleware selanjutnya, atau ke endpointnya
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
	}
}

// --------------------------------loginnya----------------------------------

func (self *UserController) HandleLogin(c *gin.Context) {
	// Melakukan proses autentikasi user dan mendapatkan ID pengguna
	userId, err := self.userUsecase.AuthenticateUser(c.Request.FormValue("username"), c.Request.FormValue("password"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Generate token JWT dan mengirimkannya sebagai response
	token, err := generateToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (self *UserController) HandleProfile(c *gin.Context) {
	// Mendapatkan ID pengguna dari JWT yang valid
	userId, err := self.userUsecase.GetUserIdFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Menampilkan data profil pengguna
	profile, err := self.userUsecase.GetById(int(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}
