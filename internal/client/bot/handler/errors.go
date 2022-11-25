package handler

type GifTypeNotSpecifiedError struct{}

func (e *GifTypeNotSpecifiedError) Error() string {
	return "gif type not specified"
}
