package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Erreur lors de la création du serveur TCP : %v", err)
	}
	defer listener.Close()

	log.Println("Serveur TCP démarré sur :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erreur lors de l'acceptation de la connexion : %v", err)
			continue
		}

		// Chaque connexion est traitée dans une goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("Connexion acceptée :", conn.RemoteAddr())

	// Lire l'opération et la valeur
	operation, value, err := readOperationAndValue(conn)
	if err != nil {
		log.Printf("Erreur lors de la lecture de l'opération et de la valeur : %v", err)
		return
	}

	// Lire l'image
	imageData, err := readImage(conn)
	if err != nil {
		log.Printf("Erreur lors de la lecture de l'image : %v", err)
		return
	}

	// Traitement de l'image
	result, err := processImage(operation, value, imageData)
	if err != nil {
		log.Printf("Erreur lors du traitement de l'image : %v", err)
		return
	}

	// Sauvegarder l'image traitée et envoyer la réponse
	err = sendProcessedImage(conn, result)
	if err != nil {
		log.Printf("Erreur lors de l'envoi de l'image traitée : %v", err)
	}
}

func readOperationAndValue(conn net.Conn) (string, int, error) {
	reader := bufio.NewReader(conn)

	// Lire l'opération
	operation, err := reader.ReadString('\n')
	if err != nil {
		return "", 0, err
	}
	operation = operation[:len(operation)-1] // Supprimer le saut de ligne

	// Lire la valeur
	valueStr, err := reader.ReadString('\n')
	if err != nil {
		return "", 0, err
	}
	valueStr = valueStr[:len(valueStr)-1] // Supprimer le saut de ligne

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return "", 0, err
	}

	return operation, value, nil
}

func readImage(conn net.Conn) ([]byte, error) {
	// Lire la taille de l'image
	var imageSize int32
	err := binary.Read(conn, binary.LittleEndian, &imageSize)
	if err != nil {
		return nil, err
	}

	// Lire l'image
	imageData := make([]byte, imageSize)
	_, err = io.ReadFull(conn, imageData)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func processImage(operation string, value int, imageData []byte) (*image.NRGBA, error) {
	src, err := imaging.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	var result *image.NRGBA
	if operation == "contrast" {
		result = imaging.AdjustContrast(src, float64(value))
	} else if operation == "quality" {
		result = imaging.Resize(src, src.Bounds().Dx()*value/100, 0, imaging.Lanczos)
	} else {
		return nil, fmt.Errorf("Opération inconnue : %s", operation)
	}

	return result, nil
}

func sendProcessedImage(conn net.Conn, result *image.NRGBA) error {
	// Sauvegarder temporairement l'image traitée
	outputPath := "output_image.jpg"
	err := imaging.Save(result, outputPath)
	if err != nil {
		return err
	}

	// Charger l'image sauvegardée
	outputFile, err := os.Open(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	outputInfo, err := outputFile.Stat()
	if err != nil {
		return err
	}

	outputSize := outputInfo.Size()
	outputData := make([]byte, outputSize)
	_, err = outputFile.Read(outputData)
	if err != nil {
		return err
	}

	// Envoyer la taille de l'image
	err = binary.Write(conn, binary.LittleEndian, int32(outputSize))
	if err != nil {
		return err
	}

	// Envoyer l'image traitée
	_, err = conn.Write(outputData)
	if err != nil {
		return err
	}

	log.Println("Image traitée envoyée avec succès.")
	return nil
}
