package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// Connexion au serveur TCP à l'adresse correcte (127.0.0.1:8080)
	conn, err := net.Dial("tcp", ":8080") // Utilisation de l'adresse IPv4
	if err != nil {
		log.Fatalf("Erreur de connexion au serveur : %v", err)
	}
	defer conn.Close()

	// Demander à l'utilisateur l'opération et la valeur
	var operation string
	var value int
	fmt.Print("Entrez l'opération (contrast ou quality): ")
	fmt.Scanln(&operation)
	fmt.Print("Entrez la valeur : ")
	fmt.Scanln(&value)

	// Vérifier que l'opération est valide
	if operation != "contrast" && operation != "quality" {
		log.Fatal("Opération invalide. Veuillez entrer 'contrast' ou 'quality'.")
	}

	// Lire l'image à envoyer
	imagePath := "input_image.jpg" // Changez le chemin de l'image selon votre fichier
	imageData, err := readImageFromFile(imagePath)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture de l'image : %v", err)
	}

	// Envoyer l'opération, la valeur et l'image au serveur
	err = sendRequest(conn, operation, value, imageData)
	if err != nil {
		log.Fatalf("Erreur lors de l'envoi de la requête au serveur : %v", err)
	}

	// Recevoir l'image traitée du serveur
	err = receiveAndSaveImage(conn)
	if err != nil {
		log.Fatalf("Erreur lors de la réception de l'image traitée : %v", err)
	}

	fmt.Println("Image traitée sauvegardée avec succès.")
}

// readImageFromFile lit le fichier image et retourne les données sous forme de slice de bytes
func readImageFromFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'ouverture de l'image : %v", err)
	}
	defer file.Close()

	// Lire l'image entière dans un buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de l'image : %v", err)
	}

	return buf.Bytes(), nil
}

// sendRequest envoie l'opération, la valeur et l'image au serveur
func sendRequest(conn net.Conn, operation string, value int, imageData []byte) error {
	// Envoyer l'opération
	_, err := fmt.Fprintf(conn, "%s\n", operation)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de l'opération : %v", err)
	}

	// Envoyer la valeur
	_, err = fmt.Fprintf(conn, "%d\n", value)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la valeur : %v", err)
	}

	// Envoyer la taille de l'image
	imageSize := int32(len(imageData))
	err = binary.Write(conn, binary.LittleEndian, imageSize)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la taille de l'image : %v", err)
	}

	// Envoyer l'image
	_, err = conn.Write(imageData)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de l'image : %v", err)
	}

	return nil
}

// receiveAndSaveImage reçoit l'image traitée du serveur et la sauvegarde localement
func receiveAndSaveImage(conn net.Conn) error {
	// Lire la taille de l'image
	var imageSize int32
	err := binary.Read(conn, binary.LittleEndian, &imageSize)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture de la taille de l'image : %v", err)
	}

	// Lire l'image traitée
	imageData := make([]byte, imageSize)
	_, err = io.ReadFull(conn, imageData)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture de l'image traitée : %v", err)
	}

	// Sauvegarder l'image traitée dans un fichier
	outputPath := "output_image.jpg"
	err = os.WriteFile(outputPath, imageData, 0644)
	if err != nil {
		return fmt.Errorf("erreur lors de la sauvegarde de l'image traitée : %v", err)
	}

	return nil
}
