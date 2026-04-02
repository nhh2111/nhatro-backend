package controllers

import (
	"doAnHTTT_go/models"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllContractHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")

	// Lấy ownerID từ context
	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	resultData, err := services.GetAllContracts(ownerID, page, pageSize, search)

	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Lỗi lấy danh sách hợp đồng: "+err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, resultData)
}

func CreateContractHandler(c *gin.Context) {
	var contract models.Contract
	if err := c.ShouldBindJSON(&contract); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu hợp đồng không hợp lệ"})
		return
	}

	// Lấy ownerID từ context
	ownerIDVal, _ := c.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	if err := services.CreateContract(ownerID, &contract); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lập hợp đồng thành công! Phòng đã được chuyển sang trạng thái Đang Thuê."})
}

// API dùng để thanh lý / kết thúc hợp đồng
func TerminateContractHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Lấy ownerID từ context
	ownerIDVal, _ := c.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	if err := services.TerminateContract(ownerID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Thanh lý hợp đồng thành công! Phòng đã được trống."})
}
