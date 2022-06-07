package trie

import (
	"fmt"
	"net/url"
	"time"
)

// TODO: this should be a lot richer
// - iface should have IsError - this is basically a Result - are you an Ok or one of the error types
//   - errors: UnknownError(err), Timeout(MaybeNil<err>, duration), Forbidden(MaybeNil<err>, permissionNeeded string) - we build these abstrations cause a) it's not obvious when holding a go error which class it is, or how to get eg its duration, and b) there might be other non-error sources of this info
//   - use IsError in Walk
//   - rename Some to Ok
// - also IsOk (the other possible type for Result), again several types:
//   - Some, which should have
//     - id (eg SomeTokenFoo)
//     - pretty name (for client / dump use)
//     - type (ie tagged union)
//     - value
//     - numeric value classes should have units - magnitude and dimensions. See if there's a nice lib for this, else write it - will be fun!
//   - NotPresent - this is a thing, and i can tell you it's not here
//   - Unknown - I can't tell, but not because of an error, probably just cause I'm not looking hard enough (eg filesystem type of unmounted partitions)

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
	// TODO: calling site to go this, cause it may be dealing with non-stdlib error types
	// - ie they should look at the errors they have in hand (or non-error data) and call Timeout, Forbidden direct
	// - util function for doing the below, ie turning common built-in errors into Timeout etc
	// - doc that this is for unknown error types (maybe rename?)
	// - this function should then call that, and if it /can/ extact something, panic, cause the call-site shouldn't be giving it here
	if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
		return Timeout(time.Second) // FIXME: duration - will go away when call sites call Timeout() direct
		// TODO 401 -> Forbidden() etc
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
