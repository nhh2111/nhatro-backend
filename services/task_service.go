package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"
	"fmt"
	"time"
)

func GetAllTask(ownerID uint, page int, pageSize int, search string) (map[string]interface{}, error) {
	var taskList []models.Task
	var totalRecords int64

	query := config.DB.Model(&models.Task{}).Preload("House").Preload("Room").Where("owner_id = ?", ownerID)

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("(status LIKE ? OR title LIKE ?)", searchKeyword, searchKeyword)
	}
	query.Count(&totalRecords)

	pageCount := utils.GetPageCount(totalRecords, pageSize)
	offset := utils.GetOffset(page, pageSize)

	result := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&taskList)
	if result.Error != nil {
		return nil, result.Error
	}
	return map[string]interface{}{
		"recordCount": totalRecords,
		"pageCount":   pageCount,
		"currentPage": page,
		"pageSize":    pageSize,
		"records":     taskList,
	}, nil
}

func CreateNewTask(ownerID uint, newTask *models.Task) error {
	if newTask.Status == "" {
		newTask.Status = "OPEN"
	}
	newTask.OwnerID = ownerID

	result := config.DB.Create(newTask)
	return result.Error
}

func UpdateTask(ownerID uint, taskID uint, updatedData map[string]interface{}) error {
	var task models.Task

	errFind := config.DB.Where("id = ? AND owner_id = ?", taskID, ownerID).First(&task).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu nhiệm vụ hoặc bạn không có quyền sửa")
	}

	if status, ok := updatedData["status"]; ok {
		if status == "DONE" {
			updatedData["finished_at"] = time.Now()
		} else {
			updatedData["finished_at"] = nil
		}
	}

	errUpdate := config.DB.Model(&task).Updates(updatedData).Error
	if errUpdate != nil {
		return errors.New("lỗi khi cập nhật thông tin nhiệm vụ")
	}

	config.DB.First(&task, taskID)

	if task.Status == "DONE" && task.Cost > 0 {
		desc := fmt.Sprintf("Tự động chi trả bảo trì - Nhiệm vụ #%d: %s", task.ID, task.Title)

		var existingTx models.Transaction
		config.DB.Where("description = ?", desc).First(&existingTx)

		if existingTx.ID == 0 {
			newTx := models.Transaction{
				HouseID:         &task.HouseID,
				RoomID:          task.RoomID,
				Type:            "EXPENSE",
				Category:        "Sửa chữa & Bảo trì",
				Amount:          task.Cost,
				TransactionDate: time.Now(),
				PayerPayeeName:  task.Assignee,
				Description:     desc,
			}
			config.DB.Create(&newTx)
		}
	}

	return nil
}

func DeleteTask(ownerID uint, taskID uint) error {
	var task models.Task

	errFind := config.DB.Where("id = ? AND owner_id = ?", taskID, ownerID).First(&task).Error
	if errFind != nil {
		return errors.New("không tìm thấy nhiệm vụ hoặc bạn không có quyền xóa")
	}

	result := config.DB.Delete(&task)
	if result.Error != nil {
		return errors.New("lỗi khi xóa nhiệm vụ")
	}
	return nil
}
