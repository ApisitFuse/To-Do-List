package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"to-do-list-app/database"
	"to-do-list-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTodos(c *gin.Context) {
	var todos []models.Todo
	result := database.DB.Order("display_order asc").Find(&todos)
	if result.Error != nil {
		log.Printf("Error fetching todos: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve todos"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func GetTrashed(c *gin.Context) {
	var todos []models.Todo
	if err := database.DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&todos).Error; err != nil {
		log.Printf("Error fetching trashed todos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve trashed items"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func CreateTodo(c *gin.Context) {

	log.Printf("create todo")
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func UpdateTodo(c *gin.Context) {

	id := c.Param("id")
	log.Printf("update todo %s", id)
	var todo models.Todo
	if err := database.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	var input models.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// database.DB.Model(&todo).Updates(input)
	if err := database.DB.Model(&todo).Select("Title", "Completed").Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func UpdateOrder(c *gin.Context) {

	var orderedTodos models.OrderObject

	if err := c.ShouldBindJSON(&orderedTodos); err != nil {
		log.Printf("Error binding JSON for order update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// OrId := orderedTodos.OrId
	// newOr := orderedTodos.NewIndex
	// oldOr := orderedTodos.OldIndex

	targetTodoID := orderedTodos.OrId
	newDisplayOrder := orderedTodos.NewIndex
	oldDisplayOrder := orderedTodos.OldIndex

	log.Printf("orderedTodos: %v", orderedTodos)
	log.Printf("orderedTodos ID: %v", targetTodoID)
	log.Printf("orderedTodos new: %v", newDisplayOrder)
	log.Printf("orderedTodos old: %v", oldDisplayOrder)

	if newDisplayOrder == oldDisplayOrder {
		log.Printf("Order unchanged for Todo ID %d.", targetTodoID)
		c.JSON(http.StatusOK, gin.H{"message": "Todo order unchanged"})
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		log.Printf("Failed to begin transaction: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction"})
		return
	}

	// --- BEGIN DEBUG LOGGING ---
	var allTodosBeforeUpdate []models.Todo
	if err := tx.Order("display_order asc").Find(&allTodosBeforeUpdate).Error; err != nil {
		tx.Rollback()
		log.Printf("DEBUG: Failed to fetch todos for logging (before any update): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data for debugging"})
		return
	}
	log.Printf("DEBUG: Todos BEFORE any update (TargetID: %d, OldOrder: %d, NewOrder: %d):", targetTodoID, oldDisplayOrder, newDisplayOrder)
	for _, t := range allTodosBeforeUpdate {
		log.Printf("DEBUG: ID: %d, Title: '%s', DisplayOrder: %d", t.ID, t.Title, t.DisplayOrder)
	}
	// --- END DEBUG LOGGING ---

	// 1. Update the DisplayOrder of the target Todo item FIRST.
	// We use UpdateColumn here to avoid GORM hooks if any, and to be explicit.
	updateResult := tx.Model(&models.Todo{}).Where("id = ?", targetTodoID).UpdateColumn("display_order", newDisplayOrder)
	if updateResult.Error != nil {
		tx.Rollback()
		log.Printf("Failed to update target todo ID %d to DisplayOrder %d: %v", targetTodoID, newDisplayOrder, updateResult.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update order for todo ID %d: %s", targetTodoID, updateResult.Error.Error())})
		return
	}
	// Check if the target item was actually found and updated.
	if updateResult.RowsAffected == 0 {
		tx.Rollback()
		log.Printf("Target todo ID %d not found for order update.", targetTodoID)
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Todo with ID %d not found for order update", targetTodoID)})
		return
	}
	log.Printf("Successfully updated target Todo ID %d to DisplayOrder %d. RowsAffected: %d", targetTodoID, newDisplayOrder, updateResult.RowsAffected)

	// 2. Adjust DisplayOrder of other affected Todo items.
	var errShift error
	var shiftResult *gorm.DB
	var queryCondition string
	var queryParams []interface{}

	if oldDisplayOrder > newDisplayOrder { // Item moved UP (e.g., from 5 to 2). Items between new (inclusive) and old (exclusive) need to be incremented.
		log.Printf("DEBUG: Item ID %d moved UP. Shifting items whose current DisplayOrder is >= %d AND < %d (excluding ID %d itself). They will be incremented.", targetTodoID, newDisplayOrder, oldDisplayOrder, targetTodoID)
		queryCondition = "display_order >= ? AND display_order < ? AND id != ?"
		queryParams = []interface{}{newDisplayOrder, oldDisplayOrder, targetTodoID}
		shiftResult = tx.Model(&models.Todo{}).
			Where(queryCondition, queryParams...).
			UpdateColumn("display_order", gorm.Expr("display_order + 1"))
		errShift = shiftResult.Error
	} else { // Item moved DOWN (e.g., from 2 to 5). Items between old (exclusive) and new (inclusive) need to be decremented.
		log.Printf("DEBUG: Item ID %d moved DOWN. Shifting items whose current DisplayOrder is > %d AND <= %d (excluding ID %d itself). They will be decremented.", targetTodoID, oldDisplayOrder, newDisplayOrder, targetTodoID)
		queryCondition = "display_order > ? AND display_order <= ? AND id != ?"
		queryParams = []interface{}{oldDisplayOrder, newDisplayOrder, targetTodoID}
		shiftResult = tx.Model(&models.Todo{}).
			Where(queryCondition, queryParams...).
			UpdateColumn("display_order", gorm.Expr("display_order - 1"))
		errShift = shiftResult.Error
	}

	if errShift != nil {
		tx.Rollback()
		log.Printf("Error shifting other items for Todo ID %d: %v. Query: %s, Params: %v", targetTodoID, errShift, queryCondition, queryParams)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subsequent item orders: " + errShift.Error()})
		return
	}

	// --- BEGIN DEBUG LOGGING ---
	var allTodosAfterUpdate []models.Todo
	if err := tx.Order("display_order asc").Find(&allTodosAfterUpdate).Error; err == nil {
		log.Printf("DEBUG: Todos AFTER all updates, BEFORE commit (TargetID: %d, OldOrder: %d, NewOrder: %d):", targetTodoID, oldDisplayOrder, newDisplayOrder)
		for _, t := range allTodosAfterUpdate {
			log.Printf("DEBUG: ID: %d, Title: '%s', DisplayOrder: %d", t.ID, t.Title, t.DisplayOrder)
		}
	} else {
		log.Printf("DEBUG: Failed to fetch todos for logging (after update, before commit): %v", err)
	}
	// --- END DEBUG LOGGING ---

	if err := tx.Commit().Error; err != nil {
		// If errShift caused a rollback, tx.Commit() will also error.
		// No need for an explicit tx.Rollback() here if it was already done,
		// but it's safe to call (GORM handles it).
		log.Printf("Transaction commit error for order update: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit order updates: " + err.Error()})
		return
	}

	// var todos []models.Todo
	// if err := database.DB.Find(&todos).Error; err != nil {
	// 	log.Printf("Failed to get all todos: %v", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all todos"})
	// 	return
	// }

	// if (oldOr - newOr) > 0 {

	// 	tx := database.DB.Begin()
	// 	if tx.Error != nil {
	// 		log.Printf("Failed to begin transaction: %v", tx.Error)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction"})
	// 		return
	// 	}

	// 	for index := newOr; index <= oldOr; index++ {
	// 		nextUpdate := tx.Model(&models.Todo{}).Where("id = ?", index).Update("display_order", gorm.Expr("display_order + 1"))
	// 		if nextUpdate.Error != nil {
	// 			tx.Rollback()
	// 			log.Printf("Failed to update order for todo ID %d: %v", index, nextUpdate.Error)
	// 			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update subsequent order (intended ID/index %d): %s", index, nextUpdate.Error.Error())})
	// 			return
	// 		}
	// 	}

	// 	firstUpdate := tx.Model(&models.Todo{}).Where("id = ?", OrId).Update("display_order", newOr)
	// 	log.Printf("first update: %v", firstUpdate)
	// 	if firstUpdate.Error != nil {
	// 		tx.Rollback()
	// 		log.Printf("Failed to update order for todo ID %d: %v", OrId, firstUpdate.Error)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update order for todo ID %d: %s", OrId, firstUpdate.Error.Error())})
	// 		return
	// 	}

	// 	if err := tx.Commit().Error; err != nil {
	// 		log.Printf("Transaction commit error for order update: %v", err)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit order updates: " + err.Error()})
	// 		return
	// 	}

	// } else {
	// 	tx := database.DB.Begin()
	// 	if tx.Error != nil {
	// 		log.Printf("Failed to begin transaction: %v", tx.Error)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction"})
	// 		return
	// 	}

	// 	firstUpdate := tx.Model(&models.Todo{}).Where("id = ?", OrId).Update("display_order", newOr)
	// 	if firstUpdate.Error != nil {
	// 		tx.Rollback()
	// 		log.Printf("Failed to update order for todo ID %d: %v", OrId, firstUpdate.Error)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update order for todo ID %d: %s", OrId, firstUpdate.Error.Error())})
	// 		return
	// 	}

	// 	for index := newOr; index >= 1; index-- {
	// 		nextUpdate := tx.Model(&models.Todo{}).Where("id = ?", index).Update("display_order", gorm.Expr("display_order - 1"))
	// 		if nextUpdate.Error != nil {
	// 			tx.Rollback()
	// 			log.Printf("Failed to update order for todo ID %d: %v", index, nextUpdate.Error)
	// 			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update subsequent order (intended ID/index %d): %s", index, nextUpdate.Error.Error())})
	// 			return
	// 		}
	// 	}

	// 	if err := tx.Commit().Error; err != nil {
	// 		log.Printf("Transaction commit error for order update: %v", err)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit order updates: " + err.Error()})
	// 		return
	// 	}
	// }

	c.JSON(http.StatusOK, gin.H{"message": "Todo order updated successfully"})

}

func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo
	if err := database.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	database.DB.Delete(&todo)

	// แบบที่ลบแล้วจะลด order ที่มากกว่ามันลงมาด้วย---------------------------------------------------------
	// tx := database.DB.Begin()
	// if tx.Error != nil {
	// 	log.Printf("Failed to begin transaction for delete: %v", tx.Error)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction"})
	// 	return
	// }

	// var todoToDelete models.Todo
	// // 1. ค้นหา Todo ที่จะลบเพื่อเอา displayOrder เดิม
	// if err := tx.First(&todoToDelete, id).Error; err != nil {
	// 	tx.Rollback()
	// 	if err == gorm.ErrRecordNotFound {
	// 		log.Printf("Todo ID %s not found for deletion.", id)
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
	// 	} else {
	// 		log.Printf("Error fetching todo ID %s for deletion: %v", id, err)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding todo to delete"})
	// 	}
	// 	return
	// }

	// // 2. ลบ Todo ที่ต้องการ
	// if err := tx.Delete(&todoToDelete).Error; err != nil {
	// 	tx.Rollback()
	// 	log.Printf("Error deleting todo ID %s: %v", id, err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
	// 	return
	// }

	// // 3. ปรับ display_order ของรายการที่มี display_order มากกว่ารายการที่ถูกลบ
	// // โดยลดค่า display_order ลง 1
	// if err := tx.Model(&models.Todo{}).Where("display_order > ?", todoToDelete.DisplayOrder).UpdateColumn("display_order", gorm.Expr("display_order - 1")).Error; err != nil {
	// 	tx.Rollback()
	// 	log.Printf("Error updating display_order after deleting todo ID %s: %v", id, err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order of subsequent todos"})
	// 	return
	// }

	// // Commit Transaction
	// if err := tx.Commit().Error; err != nil {
	// 	log.Printf("Transaction commit error for delete: %v", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit delete operation"})
	// 	return
	// }
	c.JSON(http.StatusOK, gin.H{"message": "todo deleted"})
}

func PermanentlyDeleteTodo(c *gin.Context) {
	id := c.Param("id")
	var todoToDeleteInfo models.Todo

	if err := database.DB.Unscoped().First(&todoToDeleteInfo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found for permanent deletion"})
			return
		}
		log.Printf("Error finding todo %s for permanent deletion: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding todo"})
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		log.Printf("Failed to begin transaction for permanent delete: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction"})
		return
	}

	if err := tx.Unscoped().Delete(&models.Todo{}, id).Error; err != nil {
		tx.Rollback()
		log.Printf("Error permanently deleting todo %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to permanently delete todo"})
		return
	}

	if err := tx.Unscoped().Model(&models.Todo{}).Where("display_order > ?", todoToDeleteInfo.DisplayOrder).UpdateColumn("display_order", gorm.Expr("display_order - 1")).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating display_order after permanently deleting todo ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order of subsequent todos"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Transaction commit error for permanent delete: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit delete operation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo permanently deleted"})
}

func RestoreTodo(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	// Find the todo item, including soft-deleted ones.
	// We must use Unscoped() here because we are trying to find an item that is (presumably) soft-deleted.
	if err := database.DB.Unscoped().First(&todo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found in trash or elsewhere"})
		} else {
			log.Printf("Error finding todo %s for restore: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding todo to restore"})
		}
		return
	}

	// Check if the todo is actually soft-deleted.
	// todo.DeletedAt is gorm.DeletedAt (which is sql.NullTime). It's 'Valid' if the record is soft-deleted.
	if !todo.DeletedAt.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todo is not in trash, cannot restore"})
		return
	}

	// Restore the todo by setting deleted_at to NULL.
	// Using Model(&models.Todo{}) allows updating without fetching the full model into 'todo' again,
	// and Unscoped() ensures we can target a soft-deleted record for this update.
	if err := database.DB.Unscoped().Model(&models.Todo{}).Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		log.Printf("Error restoring todo %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore todo"})
		return
	}

	// Fetch the now-restored todo to return it. It should be findable without Unscoped() now.
	var restoredTodo models.Todo
	database.DB.First(&restoredTodo, id) // Populate restoredTodo with the latest data

	c.JSON(http.StatusOK, restoredTodo)
}
