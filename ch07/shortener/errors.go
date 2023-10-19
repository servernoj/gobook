package shortener

type AppError struct {
	Error error
	Code  int
}
