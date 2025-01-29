package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <8080> <applyBlur> <C:/Users/jerem/Desktop/3TC/ELP/GO_C/input_image>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	serverAddr := os.Args[1]
	command := os.Args[2]
	imagePath := os.Args[3]

	// Read and encode image
	img, format, err := readAndEncodeImage(imagePath)
	if err != nil {
		fmt.Println("Error processing image:", err)
		os.Exit(1)
	}

	// Connect to server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Connection error:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Send command and image
	if err := sendRequest(conn, command, img); err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}

	// Receive and save processed image
	if err := receiveAndSaveImage(conn, format); err != nil {
		fmt.Println("Error receiving image:", err)
		os.Exit(1)
	}

	fmt.Println("Image processed successfully")
}

func readAndEncodeImage(path string) ([]byte, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	switch format {
	case "png":
		err = png.Encode(&buf, img)
	case "jpeg":
		err = jpeg.Encode(&buf, img, nil)
	default:
		return nil, format, fmt.Errorf("unsupported format: %s", format)
	}

	return buf.Bytes(), format, err
}

func sendRequest(conn net.Conn, command string, img []byte) error {
	// Send command
	cmdBytes := []byte(command)
	if err := binary.Write(conn, binary.BigEndian, uint32(len(cmdBytes))); err != nil {
		return err
	}
	if _, err := conn.Write(cmdBytes); err != nil {
		return err
	}

	// Send image
	if err := binary.Write(conn, binary.BigEndian, uint32(len(img))); err != nil {
		return err
	}
	_, err := conn.Write(img)
	return err
}

func receiveAndSaveImage(conn net.Conn, format string) error {
	var imgSize uint32
	if err := binary.Read(conn, binary.BigEndian, &imgSize); err != nil {
		return err
	}

	imgData := make([]byte, imgSize)
	if _, err := io.ReadFull(conn, imgData); err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return err
	}

	outFile, err := os.Create("processed." + format)
	if err != nil {
		return err
	}
	defer outFile.Close()

	switch format {
	case "png":
		return png.Encode(outFile, img)
	case "jpeg":
		return jpeg.Encode(outFile, img, nil)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}
