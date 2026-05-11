package usecase

// ValidationError は入力バリデーションエラー
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field != "" {
		return "validation error: " + e.Field + " - " + e.Message
	}
	return "validation error: " + e.Message
}

// NotFoundError はリソース未検出エラー
type NotFoundError struct {
	Resource string
	ID       interface{}
}

func (e NotFoundError) Error() string {
	return e.Resource + " not found"
}

// IsValidationError はバリデーションエラーか判定
func IsValidationError(err error) bool {
	_, ok := err.(ValidationError)
	return ok
}

// IsNotFoundError はリソース未検出エラーか判定
func IsNotFoundError(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}

// AsValidationError はエラーをバリデーションエラーにキャスト
func AsValidationError(err error) (ValidationError, bool) {
	valErr, ok := err.(ValidationError)
	return valErr, ok
}

// AsNotFoundError はエラーをリソース未検出エラーにキャスト
func AsNotFoundError(err error) (NotFoundError, bool) {
	notFoundErr, ok := err.(NotFoundError)
	return notFoundErr, ok
}
