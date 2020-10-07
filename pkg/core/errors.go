package core

type ErrorCode string

const (
	ETODO     ErrorCode = "todo"
	EINVALID  ErrorCode = "invalid"
	EINTERNAL ErrorCode = "internal"
)

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`

	Err error  `json:"-"`
	Op  string `json:"-"`
}

func (e Error) Error() string {
	return e.Message
}
