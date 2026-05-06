package controllers

import (
	"math"
	"net/http"
	"pubcore/database"
	"pubcore/helpers"
	"pubcore/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPublishers godoc
// @Tags Publishers
// @Summary List all publishers
// @Description Retrieve a paginated, filterable, sortable list of publishers
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param sort query string false "Sort field" default(created_at)
// @Param order query string false "Sort order (asc or desc)" default(desc)
// @Param name query string false "Filter by name"
// @Success 200 {object} helpers.CommonResponse
// @Failure 401 {object} helpers.CommonErrorResponse
// @Router /publishers [get]
func GetPublishers(c *gin.Context) {
	db := database.GetDB()
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	offset := (page - 1) * limit

	sort := helpers.NormalizeSortField(c.DefaultQuery("sort", "created_at"))
	order := c.DefaultQuery("order", "desc")
	if order != "asc" && order != "desc" { order = "desc" }

	name := c.Query("name")
	query := db.Model(&models.Publishers{})
	if name != "" { query = query.Where("name ILIKE ?", "%"+name+"%") }

	var total int64
	query.Count(&total)

	var publishers []models.Publishers
	result := query.Order(sort+" "+order).Limit(limit).Offset(offset).Find(&publishers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.CommonErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      "Internal Server Error",
			Message:    result.Error.Error(),
		})
		return
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	from := offset + 1
	to := offset + len(publishers)
	if len(publishers) == 0 { from = 0 }

	c.JSON(http.StatusOK, helpers.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    "Success",
		Data:       publishers,
		Meta: &helpers.Meta{
			Total:       total,
			Limit:       limit,
			CurrentPage: page,
			LastPage:    lastPage,
			From:        from,
			To:          to,
		},
	})
}

// GetPublisherByID godoc
// @Tags Publishers
// @Summary Get a publisher by ID
// @Security BearerAuth
// @Produce json
// @Param id path int true "Publisher ID"
// @Success 200 {object} models.Publishers
// @Failure 404 {object} map[string]string
// @Router /publishers/{id} [get]
func GetPublisherByID(c *gin.Context) {
	db := database.GetDB()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid publisher ID"})
		return
	}
	var publisher models.Publishers
	if result := db.Preload("Books").First(&publisher, uint(id)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Publisher not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": publisher})
}

// CreatePublisher godoc
// @Tags Publishers
// @Summary Create a new publisher
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param publisher body models.Publishers true "Publisher Input"
// @Success 201 {object} models.Publishers
// @Failure 400 {object} map[string]string
// @Router /publishers [post]
func CreatePublisher(c *gin.Context) {
	db := database.GetDB()
	var publisher models.Publishers
	if err := c.ShouldBindJSON(&publisher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	if err := db.Create(&publisher).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Publisher created successfully", "data": publisher})
}

// UpdatePublisher godoc
// @Tags Publishers
// @Summary Update a publisher
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Publisher ID"
// @Param publisher body models.Publishers true "Publisher Update Input"
// @Success 200 {object} models.Publishers
// @Failure 404 {object} map[string]string
// @Router /publishers/{id} [put]
func UpdatePublisher(c *gin.Context) {
	db := database.GetDB()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid publisher ID"})
		return
	}
	var publisher models.Publishers
	if result := db.First(&publisher, uint(id)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Publisher not found"})
		return
	}
	if err := c.ShouldBindJSON(&publisher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	if err := db.Save(&publisher).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Publisher updated successfully", "data": publisher})
}

// DeletePublisher godoc
// @Tags Publishers
// @Summary Delete a publisher
// @Security BearerAuth
// @Produce json
// @Param id path int true "Publisher ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /publishers/{id} [delete]
func DeletePublisher(c *gin.Context) {
	db := database.GetDB()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid publisher ID"})
		return
	}
	var publisher models.Publishers
	if result := db.First(&publisher, uint(id)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Publisher not found"})
		return
	}
	var bookCount int64
	db.Model(&models.Books{}).Where("publisher_id = ?", id).Count(&bookCount)
	if bookCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Cannot delete publisher with existing books."})
		return
	}
	if err := db.Delete(&publisher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Publisher deleted successfully"})
}
