package server

import (
	"context"
	"errors"

	pb "github.com/agrimel-0/rio-grpc"
	"github.com/sirupsen/logrus"
)

// SetGPIObyOffset gRPC call for setting IO pin by offset
func (s *server) SetGPIObyOffset(ctx context.Context, in *pb.GPIOselected) (*pb.ServerResponse, error) {
	offsetSelected := in.GPIOLineOffset

	// Select the matching exported io
	ioSelected, err := s.findPinByOffset(offsetSelected)
	if err != nil {
		return &pb.ServerResponse{ResponseString: "error found"}, err
	}

	logrus.Infof("setting '%s' GPIO by offset %d  with value %d\n", ioSelected.Alias, ioSelected.Line.Offset(), in.GetGPIOLineValue())

	// Set the line value. Should it throw an error if you are setting a value that it's already set at?
	err = SetLineValue(*ioSelected.Line, in.GPIOLineValue)
	if err != nil {
		logrus.Errorf("error setting by offset", err)
		return &pb.ServerResponse{ResponseString: err.Error()}, err
	}

	return &pb.ServerResponse{ResponseString: "none"}, nil
}

// SetGPIObyAlias gRPC call for setting IO pin by alias. Would be cool to be able to request a list of aliases from a client
func (s *server) SetGPIObyAlias(ctx context.Context, in *pb.GPIOselected) (*pb.ServerResponse, error) {

	aliasSelected := in.GetGPIOLineAlias()

	// Select the matching exported io
	ioSelected, err := s.findPinByAlias(aliasSelected)
	if err != nil {
		return &pb.ServerResponse{ResponseString: err.Error()}, err
	}

	logrus.Infof("setting '%s' GPIO by alias %d  with value %d\n", ioSelected.Alias, ioSelected.Line.Offset(), in.GetGPIOLineValue())

	// Set the line value. Should it throw an error if you are setting a value that it's already set at?
	err = SetLineValue(*ioSelected.Line, in.GPIOLineValue)
	if err != nil {
		logrus.Errorf("error setting by alias", err)
		return &pb.ServerResponse{ResponseString: err.Error()}, err
	}

	return &pb.ServerResponse{ResponseString: "none"}, nil
}

// GetGPIOList gRPC call for listing IO pins.
func (s *server) GetGPIOList(in *pb.ClientRequest, stream pb.Rio_GetGPIOListServer) error {

	for _, x := range s.exportedPins {

		logrus.Debugf("line: %d", x.Line.Offset())
		logrus.Debugf("value: %d", x.Value)
		logrus.Debugf("alias: %s", x.Alias)
		logrus.Debugf("output: %v", x.AsOutput)
		logrus.Debugf("gpiochip: %s", x.GpioChip)

		// go from x (type IOpin) to gpioToStream (GPIOselected) ...
		gpioToStream := pb.GPIOselected{
			GPIOLineOffset: int32(x.Line.Offset()),
			GPIOLineValue:  int32(x.Value),
			GPIOLineAlias:  x.Alias,
			GPIOOutput:     x.AsOutput,
			GPIOChip:       x.GpioChip,
		}

		logrus.Debugf("sending: %v\n", gpioToStream)

		err := stream.Send(&gpioToStream)
		if err != nil {
			return err
		}
	}
	return nil
}

// find the matching pin in the exported pins when searching by offset value
func (s *server) findPinByOffset(offsetSelected int32) (*IoPin, error) {
	for _, x := range s.exportedPins {
		if x.Line.Offset() == int(offsetSelected) {
			return x, nil
		}
	}
	return s.exportedPins[0], errors.New("no line found by offset")
}

// find the matching pin in the exported pins when searching by alias value
func (s *server) findPinByAlias(aliasSelected string) (*IoPin, error) {
	for _, x := range s.exportedPins {
		if x.Alias == aliasSelected {
			return x, nil
		}
	}
	return s.exportedPins[0], errors.New("no line found by alias")
}
