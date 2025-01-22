package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("Connexion acceptée :", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	// Lire l'opération
	operation, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Erreur lors de la lecture de l'opération : %v", err)
		return
	}
	operation = operation[:len(operation)-1] // Supprimer le saut de ligne

	// Lire la valeur
	valueStr, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Erreur lors de la lecture de la valeur : %v", err)
		return
	}
	valueStr = valueStr[:len(valueStr)-1] // Supprimer le saut de ligne

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Valeur invalide : %v", err)
		return
	}

	// Lire la taille de l'image
	var imageSize int32
	err = binary.Read(reader, binary.LittleEndian, &imageSize)
	if err != nil {
		log.Printf("Erreur lors de la lecture de la taille de l'image : %v", err)
		return
	}

	// Lire l'image
	imageData := make([]byte, imageSize)
	_, err = io.ReadFull(reader, imageData)
	if err != nil {
		log.Printf("Erreur lors de la lecture des données de l'image : %v", err)
		return
	}

	src, err := imaging.Decode(bytes.NewReader(imageData))
	if err != nil {
		log.Printf("Erreur lors du décodage de l'image : %v", err)
		return
	}

	var result *image.NRGBA
	if operation == "contrast" {
		result = imaging.AdjustContrast(src, float64(value))
	} else if operation == "quality" {
		result = imaging.Resize(src, src.Bounds().Dx()*value/100, 0, imaging.Lanczos)
	} else {
		log.Printf("Opération inconnue : %s", operation)
		return
	}

	// Sauvegarder l'image temporairement
	outputPath := "output_image.jpg"
	err = imaging.Save(result, outputPath)
	if err != nil {
		log.Printf("Erreur lors de la sauvegarde de l'image : %v", err)
		return
	}

	// Charger l'image sauvegardée
	outputFile, err := os.Open(outputPath)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture du fichier traité : %v", err)
		return
	}
	defer outputFile.Close()

	outputInfo, err := outputFile.Stat()
	if err != nil {
		log.Printf("Erreur lors de la récupération des informations du fichier : %v", err)
		return
	}

	outputSize := outputInfo.Size()
	outputData := make([]byte, outputSize)
	_, err = outputFile.Read(outputData)
	if err != nil {
		log.Printf("Erreur lors de la lecture du fichier traité : %v", err)
		return
	}

	// Envoyer la taille de l'image
	err = binary.Write(conn, binary.LittleEndian, int32(outputSize))
	if err != nil {
		log.Printf("Erreur lors de l'envoi de la taille de l'image : %v", err)
		return
	}

	// Envoyer l'image traitée
	_, err = conn.Write(outputData)
	if err != nil {
		log.Printf("Erreur lors de l'envoi de l'image : %v", err)
		return
	}

	log.Println("Image traitée envoyée avec succès.")
}
