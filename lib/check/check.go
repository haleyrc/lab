package check

import (
	"errors"
	"reflect"
	"testing"
)

type EmptyReporter interface {
	IsEmpty() bool
}

type Checker struct {
	t      *testing.T
	failed bool
}

func New(t *testing.T) *Checker {
	return &Checker{t: t}
}

func (c *Checker) Then(f func()) *Checker {
	c.t.Helper()

	if !c.failed {
		f()
	}

	return c
}

func (c *Checker) Fatal() {
	c.t.Helper()

	failed := c.failed
	c.failed = false

	if failed {
		c.t.FailNow()
	}
}

func (c *Checker) Nil(i interface{}) *Checker {
	c.t.Helper()

	c.failed = false

	if !isNil(i) {
		c.failed = true
		c.t.Errorf("expected i to be nil, but got %#v", i)
	}

	return c
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func (c *Checker) ErrorCause(err, want error) *Checker {
	c.t.Helper()

	c.failed = false

	cause := errors.Unwrap(err)
	c.Equals(cause, want)

	return c
}

func (c *Checker) Error(err error) *Checker {
	c.t.Helper()

	c.failed = false

	if err == nil {
		c.failed = true
		c.t.Error("expected an error, but got none")
	}

	return c
}

func (c *Checker) OK(err error) *Checker {
	c.t.Helper()

	c.failed = false

	if err != nil {
		c.failed = true
		c.t.Error(err)
	}

	return c
}

func (c *Checker) True(pred bool) *Checker {
	c.t.Helper()

	c.failed = false

	if !pred {
		c.failed = true
		c.t.Error("expected predicate to be true, but wasn't")
	}

	return c
}

func (c *Checker) False(pred bool) *Checker {
	c.t.Helper()

	c.failed = false

	if pred {
		c.failed = true
		c.t.Error("expected predicate to be false, but wasn't")
	}

	return c
}

func (c *Checker) NotEmpty(val interface{}) *Checker {
	c.t.Helper()

	c.failed = false

	switch v := val.(type) {
	case EmptyReporter:
		if v.IsEmpty() {
			c.failed = true
			c.t.Error("expected val to not be empty, but was")
		}
	case string:
		if v == "" {
			c.failed = true
			c.t.Error("expected val to not be empty, but was")
		}
	default:
		c.failed = true
		c.t.Errorf("unhandled type %T for val %v", val, val)
	}

	return c
}

func (c *Checker) NotEquals(got, want interface{}) *Checker {
	c.t.Helper()

	c.failed = false

	if reflect.DeepEqual(got, want) {
		c.failed = true
		c.t.Error("expected values to be different, but weren't")
	}

	return c
}

// TODO (RCH): Fail early if types aren't equal
// TODO (RCH): Type switch for basic types
func (c *Checker) Equals(got, want interface{}) *Checker {
	c.t.Helper()

	c.failed = false

	switch w := want.(type) {
	case int64:
		g, ok := got.(int64)
		if !ok {
			c.failed = true
			c.t.Errorf("mismatched types.\nwanted:\n\t(%T)(%v\ngot:\n\t(%T)(%v)", want, want, got, got)
			return c
		}
		if g != w {
			c.failed = true
		}
	default:
		if !reflect.DeepEqual(got, want) {
			c.failed = true
		}
	}

	if c.failed {
		c.t.Errorf("expected values to be equal.\nwanted:\n\t%v\ngot:\n\t%v", want, got)
	}

	return c
}
