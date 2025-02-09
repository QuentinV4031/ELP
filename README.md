# Image Processing Server and Client
Ce projet consiste en un serveur et un client en Go qui permettent de traiter des images en fonction de commandes spécifiques. Le serveur reçoit une image et une commande de traitement, effectue l'opération demandée, et renvoie l'image traitée au client. Le client envoie une image et une commande au serveur, puis sauvegarde l'image traitée.

## Fonctionnalités
Le serveur supporte les opérations suivantes :

Flou (Blur) : Applique un flou à l'image avec un rayon spécifié.

Redimensionnement (Resize) : Redimensionne l'image aux dimensions spécifiées.

Contraste (Contrast) : Ajuste le contraste de l'image avec un facteur spécifié.

### Prérequis
Go 1.16 ou supérieur installé sur votre machine.

Une image au format PNG ou JPEG.

### Installation
Clonez le dépôt contenant les fichiers CLIENT.go et SERVEUR.go.

Ouvrez un terminal et naviguez jusqu'au répertoire du projet.

### Utilisation
#### 1. Compilation
Compilez les fichiers CLIENT.go et SERVEUR.go :

go build -o client CLIENT.go
go build -o serveur SERVEUR.go


#### 2. Lancement du Serveur
Lancez le serveur en exécutant le fichier compilé serveur :

./serveur
Le serveur écoutera sur le port 8080 par défaut.

#### 3. Utilisation du Client
Le client prend trois arguments :

L'adresse du serveur (par exemple, localhost:8080).

La commande de traitement.

Le chemin de l'image à traiter.

Exemples de commandes
Flou : Applique un flou avec un rayon de 5.

./client localhost:8080 "blur:5" image.jpg
Redimensionnement : Redimensionne l'image à 800x600 pixels.

./client localhost:8080 "resize:800x600" image.png
Contraste : Ajuste le contraste avec un facteur de 1.5.

./client localhost:8080 "contrast:1.5" image.jpg
#### 4. Résultat
Après exécution, le client sauvegarde l'image traitée sous le nom processed.<format> dans le répertoire courant, où <format> est le format d'origine de l'image (par exemple, processed.jpg ou processed.png).

## Structure du Code
CLIENT.go : Contient le code du client qui envoie une image et une commande au serveur, puis sauvegarde l'image traitée.
  
SERVEUR.go : Contient le code du serveur qui reçoit une image et une commande, traite l'image, et renvoie le résultat au client.

### Fonctions Principales
readAndEncodeImage : Lit et encode l'image dans le format approprié.

sendRequest : Envoie la commande et l'image au serveur.

receiveAndSaveImage : Reçoit et sauvegarde l'image traitée.

processImage : Traite l'image en fonction de la commande reçue.

applyBlur, resizeImage, adjustContrast : Fonctions de traitement d'image.

### Remarques
Le serveur supporte un pool de workers pour gérer plusieurs connexions simultanément.

Les formats d'image supportés sont PNG et JPEG
