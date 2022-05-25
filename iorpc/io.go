/*
Copyright © 2022 Federico Giovine <giovine.federico@gmail.com>

*/

package iorpc

import "github.com/warthog618/gpiod"

type ioPin struct {
	line     *gpiod.Line // gpio line struct
	alias    string      // line alias such as pump, fan, etc
	value    int         // current value of the pin.
	AsOutput int         // 1 is output, 0 is input
}

// Initialize a new pin. Requires gpiochip, line offset, default value, and an alias.
func newPin(lineOffset, startValue, AsOutput int, gpioChip, alias string) (*ioPin, error) {
	newLine, err := gpiod.RequestLine(gpioChip, lineOffset, gpiod.AsOutput(AsOutput))
	if err != nil {
		return nil, err
	}
	defer newLine.Close() // Not sure if actually needed

	newLine.SetValue(startValue) // default value for the line

	io := ioPin{
		line:     newLine,
		alias:    alias,
		value:    startValue,
		AsOutput: AsOutput,
	}

	return &io, nil
}
