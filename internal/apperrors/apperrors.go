package apperrors

type AppError struct {
	LogErr  error
	UserErr error
}

func (e *AppError) Error() string {
	return e.UserErr.Error()
}

func New(logErr, userErr error) *AppError {
	return &AppError{
		LogErr:  logErr,
		UserErr: userErr,
	}
}
