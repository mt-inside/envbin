package trie

import (
	"fmt"
	"net/url"
	"time"
)

type Value interface {
	Render() string
}

type some struct {
	Value string `json:"value"`
}

func Some(value string) some {
	return some{value}
}
func Optional(value string) Value {
	if value != "" {
		return some{value}
	}
	return NotPresent()
}

func (s some) Render() string {
	return s.Value
}

type notPresent struct{}

func NotPresent() notPresent {
	return notPresent{}
}

func (np notPresent) Render() string {
	return "NotPresent"
}

type erro struct {
	Err error
}

func Error(err error) Value {
	// TODO 401 -> Forbidden() etc
	if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
		return Timeout(time.Second) // FIXME: duration
	} else {
		return erro{err}
	}
}

func (e erro) Render() string {
	return fmt.Sprintf("Error: %v", e.Err)
}

type timeout struct {
	D time.Duration
}

func Timeout(d time.Duration) timeout {
	return timeout{d}
}

func (t timeout) Render() string {
	return fmt.Sprintf("Timed Out (waited %v)", t.D)
}

type forbidden struct{}

func Forbidden() forbidden {
	return forbidden{}
}

func (f forbidden) Render() string {
	return "Forbidden"
}
