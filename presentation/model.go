package presentation

// Register from JSON Model
type Register struct {
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Email ...
type Email struct {
	Email string `json:"email" binding:"required"`
}

// ActivationToken ...
type ActivationToken struct {
	Token string `json:"token" binding:"required"`
}

// ErrorModel return this if error
type ErrorModel struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// EmailLogin ...
type EmailLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
