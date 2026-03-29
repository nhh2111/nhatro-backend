package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"
	"time"
)

func GetAllTask(page int, pageSize int, search string) (map[string]interface{}, error) {
	var taskList []models.Task
	var totalRecords int64

	query := config.DB.Model(&models.Task{}).Preload("House").Preload("Room")

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("status LIKE ? OR title LIKE ?", searchKeyword, searchKeyword)
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

func CreateNewTask(newTask *models.Task) error {
	if newTask.Status == "" {
		newTask.Status = "OPEN"
	}
	result := config.DB.Create(newTask)
	return result.Error
}

func UpdateTask(taskID uint, updatedData map[string]interface{}) error {
	var task models.Task

	errFind := config.DB.First(&task, taskID).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu nhiệm vụ cần sửa")
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

	return nil
}

func DeleteTask(taskID uint) error {
	var task models.Task
	result := config.DB.Delete(&task, taskID)
	if result.Error != nil {
		return errors.New("lỗi khi xóa nhiệm vụ")
	}
	if result.RowsAffected == 0 {
		return errors.New("không tìm thấy nhiệm vụ để xóa")
	}
	return nil
}
