package buspirate

import "github.com/davecheney/buspirate"
import "time"

// ExampleBusPirate pulses a LED connected to the AUX pin.
func ExampleBusPirate_SetPWM() {
	bp, err := buspirate.Open("/dev/ttyACM0")
	if err != nil {
		panic(err)
	}
	duty := 0.1
	delta := 0.1
	for {
		bp.SetPWM(duty)
		time.Sleep(50 * time.Millisecond)
		duty += delta
		if duty > 1.0 {
			duty = 1.0
			delta = -delta
		}
		if duty < 0 {
			duty = 0.0
			delta = -delta
		}
	}
}
