package ArgReader

func NewInsufficientArgumentError() *insufficientArgumentError {
	return &insufficientArgumentError{"You have to provide url as your first argument"}
}

type insufficientArgumentError struct {
	msg string
}

func (e *insufficientArgumentError) Error() string {
	return e.msg
}