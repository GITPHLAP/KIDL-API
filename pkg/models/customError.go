package models

type CustomError struct {
	Message string
	Custom  bool
}

func (c CustomError) Error() string {
	return c.Message
}
