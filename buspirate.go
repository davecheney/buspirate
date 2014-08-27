// Package buspirate interfaces with the binary mode of the BusPirate.
package buspirate

import (
	"fmt"
	"github.com/pkg/term"
	"time"
)

// Open opens a connection to a BusPirate module and places it into binary mode.
func Open(dev string) (*BusPirate, error) {
	t, err := term.Open(dev, term.Speed(115200), term.RawMode)
	if err != nil {
		return nil, err
	}
	bp := BusPirate{term: t}
	return &bp, bp.enterBinaryMode()
}

// BusPirate represents a connection to a remote BusPirate v4 device.
type BusPirate struct {
	term *term.Term
}

// enterBinaryMode resets the BusPirate and enters binary mode.
// http://dangerousprototypes.com/docs/Bitbang
func (bp *BusPirate) enterBinaryMode() error {
	bp.term.Flush()
	bp.term.Write([]byte{'\n', '\n', '\n', '\n', '\n', '\n', '\n', '\n', '\n', '\n'})
	const n = 30
	for i := 0; i < n; i++ {
		// send binary reset
		bp.term.Write([]byte{0x00})
		time.Sleep(10 * time.Millisecond)
		n, err := bp.term.Available()
		if err != nil {
			return err
		}
		buf := make([]byte, n)
		_, err = bp.term.Read(buf)
		if err != nil {
			return err
		}
		if string(buf) == "BBIO1" {
			return nil
		}
	}
	return fmt.Errorf("could not enter binary mode")
}

// PowerOn turns on the 5v and 3v3 regulators.
func (bp *BusPirate) PowerOn() {
	buf := []byte{0xc0}
	bp.term.Write(buf)
	bp.term.Read(buf)
}

// PowerOff turns off the 5v and 3v3 refulators.
func (bp *BusPirate) PowerOff() {
	buf := []byte{0x80}
	bp.term.Write(buf)
	bp.term.Read(buf)
}

func (bp *BusPirate) SetPWM(duty float64) {
	PRy := uint16(0x3e7f)
	OCR := uint16(float64(PRy) * duty)
	buf := []byte{0x12, 0x00, uint8(OCR >> 8), uint8(OCR), uint8(PRy >> 8), uint8(PRy)}
	bp.term.Write(buf)
	bp.term.Read(buf[:1])
}
