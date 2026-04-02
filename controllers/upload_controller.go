package controllers

import (
	"doAnHTTT_go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadImageHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Vui lòng chọn một file ảnh")
		return
	}

	folder := c.PostForm("folder")
	if folder == "" {
		folder = "general"
	}

	oldImageURL := c.PostForm("old_url")

	fileURL, errUpload := utils.ProcessImageUpload(c, file, folder, oldImageURL)
	if errUpload != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, errUpload.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, map[string]string{"url": fileURL})
}

func UploadMultipleImagesHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Dữ liệu form không hợp lệ")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Vui lòng chọn ít nhất 1 ảnh")
		return
	}

	folder := c.PostForm("folder")
	if folder == "" {
		folder = "gallery"
	}

	fileURLs, errUpload := utils.ProcessMultipleImagesUpload(c, files, folder)
	if errUpload != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, errUpload.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{"urls": fileURLs})
}

func DeleteImageHandler(c *gin.Context) {
	var req struct {
		Folder  string `json:"folder" binding:"required"`
		FileURL string `json:"file_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Dữ liệu không hợp lệ")
		return
	}

	utils.DeleteOldFile(req.Folder, req.FileURL)

	utils.SuccessResponse(c, http.StatusOK, map[string]string{"message": "Đã xóa file ảnh vật lý"})
}
