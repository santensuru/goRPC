// client.go
package main

import (
	"fmt"
	"net/rpc"
	"log"
	"image"
	"os"
	// "encoding/base64"
	"image/png"
	"image/jpeg"
	"image/color"
	// "math"
	// "strconv"
	"io/ioutil"
	"strings"
	"time"
	"container/list"
	"sync"
)

type Args struct {
	Red []uint32
	Green []uint32
	Blue []uint32
	Length int
}

type Pair struct {
	Name string
	Format string
}

const rootPath = "C:/cygwin64/home/user/coba/SISTER/"

func handleList(l *list.List, index string, dest string, wg *sync.WaitGroup) (result string) {
	defer wg.Done()
	client, err := rpc.Dial("tcp", dest)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	timeStart := time.Now()

	for i := l.Front(); i != nil; i = i.Next() {

		fmt.Print(index + " " + i.Value.(Pair).Name + "\n")
		reader, err := os.Open(rootPath+i.Value.(Pair).Name) // os.Open(strconv.Itoa(q)+".jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		m, _, err := image.Decode(reader)
		if err != nil {
			log.Fatal(err)
		}
		bounds := m.Bounds()

		var red = make([]uint32, bounds.Max.Y * bounds.Max.X)
		var green = make([]uint32, bounds.Max.Y * bounds.Max.X)
		var blue = make([]uint32, bounds.Max.Y * bounds.Max.X)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := m.At(x, y).RGBA()
				// fmt.Printf("%d %d %d\n", r, g, b)
				red[x + y*bounds.Max.X] = r
				green[x + y*bounds.Max.X] = g
				blue[x + y*bounds.Max.X] = b
				// fmt.Print(r , "\n")
				// fmt.Print(float64(r) * 0.299)
			}
		}

		// fmt.Print(bounds.Min.Y, " ",  bounds.Max.Y)

		args := &Args{red, green, blue, bounds.Max.Y * bounds.Max.X}

		// fmt.Print(args.Length)

		var replay = make([]uint8, bounds.Max.Y * bounds.Max.X)
		err = client.Call("Image.ToGrayscale", args, &replay)
		if err != nil {
			log.Fatal("arith error:", err)
		} else {
			dest := image.NewGray( image.Rect(0,0,bounds.Max.X, bounds.Max.Y))
				


			for y := 0; y < bounds.Max.Y; y++ {

				for x := 0; x < bounds.Max.X; x++ {
					// fmt.Printf("%d", replay[y])
					newcolor := color.Gray{replay[x + y*bounds.Max.X]}
					// fmt.Print(math.Sqrt((float64) (replay[x + y*bounds.Max.X])), " ", replay[x + y*bounds.Max.X], "\n")
					dest.SetGray(x, y, newcolor)
				}

				
			}

			toimg, _ := os.Create(rootPath+"BW2/BW-"+i.Value.(Pair).Name) // os.Create("BW-"+strconv.Itoa(q)+".jpg")
			defer toimg.Close()

			if i.Value.(Pair).Format == "jpg" {
				jpeg.Encode(toimg, dest, &jpeg.Options{jpeg.DefaultQuality})
			} else if i.Value.(Pair).Format == "png" {
				png.Encode(toimg, dest)
			}
		}
	}

	timeEnd := time.Now()

	fmt.Printf("%s %v\n", index, timeEnd.Sub(timeStart))

	result = "yey"
	return
	
}
 
