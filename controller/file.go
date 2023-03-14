package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"enigmacamp.com/fine_dms/config"
	"enigmacamp.com/fine_dms/middleware"
	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
	"enigmacamp.com/fine_dms/usecase"
	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileUsecase usecase.FileUsecase
}

func NewFileController(router *gin.RouterGroup, f usecase.FileUsecase, cfg *config.Secret) {
	fc := &FileController{fileUsecase: f}
	authMiddleware := middleware.ValidateToken(cfg.Key)

	router.GET("/files", authMiddleware, fc.getFilesByUserId)
	router.PUT("/files/:id", authMiddleware, fc.updateFile)
	router.DELETE("/files/:id", authMiddleware, fc.deleteFile)
	router.GET("/files/search", authMiddleware, fc.searchFiles)
}

func (fc *FileController) getFilesByUserId(ctx *gin.Context) {
	id, err := GetUserId(ctx)
	if err != nil {
		return
	}
	files, err := fc.fileUsecase.GetFilesByUserId(id)
	if err != nil {
		if errors.Is(err, usecase.ErrUsecaseNoData) {
			FailedJSONResponse(ctx, http.StatusNotFound, "no data")
		} else {
			FailedJSONResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	SuccessJSONResponse(ctx, http.StatusOK, "files retrieved successfully", files)
}

func (fc *FileController) updateFile(ctx *gin.Context) {
	userID, err := GetUserId(ctx)
	if err != nil {
		return
	}

	if err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var file model.File

	if err := ctx.BindJSON(&file); err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = fc.fileUsecase.UpdateFile(userID, file.Path, file.Ext)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidFileData) {
			FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		} else if errors.Is(err, usecase.ErrUsecaseNoData) {
			FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		} else {
			FailedJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		}

		return
	}
	SuccessJSONResponse(ctx, http.StatusOK,
		fmt.Sprintf("file with id = %d has been updated", userID),
		nil,
	)
}

func (fc *FileController) deleteFile(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	fileId, err := strconv.Atoi(ctx.Query("file_id"))
	if err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, "invalid file id")
		return
	}

	err = fc.fileUsecase.DeleteFile(userId, fileId)
	if err != nil {
		if errors.Is(err, usecase.ErrUsecaseNoData) {
			FailedJSONResponse(ctx, http.StatusNotFound, "no data")
		} else {
			FailedJSONResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	SuccessJSONResponse(ctx, http.StatusOK, "file deleted successfully", nil)
}

func (fc *FileController) searchFiles(ctx *gin.Context) {
	var files []model.File
	var err error

	idParam := ctx.Query("id")
	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			FailedJSONResponse(ctx, http.StatusBadRequest, "invalid id")
			return
		}
		files, err = fc.fileUsecase.SearchByUserId(id, ctx.Query("q"))
	} else {
		nameParam := ctx.Query("name")
		tagsParam := ctx.Query("tags")

		if nameParam != "" {
			files, err = fc.fileUsecase.SearchByName(nameParam)
		} else if tagsParam != "" {
			tags := strings.Split(tagsParam, ",")
			files, err = fc.fileUsecase.SearchByTags(tags)
		} else {
			SuccessJSONResponse(ctx, http.StatusOK, "Success", nil)
			return
		}
	}

	if err != nil {
		switch {
		case errors.Is(err, repo.ErrRepoNoData):
			FailedJSONResponse(ctx, http.StatusNotFound, "no data")
		case errors.Is(err, usecase.ErrInvalidUserID), errors.Is(err, usecase.ErrInvalidQuery), errors.Is(err, usecase.ErrUsecaseInvalidTag):
			FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		default:
			FailedJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	SuccessJSONResponse(ctx, http.StatusOK, "file search successfully", files)
}
