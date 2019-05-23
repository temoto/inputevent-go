package inputevent

import (
	"encoding/hex"
	"math/rand"
	"strings"
	"syscall"
	"testing"
	"unsafe"
)

func TestParse(t *testing.T) {
	t.Parallel()

	type Case struct {
		name   string
		evhex  string
		expect InputEvent
	}
	cases := []Case{
		// FIXME should break on big-endian arch
		{"key-press", "0100 8d00 01000000", InputEvent{Type: EV_KEY, Code: 141, Value: 1}},
		{"key-up___", "0100 2000 00000000", InputEvent{Type: EV_KEY, Code: 32, Value: 0}},
	}
	const timevalSizeof = int(unsafe.Sizeof(syscall.Timeval{}))
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			evb, err := hex.DecodeString(strings.Replace(c.evhex, " ", "", -1))
			if err != nil {
				t.Fatalf("hex=%s err=%v", c.evhex, err)
			}

			var src [EventSizeof]byte
			tv := syscall.NsecToTimeval(rand.Int63())
			c.expect.Time = tv

			// Sorry for not using proper serialization.
			copy(src[:timevalSizeof], (*(*[timevalSizeof]byte)(unsafe.Pointer(&tv)))[:])

			copy(src[timevalSizeof:], evb)

			event, err := parse(src)
			if err != nil {
				t.Fatal(err)
			}
			if event != c.expect {
				t.Errorf("parsed=%#v expected=%#v", event, c.expect)
			}
		})
	}
}
