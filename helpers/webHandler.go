package helpers

// Auth Input/Response Types

// RegisterInput is the request body for user registration.
type RegisterInput struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"secret123"`
	Bio      string `json:"bio" example:"A passionate reader"`
}

// RegisterResponse is the response body after successful registration.
type RegisterResponse struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
}

// LoginInput is the request body for user login.
type LoginInput struct {
	Email    string `json:"email" example:"john@example.com" valid:"required"`
	Password string `json:"password" example:"secret123" valid:"required"`
}

// LoginResponse is the response body after successful login.
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOi..."`
}

// Author Input/Response Types

// AuthorInput is the request body for creating or updating an author.
type AuthorInput struct {
	Name        string `json:"name" example:"J.K. Rowling" valid:"required~Author name is required"`
	Bio         string `json:"bio" example:"British author best known for Harry Potter."`
	Nationality string `json:"nationality" example:"British"`
}

// Publisher Input/Response Types

// PublisherInput is the request body for creating or updating a publisher.
type PublisherInput struct {
	Name    string `json:"name" example:"Bloomsbury" valid:"required~Publisher name is required"`
	Address string `json:"address" example:"50 Bedford Square, London"`
	Website string `json:"website" example:"https://www.bloomsbury.com"`
}

// Book Input/Response Types

// BookInput is the request body for creating or updating a book.
type BookInput struct {
	Title       string `json:"title" example:"Harry Potter and the Philosopher's Stone" valid:"required~Title is required"`
	ISBN        string `json:"isbn" example:"978-0-7475-3269-9"`
	Description string `json:"description" example:"A young boy discovers he is a wizard."`
	Genre       string `json:"genre" example:"Fantasy"`
	Year        int    `json:"year" example:"1997"`
	AuthorID    uint   `json:"author_id" example:"1" valid:"required~Author is required"`
	PublisherID uint   `json:"publisher_id" example:"1" valid:"required~Publisher is required"`
}

// Meta holds metadata for paginated responses.
type Meta struct {
	Total       int64 `json:"total" example:"42"`
	Limit       int   `json:"limit" example:"10"`
	CurrentPage int   `json:"currentPage" example:"1"`
	LastPage    int   `json:"lastPage" example:"5"`
	From        int   `json:"from" example:"1"`
	To          int   `json:"to" example:"10"`
}

// CommonResponse is the standard success response body.
type CommonResponse struct {
	Data       interface{} `json:"data"`
	Message    string      `json:"message" example:"Success"`
	StatusCode int         `json:"statusCode" example:"200"`
	Meta       *Meta       `json:"meta,omitempty"`
}

// CommonErrorResponse is the standard error response body.
type CommonErrorResponse struct {
	Error      string      `json:"error" example:"Bad Request"`
	Message    interface{} `json:"message"`
	StatusCode int         `json:"statusCode" example:"400"`
}

// DeleteResponse is returned after a successful delete operation.
type DeleteResponse struct {
	Message string `json:"message" example:"Data deleted successfully"`
}
