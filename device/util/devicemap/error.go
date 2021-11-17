package devicemap

type errorType int

// const (
// 	ErrorDuplicate errorType = iota
// 	ErrorNotFound
// )

// OpError ...
// type OpError struct {
// 	err string
// }

// func (e OpError) Error() string { return e.err }

type OpErrorDuplicate struct {
	err string
}

func (e OpErrorDuplicate) Error() string { return e.err }

type OpErrorNotFound struct {
	err string
}

func (e OpErrorNotFound) Error() string { return e.err }
