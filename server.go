// server.go
package main

import (
	"net"
	"net/rpc"
	"log"
	// "fmt"
	// "image"
	// "os"
	// _ "image/png"
	// _ "image/jpeg"
	// "os"
	// "math"
)

type Args struct {
	Red []uint32
	Green []uint32
	Blue []uint32
	Length int
}

type Image struct{}

func (t *Image) ToGrayscale(args *Args, reply *[]uint8) error {
	var temp = make([]uint8, args.Length)
	// var red, green, blue float64

	for i := 0 ; i < args.Length ; i++ {
		temp[i] = (uint8)(((2125*args.Red[i] + 7154*args.Green[i] + 721*args.Blue[i] + 5000) / 10000) >> 8 )

		// red = math.Sqrt((float64) (args.Red[i]))
		// green = math.Sqrt((float64) (args.Green[i]))
		// blue = math.Sqrt((float64) (args.Blue[i]))
		// temp[i] = (uint8)((red * 0.299) + (green * 0.587) + (blue * 0.114)) // ((red + green + blue) / 3)// 
	}

	*reply = temp
	return nil
}

func main(){
	cal := new(Image)
	rpc.Register(cal)
	listener, e := net.Listen("tcp", ":6060")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			log.Printf("new connection established\n")
			go rpc.ServeConn(conn)
		}
	}
}
