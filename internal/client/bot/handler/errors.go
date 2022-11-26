package handler

type GifTypeNotSpecifiedError struct{}

func (e *GifTypeNotSpecifiedError) Error() string {
	return "gif type not specified"
}

type GifsNotFoundError struct{}

func (e *GifsNotFoundError) Error() string {
	return "gifs not found"
}
