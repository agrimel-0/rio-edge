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

func SetLineValue(line gpiod.Line, value int32) error {

	err := line.SetValue(int(value))

	return err
}

// Initialize a new pin. Requires gpiochip, line offset, default value, and an alias.
func newPin(pinstStat Pins) (*IoPin, error) {

	newLine, err := gpiod.RequestLine(pinstStat.GpioChip, pinstStat.LineOffset, gpiod.AsOutput(b2i(pinstStat.Output)))
	if err != nil {
		return nil, err
	}
	// defer newLine.Close()

	io := IoPin{
		Line:     newLine,
		Alias:    pinstStat.Alias,
		Value:    pinstStat.Value,
		AsOutput: pinstStat.Output,
		GpioChip: pinstStat.GpioChip,
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
