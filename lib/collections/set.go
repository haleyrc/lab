package collections

type Set interface {
	Add(s string)
	AddN(ss ...string)
	Slice() []string
}
