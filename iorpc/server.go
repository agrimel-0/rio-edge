/*
Copyright © 2022 Federico Giovine <giovine.federico@gmail.com>

*/

package iorpc

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	pb "192.168.1.153/blue-box/blueboxproto"
	"google.golang.org/grpc"
)

/*
	gRPC Server stuff
*/

// Server struct
type server struct {
	pb.UnimplementedBlueBoxServer

	exportedPins []*ioPin // Slice containing the exported pins
}

// gRPC call for setting IO pin by offset
func (s *server) SetGPIObyOffset(ctx context.Context, lineOffset *pb.GPIOselected) (*pb.ServerResponse, error) {

	offsetSelected := lineOffset.GPIOLineOffset

	log.Printf("setting '%s' GPIO by offset %d  with value %d\n", lineOffset.GPIOLineAlias, offsetSelected, lineOffset.GetGPIOLineValue())

	// Select the matching exported io
	ioSelected, err := s.findPinByOffset(offsetSelected)
	if err != nil {
		return &pb.ServerResponse{ErrorString: "error found"}, err
	}

	// Set the line value. Should it throw an error if you are setting a value that it's already set at?
	ioSelected.line.SetValue(int(lineOffset.GetGPIOLineValue()))

	return &pb.ServerResponse{ErrorString: "none"}, nil
}

// gRPC call for setting IO pin by alias. Would be cool to be able to request a list of aliases from a client
func (s *server) SetGPIObyAlias(ctx context.Context, lineAlias *pb.GPIOselected) (*pb.ServerResponse, error) {

	aliasSelected := lineAlias.GetGPIOLineAlias()

	log.Printf("setting '%s' GPIO by alias %d  with value %d\n", lineAlias.GPIOLineAlias, lineAlias.GetGPIOLineOffset(), lineAlias.GetGPIOLineValue())

	// Select the matching exported io
	ioSelected, err := s.findPinByAlias(aliasSelected)
	if err != nil {
		return &pb.ServerResponse{ErrorString: "error found"}, err
	}

	// Set the line value. Should it throw an error if you are setting a value that it's already set at?
	ioSelected.line.SetValue(int(lineAlias.GetGPIOLineValue()))

	return &pb.ServerResponse{ErrorString: "none"}, nil
}

// find the matching pin in the exported pins when searching by offset value
func (s *server) findPinByOffset(offsetSelected int32) (*ioPin, error) {
	for _, x := range s.exportedPins {
		if x.line.Offset() == int(offsetSelected) {
			return x, nil
		}
	}
	return s.exportedPins[0], errors.New("no line found by offset")
}

// find the matching pin in the exported pins when searching by alias value
func (s *server) findPinByAlias(aliasSelected string) (*ioPin, error) {
	for _, x := range s.exportedPins {
		if x.alias == aliasSelected {
			return x, nil
		}
	}
	return s.exportedPins[0], errors.New("no line found by alias")
}

// Serve the server
func StartServer(port int) {

	flag.Parse()

	/*
		GPIO setup stuff. Must be done before gRPC server setup!
	*/
	lights, err := newPin(17, 1, 1, "gpiochip0", "lights") // 1 so it is an output
	if err != nil {
		log.Fatalf("failed to set a new pin: %v", err) // 1 so it is an output
	}
	log.Printf("setup GPIO with alias %s on line %d of chip %s\n", lights.alias, lights.line.Offset(), lights.line.Chip())
	pump, err := newPin(18, 1, 1, "gpiochip0", "pump")
	if err != nil {
		log.Fatalf("failed to set a new pin: %v", err)
	}
	log.Printf("setup GPIO with alias %s on line %d of chip %s\n", pump.alias, pump.line.Offset(), pump.line.Chip())

	// Assign the exported pins
	exportedPins := []*ioPin{lights, pump}
	server := server{exportedPins: exportedPins}

	/*
		gRPC stuff
	*/
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBlueBoxServer(s, &server) // instead opf &server{} we use &server since it is already initialized
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
