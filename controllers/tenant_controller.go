package controllers

import (
	"doAnHTTT_go/models"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllTenantHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	resultData, err := services.GetAllTenant(ownerID, page, pageSize, search)
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Không thể lấy danh sách: "+err.Error())
		return
	}
	utils.SuccessResponse(ginContext, http.StatusOK, resultData)
}

func CreateTenantHandler(ginContext *gin.Context) {
	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	// 1. ĐỌC DỮ LIỆU TỪ FORM-DATA THAY VÌ JSON
	motorbikeCount, _ := strconv.Atoi(ginContext.PostForm("motorbike_count"))
	carCount, _ := strconv.Atoi(ginContext.PostForm("car_count"))

	newTenant := models.Tenant{
		FullName:       ginContext.PostForm("full_name"),
		CCCD:           ginContext.PostForm("cccd"),
		Phone:          ginContext.PostForm("phone"),
		Dob:            ginContext.PostForm("dob"),
		Gender:         ginContext.PostForm("gender"),
		Address:        ginContext.PostForm("address"),
		LicensePlates:  ginContext.PostForm("license_plates"),
		MotorbikeCount: motorbikeCount,
		CarCount:       carCount,
		OwnerID:        ownerID,
	}

	// 2. XỬ LÝ UPLOAD ẢNH (Nếu có)
	file, errFile := ginContext.FormFile("image")
	if errFile == nil {
		// Gọi hàm từ upload_utils.go của bạn
		imageURL, errUpload := utils.ProcessImageUpload(ginContext, file, "tenants", "")
		if errUpload == nil {
			newTenant.ImageUrl = imageURL
		}
	}

	// 3. LƯU VÀO DATABASE
	errCreate := services.CreateNewTenant(&newTenant)
	if errCreate != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, errCreate.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusCreated, newTenant)
}

func UpdateTenantHandler(ginContext *gin.Context) {
	tenantID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID khách hàng không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	// 1. LẤY CÁC TRƯỜNG DỮ LIỆU MUỐN CẬP NHẬT
	updateData := make(map[string]interface{})

	textFields := []string{"full_name", "cccd", "phone", "dob", "gender", "address", "license_plates"}
	for _, field := range textFields {
		if val := ginContext.PostForm(field); val != "" {
			updateData[field] = val
		}
	}

	if val := ginContext.PostForm("motorbike_count"); val != "" {
		mCount, _ := strconv.Atoi(val)
		updateData["motorbike_count"] = mCount
	}
	if val := ginContext.PostForm("car_count"); val != "" {
		cCount, _ := strconv.Atoi(val)
		updateData["car_count"] = cCount
	}

	// 2. XỬ LÝ ẢNH MỚI NẾU CÓ THAY ĐỔI
	file, errFile := ginContext.FormFile("image")
	if errFile == nil {
		// old_image_url được Angular gửi lên để Go xóa file cũ đi cho nhẹ server
		oldImageURL := ginContext.PostForm("old_image_url")
		imageURL, errUpload := utils.ProcessImageUpload(ginContext, file, "tenants", oldImageURL)
		if errUpload == nil {
			updateData["image_url"] = imageURL
		}
	}

	// 3. GỌI SERVICE ĐỂ LƯU
	errService := services.UpdateTenant(ownerID, tenantID, updateData)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Cập nhật khách hàng thành công"})
}

func DeleteTenantHandler(ginContext *gin.Context) {
	tenantID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID khách hàng không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errService := services.DeleteTenant(ownerID, tenantID)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Xóa khách hàng thành công"})
}
