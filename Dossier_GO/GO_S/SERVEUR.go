package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/image/draw"
)

const (
	workerPoolSize = 4
	tcpPort        = ":8080"
)

func main() {
	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s\n", tcpPort)

	jobs := make(chan net.Conn)
	var wg sync.WaitGroup

	// Initialize worker pool
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	// Accept connections and dispatch to workers
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		jobs <- conn
	}

	wg.Wait()
}

func worker(jobs <-chan net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for conn := range jobs {
		handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	// Read command
	cmd, err := readCommand(conn)
	if err != nil {
		fmt.Println("Error reading command:", err)
		return
	}

	// Read image data
	img, format, err := readImage(conn)
	if err != nil {
		fmt.Println("Error reading image:", err)
		return
	}

	// Process image
	processedImg, err := processImage(img, cmd)
	if err != nil {
		fmt.Println("Error processing image:", err)
		return
	}

	// Send processed image
	if err := sendImage(conn, processedImg, format); err != nil {
		fmt.Println("Error sending image:", err)
		return
	}
}

func readCommand(conn net.Conn) (string, error) {
	var cmdLen uint32
	if err := binary.Read(conn, binary.BigEndian, &cmdLen); err != nil {
		return "", err
	}

	cmdBuf := make([]byte, cmdLen)
	if _, err := io.ReadFull(conn, cmdBuf); err != nil {
		return "", err
	}

	return string(cmdBuf), nil
}

func readImage(conn net.Conn) (image.Image, string, error) {
	var imgSize uint32
	if err := binary.Read(conn, binary.BigEndian, &imgSize); err != nil {
		return nil, "", err
	}

	imgBuf := make([]byte, imgSize)
	if _, err := io.ReadFull(conn, imgBuf); err != nil {
		return nil, "", err
	}

	img, format, err := image.Decode(bytes.NewReader(imgBuf))
	return img, format, err
}

func sendImage(conn net.Conn, img image.Image, format string) error {
	var buf bytes.Buffer
	switch format {
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			return err
		}
	case "jpeg":
		if err := jpeg.Encode(&buf, img, nil); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	processedBytes := buf.Bytes()
	processedSize := uint32(len(processedBytes))

	if err := binary.Write(conn, binary.BigEndian, processedSize); err != nil {
		return err
	}

	_, err := conn.Write(processedBytes)
	return err
}

func processImage(img image.Image, command string) (image.Image, error) {
	parts := strings.Split(command, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid command format")
	}

	action, params := parts[0], parts[1]

	switch action {
	case "blur":
		radius, err := strconv.Atoi(params)
		if err != nil {
			return nil, fmt.Errorf("invalid blur radius")
		}
		return applyBlur(img, radius), nil

	case "resize":
		dims := strings.Split(params, "x")
		if len(dims) != 2 {
			return nil, fmt.Errorf("invalid resize dimensions")
		}
		width, _ := strconv.Atoi(dims[0])
		height, _ := strconv.Atoi(dims[1])
		return resizeImage(img, width, height), nil

	case "contrast":
		factor, err := strconv.ParseFloat(params, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid contrast factor")
		}
		return adjustContrast(img, factor), nil

	default:
		return nil, fmt.Errorf("unsupported action: %s", action)
	}
}

func applyBlur(img image.Image, radius int) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, a uint32
			count := 0

			// Sample surrounding pixels
			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					xx, yy := x+dx, y+dy
					if xx >= bounds.Min.X && xx < bounds.Max.X && yy >= bounds.Min.Y && yy < bounds.Max.Y {
						pr, pg, pb, pa := img.At(xx, yy).RGBA()
						r += pr
						g += pg
						b += pb
						a += pa
						count++
					}
				}
			}

			// Average the values (no bit-shifting)
			if count > 0 {
				r /= uint32(count)
				g /= uint32(count)
				b /= uint32(count)
				a /= uint32(count)
			}

			// Preserve 16-bit color depth
			dst.Set(x, y, color.RGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: uint16(a),
			})
		}
	}
	return dst
}

func resizeImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

func adjustContrast(img image.Image, factor float64) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rf := float64(r)/65535.0 - 0.5
			gf := float64(g)/65535.0 - 0.5
			bf := float64(b)/65535.0 - 0.5

			rf = clamp((rf*factor)+0.5) * 65535.0
			gf = clamp((gf*factor)+0.5) * 65535.0
			bf = clamp((bf*factor)+0.5) * 65535.0

			dst.Set(x, y, color.RGBA64{
				R: uint16(rf),
				G: uint16(gf),
				B: uint16(bf),
				A: uint16(a),
			})
		}
	}
	return dst
}

func clamp(v float64) float64 {
	if v < 0.0 {
		return 0.0
	}
	if v > 1.0 {
		return 1.0
	}
	return v
}
