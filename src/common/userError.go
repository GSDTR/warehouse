package Common

type errorString struct {
	s string
}
func (e *errorString) Error() string {
	return e.s
}
func UserError(text string) error {
	return &errorString{text}
}
