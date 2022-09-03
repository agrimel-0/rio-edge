package server

import (
	"github.com/warthog618/gpiod"
)

type IoPin struct {
	Alias    string      // line alias such as pump, fan, etc
	Line     *gpiod.Line // gpio line struct
	GpioChip string      // gpio chip
	Value    int         // current value of the pin.
	AsOutput bool        // 1 is output, 0 is input
}

// IoFromConfig, generate a gpiod pin array from pin config
func IoFromConfig(pinMap []map[string]Pins) ([]*IoPin, []error) {

	var ioPins []*IoPin
	var errorList []error

	for _, pinStats := range pinMap {
		for _, stat := range pinStats {
			ioPin, err := newPin(stat) // create a new pin
			if err != nil {
				errorList = append(errorList, err) // append errors to the list
			}

			ioPins = append(ioPins, ioPin) // add the new pin to the list of pins

		}
	}

	return ioPins, errorList
}

func (io *IoPin) SetLineValue(value int32) error {
	err := io.Line.SetValue(int(value))

	io.Value = int(value)

	return err
}

// Initialize a new pin. Requires gpiochip, line offset, default value, and an alias.
func newPin(pinStat Pins) (*IoPin, error) {

	newLine, err := gpiod.RequestLine(pinStat.GpioChip, pinStat.LineOffset, gpiod.AsOutput(b2i(pinStat.Output)))
	if err != nil {
		return nil, err
	}
	// defer newLine.Close()

	io := IoPin{
		Line:     newLine,
		Alias:    pinStat.Alias,
		Value:    pinStat.Value,
		AsOutput: pinStat.Output,
		GpioChip: pinStat.GpioChip,
	}

	if io.AsOutput {
		err = newLine.SetValue(io.Value) // default value for the line
		return &io, err
	}

	return &io, nil
}

func b2i(b bool) int {
	if b {
		return 0
	}
	return 1
}