func main(){
	
	// client, err := rpc.Dial("tcp", "127.0.0.1:6060")
	// if err != nil {
	// 	log.Fatal("dialing:", err)
	// }

	// Synchronous call

	var format string
	// var rootPath string
	// rootPath = "C:/cygwin64/home/user/coba/SISTER/"

	l1 := list.New()
	l2 := list.New()
	l3 := list.New()
	l4 := list.New()
	var i int
	i = 0

	// timeStart := time.Now()

	files, _ := ioutil.ReadDir(rootPath)
	for _, f := range files {
			// fmt.Println(f.Name())
		if strings.Contains(f.Name(), ".") {
			s := strings.Split(f.Name(), ".")
			// fmt.Print(f.Name()+"\n")
			if s[1] == "jpg" || s[1] == "JPG" {
				format = "jpg"
			} else if s[1] == "png" || s[1] == "PNG" {
				format = "png"
			} else {
				continue
			}
			pair := Pair{f.Name(), format}
			if i%4 == 0 {
				l1.PushBack(pair)
				// fmt.Print(l1.Front().Value.(Pair).Name)
			} else if i%4 == 1 {
				l2.PushBack(pair)
			} else if i%4 == 2 {
				l3.PushBack(pair)
			} else if i%4 == 3 {
				l4.PushBack(pair)
			}
			i++

		} else {
			continue
		}

 //    }

	// for q := 1; q < 11 ; q++ {
		// fmt.Print((string)(q)+".jpg")
		// reader, err := os.Open(rootPath+f.Name()) // os.Open(strconv.Itoa(q)+".jpg")
		// if err != nil {
		//     log.Fatal(err)
		// }
		// defer reader.Close()

		// m, _, err := image.Decode(reader)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// bounds := m.Bounds()

		// var red = make([]uint32, bounds.Max.Y * bounds.Max.X)
		// var green = make([]uint32, bounds.Max.Y * bounds.Max.X)
		// var blue = make([]uint32, bounds.Max.Y * bounds.Max.X)
		// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// 	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		// 		r, g, b, _ := m.At(x, y).RGBA()
		// 		// fmt.Printf("%d %d %d\n", r, g, b)
		// 		red[x + y*bounds.Max.X] = r
		// 		green[x + y*bounds.Max.X] = g
		// 		blue[x + y*bounds.Max.X] = b
		// 		// fmt.Print(r , "\n")
		// 		// fmt.Print(float64(r) * 0.299)
		// 	}
		// }

		// // fmt.Print(bounds.Min.Y, " ",  bounds.Max.Y)

		// args := &Args{red, green, blue, bounds.Max.Y * bounds.Max.X}

		// // fmt.Print(args.Length)

		// var replay = make([]uint8, bounds.Max.Y * bounds.Max.X)
		// err = client.Call("Image.ToGrayscale", args, &replay)
		// if err != nil {
		// 	log.Fatal("arith error:", err)
		// } else {
		// 	dest := image.NewGray( image.Rect(0,0,bounds.Max.X, bounds.Max.Y))
				


		// 	for y := 0; y < bounds.Max.Y; y++ {

		// 		for x := 0; x < bounds.Max.X; x++ {
		// 			// fmt.Printf("%d", replay[y])
		// 			newcolor := color.Gray{replay[x + y*bounds.Max.X]}
		// 			// fmt.Print(math.Sqrt((float64) (replay[x + y*bounds.Max.X])), " ", replay[x + y*bounds.Max.X], "\n")
		// 			dest.SetGray(x, y, newcolor)
		// 		}

				
		// 	}

		// 	toimg, _ := os.Create(rootPath+"BW2/BW-"+f.Name()) // os.Create("BW-"+strconv.Itoa(q)+".jpg")
		// 	defer toimg.Close()

		// 	if format == "jpg" {
		// 		jpeg.Encode(toimg, dest, &jpeg.Options{jpeg.DefaultQuality})
		// 	} else if format == "png" {
		// 		png.Encode(toimg, dest)
		// 	}
		// }
	}

	var wg sync.WaitGroup
	wg.Add(4)

	out1 := make(chan string)
	out2 := make(chan string)
	out3 := make(chan string)
	out4 := make(chan string)

	go func() {
		out1 <- handleList(l1, "1", "10.151.12.202:6060", &wg)
	}()
	go func() {
		out2 <- handleList(l2, "2", "10.151.12.202:6060", &wg)
	}()
	go func() {
		out3 <- handleList(l3, "3", "10.151.12.202:6060", &wg)
	}()
	go func() {
		out4 <- handleList(l4, "4", "10.151.12.201:6060", &wg)
	}()

	fmt.Print(<-out1 + " " + <-out2 + " " + <-out3 + " " + <-out4 + "\n")

	// for true {

	// }

	// timeEnd := time.Now()

	// fmt.Printf("%v\n", timeEnd.Sub(timeStart))

	// args := &Args{7,8}
	// var reply int
	// err = client.Call("Calculator.Add", args, &reply)
	// if err != nil {
	// 	log.Fatal("arith error:", err)
	// }
	// fmt.Printf("Result: %d+%d=%d", args.X, args.Y, reply)
}
