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

	resultData, err := services.GetAllTenant(page, pageSize, search)
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Không thể lấy danh sách: "+err.Error())
		return
	}
	utils.SuccessResponse(ginContext, http.StatusOK, resultData)
}

func CreateTenantHandler(ginContext *gin.Context) {
	var newTenant models.Tenant

	errBind := ginContext.ShouldBindJSON(&newTenant)
	if errBind != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu đầu vào không hợp lệ")
		return
	}

	errCreate := services.CreateNewTenant(&newTenant)
	if errCreate != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, errCreate.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusCreated, newTenant)
}

func UpdateTenantHandler(ginContext *gin.Context) {
	tenantID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID khách hàng không hợp lệ"})
		return
	}

	var updateData map[string]interface{}
	if errBind := ginContext.ShouldBindJSON(&updateData); errBind != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu cập nhật không hợp lệ"})
		return
	}

	errService := services.UpdateTenant(uint(tenantID), updateData)
	if errService != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Cập nhật khách hàng thành công"})
}

func DeleteTenantHandler(ginContext *gin.Context) {
	tenantID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID khách hàng không hợp lệ"})
		return
	}

	errService := services.DeleteTenant(uint(tenantID))
	if errService != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Xóa khách hàng thành công"})
}
