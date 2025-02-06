const fs = require('fs');
const readline = require('readline');

// Liste de mots à deviner
const mots = ["pomme", "banane", "ordinateur", "chaise", "soleil", "livre", "montagne"];

// Fonction pour choisir un mot aléatoire
function choisirMot() {
    return mots[Math.floor(Math.random() * mots.length)];
}

// Fonction pour jouer à Just One
async function jouerJustOne() {
    const motSecret = choisirMot();
    console.log(`Le mot secret est : ${motSecret} (à ne pas révéler aux joueurs !)`);

    const joueurs = ["Alice", "Bob", "Charlie", "Diana", "Eve"];
    const indices = [];

    // Interface pour lire les entrées des joueurs
    const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout,
    });

    // Fonction pour demander un indice à un joueur
    function demanderIndice(joueur) {
        return new Promise((resolve) => {
            rl.question(`${joueur}, propose un indice pour le mot secret : `, (indice) => {
                indices.push({ joueur, indice: indice.toLowerCase() });
                resolve();
            });
        });
    }

    // Demander un indice à chaque joueur
    for (const joueur of joueurs) {
        await demanderIndice(joueur);
    }

    // Vérifier les doublons et les annuler
    const indicesUniques = [];
    const doublons = new Set();

    indices.forEach(({ indice }, index) => {
        if (indices.slice(0, index).some((i) => i.indice === indice)) {
            doublons.add(indice);
        }
    });

    indices.forEach(({ joueur, indice }) => {
        if (!doublons.has(indice)) {
            indicesUniques.push({ joueur, indice });
        }
    });

    // Afficher les indices restants
    console.log("\nIndices restants :");
    indicesUniques.forEach(({ joueur, indice }) => {
        console.log(`${joueur} : ${indice}`);
    });

    // Demander au devineur de deviner le mot
    rl.question("\nDevineur, quel est le mot secret ? ", (tentative) => {
        if (tentative.toLowerCase() === motSecret) {
            console.log("Bravo ! Vous avez trouvé le mot secret !");
        } else {
            console.log(`Dommage, le mot secret était : ${motSecret}`);
        }

        // Enregistrer les propositions dans un fichier
        const data = JSON.stringify({
            motSecret,
            indices,
            indicesUniques,
            tentative,
            resultat: tentative.toLowerCase() === motSecret ? "Gagné" : "Perdu",
        }, null, 2);

        fs.writeFileSync('propositions.json', data);
        console.log("Les propositions ont été enregistrées dans 'propositions.json'.");

        rl.close();
    });
}

// Lancer le jeu
jouerJustOne();