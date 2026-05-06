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

// GetAuthors godoc
// @Tags Authors
// @Summary List all authors
// @Description Retrieve a paginated, filterable, sortable list of authors
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param sort query string false "Sort field (e.g. name, created_at)" default(created_at)
// @Param order query string false "Sort order (asc or desc)" default(desc)
// @Param name query string false "Filter by name (partial match)"
// @Param nationality query string false "Filter by nationality"
// @Success 200 {object} helpers.CommonResponse
// @Failure 401 {object} helpers.CommonErrorResponse
// @Failure 500 {object} helpers.CommonErrorResponse
// @Router /authors [get]
func GetAuthors(c *gin.Context) {
	db := database.GetDB()

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Sorting
	sort := helpers.NormalizeSortField(c.DefaultQuery("sort", "created_at"))
	order := c.DefaultQuery("order", "desc")
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// Filtering
	name := c.Query("name")
	nationality := c.Query("nationality")

	query := db.Model(&models.Authors{})
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if nationality != "" {
		query = query.Where("nationality ILIKE ?", "%"+nationality+"%")
	}

	var total int64
	query.Count(&total)

	var authors []models.Authors
	result := query.
		Order(sort + " " + order).
		Limit(limit).
		Offset(offset).
		Find(&authors)

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
	to := offset + len(authors)
	if len(authors) == 0 {
		from = 0
	}

	c.JSON(http.StatusOK, helpers.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    "Success",
		Data:       authors,
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

// GetAuthorByID godoc
// @Tags Authors
// @Summary Get a single author by ID
// @Description Retrieve an author and their books by ID
// @Security BearerAuth
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} models.Authors
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /authors/{id} [get]
func GetAuthorByID(c *gin.Context) {
	db := database.GetDB()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid author ID",
		})
		return
	}

	var author models.Authors
	result := db.Preload("Books").First(&author, uint(id))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, helpers.CommonErrorResponse{
			StatusCode: http.StatusNotFound,
			Error:      "Not Found",
			Message:    "Author not found",
		})
		return
	}

	c.JSON(http.StatusOK, helpers.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    "Success",
		Data:       author,
	})
}

// CreateAuthor godoc
// @Tags Authors
// @Summary Create a new author
// @Description Add a new author to the platform
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param author body models.Authors true "Author Input"
// @Success 201 {object} helpers.CommonResponse
// @Failure 400 {object} helpers.CommonErrorResponse
// @Router /authors [post]
func CreateAuthor(c *gin.Context) {
	db := database.GetDB()

	var author models.Authors
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	if err := db.Create(&author).Error; err != nil {
		c.JSON(http.StatusBadRequest, helpers.CommonErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Bad Request",
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, helpers.CommonResponse{
		StatusCode: http.StatusCreated,
		Message:    "Author created successfully",
		Data:       author,
	})
}

// UpdateAuthor godoc
// @Tags Authors
// @Summary Update an existing author
// @Description Modify author details by ID
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Param author body models.Authors true "Author Update Input"
// @Success 200 {object} models.Authors
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /authors/{id} [put]
func UpdateAuthor(c *gin.Context) {
	db := database.GetDB()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid author ID",
		})
		return
	}

	var author models.Authors
	if result := db.First(&author, uint(id)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Author not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	if err := db.Save(&author).Error; err != nil {
		c.JSON(http.StatusBadRequest, helpers.CommonErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Bad Request",
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, helpers.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    "Author updated successfully",
		Data:       author,
	})
}

// DeleteAuthor godoc
// @Tags Authors
// @Summary Delete an author
// @Description Remove an author by ID (only if they have no books)
// @Security BearerAuth
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} helpers.CommonResponse
// @Failure 400 {object} helpers.CommonErrorResponse
// @Failure 404 {object} helpers.CommonErrorResponse
// @Router /authors/{id} [delete]
func DeleteAuthor(c *gin.Context) {
	db := database.GetDB()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid author ID",
		})
		return
	}

	var author models.Authors
	if result := db.First(&author, uint(id)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Author not found",
		})
		return
	}

	// Check if author has books
	var bookCount int64
	db.Model(&models.Books{}).Where("author_id = ?", id).Count(&bookCount)
	if bookCount > 0 {
		c.JSON(http.StatusBadRequest, helpers.CommonErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Bad Request",
			Message:    "Cannot delete author with existing books. Remove or reassign books first.",
		})
		return
	}

	if err := db.Delete(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.CommonErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      "Internal Server Error",
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, helpers.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    "Author deleted successfully",
	})
}
