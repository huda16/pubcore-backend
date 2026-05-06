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

// GetBooks godoc
// @Tags Books
// @Summary List all books
// @Description Retrieve a paginated, filterable, sortable list of books
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param sort query string false "Sort field (e.g. title, year, created_at)" default(created_at)
// @Param order query string false "Sort order (asc or desc)" default(desc)
// @Param title query string false "Filter by title"
// @Param genre query string false "Filter by genre"
// @Param year query int false "Filter by publication year"
// @Param author_id query int false "Filter by author ID"
// @Param publisher_id query int false "Filter by publisher ID"
// @Success 200 {object} helpers.CommonResponse
// @Failure 401 {object} helpers.CommonErrorResponse
// @Router /books [get]
func GetBooks(c *gin.Context) {
	db := database.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	offset := (page - 1) * limit

	sort := helpers.NormalizeSortField(c.DefaultQuery("sort", "created_at"))
	order := c.DefaultQuery("order", "desc")
	if order != "asc" && order != "desc" { order = "desc" }

	query := db.Model(&models.Books{}).Preload("Author").Preload("Publisher")

	if title := c.Query("title"); title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if genre := c.Query("genre"); genre != "" {
		query = query.Where("genre ILIKE ?", "%"+genre+"%")
	}
	if year := c.Query("year"); year != "" {
		query = query.Where("year = ?", year)
	}
	if authorID := c.Query("author_id"); authorID != "" {
		query = query.Where("author_id = ?", authorID)
	}
	if publisherID := c.Query("publisher_id"); publisherID != "" {
		query = query.Where("publisher_id = ?", publisherID)
	}

	var total int64
	query.Count(&total)

	var books []models.Books
	result := query.Order(sort+" "+order).Limit(limit).Offset(offset).Find(&books)
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
	to := offset + len(books)
	if len(books) == 0 { from = 0 }

	c.JSON(http.StatusOK, helpers.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    "Success",
		Data:       books,
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

// GetBookByID godoc
// @Tags Books
// @Summary Get a book by ID
// @Description Retrieve a book with its Author and Publisher
// @Security BearerAuth
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.Books
// @Failure 404 {object} map[string]string
// @Router /books/{id} [get]
func GetBookByID(c *gin.Context) {
	db := database.GetDB()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid book ID"})
		return
	}
	var book models.Books
	result := db.Preload("Author").Preload("Publisher").First(&book, uint(id))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, helpers.CommonErrorResponse{
			StatusCode: http.StatusNotFound,
			Error:      "Not Found",
			Message:    "Book not found",
		})
		return
	}
	c.JSON(http.StatusOK, helpers.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    "Success",
		Data:       book,
	})
}

// CreateBook godoc
// @Tags Books
// @Summary Create a new book
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param book body models.Books true "Book Input"
// @Success 201 {object} models.Books
// @Failure 400 {object} map[string]string
// @Router /books [post]
func CreateBook(c *gin.Context) {
	db := database.GetDB()
	var book models.Books
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}

	// Verify author exists
	var author models.Authors
	if err := db.First(&author, book.AuthorID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Author not found"})
		return
	}
	// Verify publisher exists
	var publisher models.Publishers
	if err := db.First(&publisher, book.PublisherID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Publisher not found"})
		return
	}

	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, helpers.CommonErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Bad Request",
			Message:    err.Error(),
		})
		return
	}

	db.Preload("Author").Preload("Publisher").First(&book, book.ID)
	c.JSON(http.StatusCreated, helpers.CommonResponse{
		StatusCode: http.StatusCreated,
		Message:    "Book created successfully",
		Data:       book,
	})
}

// UpdateBook godoc
// @Tags Books
// @Summary Update a book
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.Books true "Book Update Input"
// @Success 200 {object} models.Books
// @Failure 404 {object} map[string]string
// @Router /books/{id} [put]
func UpdateBook(c *gin.Context) {
	db := database.GetDB()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid book ID"})
		return
	}
	var book models.Books
	if result := db.First(&book, uint(id)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Book not found"})
		return
	}
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}

	// Verify author & publisher exist
	var author models.Authors
	if err := db.First(&author, book.AuthorID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Author not found"})
		return
	}
	var publisher models.Publishers
	if err := db.First(&publisher, book.PublisherID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Publisher not found"})
		return
	}

	if err := db.Save(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, helpers.CommonErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Bad Request",
			Message:    err.Error(),
		})
		return
	}

	db.Preload("Author").Preload("Publisher").First(&book, book.ID)
	c.JSON(http.StatusOK, helpers.CommonResponse{
		StatusCode: http.StatusOK,
		Message:    "Book updated successfully",
		Data:       book,
	})
}

// DeleteBook godoc
// @Tags Books
// @Summary Delete a book
// @Security BearerAuth
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id} [delete]
func DeleteBook(c *gin.Context) {
	db := database.GetDB()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid book ID"})
		return
	}
	var book models.Books
	if result := db.First(&book, uint(id)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Book not found"})
		return
	}
	if err := db.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.CommonErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      "Internal Server Error",
			Message:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
