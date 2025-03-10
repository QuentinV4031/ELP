const fs = require('fs');
const readline = require('readline');
const rl = readline.createInterface({ input: process.stdin, output: process.stdout });

class JustOneGame {
    constructor() {
        this.players = ["Louis", "Ivan", "Jérémie", "Karel", "Hatim"];
        this.score = 0;
        this.deck = [];
        this.discarded = [];
        this.currentCard = null;
        this.activePlayerIndex = 0;
    }

    async initialize() {
        // Charger 13 cartes avec 5 mots chacune (exemple simplifié)
        this.deck = Array.from({length: 1}, () => ({
            words: ["Europe", "Cirque", "Virus", "Crocodile", "Moutarde"]
        }));
        
        console.log("=== NOUVELLE PARTIE DE JUST ONE ===");
        await this.playRound();
    }

    async playRound() {
        if(this.deck.length === 0) return this.endGame();
        
        // Phase 1: Sélection du mot mystère
        this.currentCard = this.deck.pop();
        const activePlayer = this.players[this.activePlayerIndex];
        
        console.log(`\n--- Tour de ${activePlayer} ---`);
        console.log("Mots disponibles:", this.currentCard.words.join(", "));
        
        const wordIndex = await this.askQuestion(`${activePlayer}, choisissez un numéro entre 1 et 5 : `);
        this.mysteryWord = this.currentCard.words[wordIndex - 1];
        
        // Phase 2: Collecte des indices
        const clues = await this.collectClues(activePlayer);
        
        // Phase 3: Validation des indices
        const { validClues, invalidClues } = this.validateClues(clues);
        
        if(validClues.length === 0) {
            console.log("Tous les indices sont invalides ! Carte annulée.");
            this.discarded.push(this.currentCard);
            return this.nextRound();
        }
        
        // Phase 4: Deviner le mot avec utilisation des indices
        await this.guessWithClues(activePlayer, validClues);
        
        // Sauvegarde des données
        this.saveRoundData({
            mysteryWord: this.mysteryWord,
            clues,
            validClues,
            result: this.score
        });
        
        await this.nextRound();
    }

    async guessWithClues(activePlayer, validClues) {
        let remainingClues = [...validClues];
        
        while(remainingClues.length > 0) {
            console.log(`\nIndices restants (${remainingClues.length}) : ${remainingClues.join(", ")}`);
            const guess = await this.askQuestion(`${activePlayer}, quel est le mot mystère ? `);
            
            if(guess.toLowerCase() === this.mysteryWord.toLowerCase()) {
                this.score++;
                console.log(`Correct ! Score: ${this.score}`);
                this.discarded.push(this.currentCard);
                return;
            } else {
                console.log("Incorrect.\n");
                console.log("Voici un indice supplémentaire: ");
                const usedClue = remainingClues.shift(); // Retire le premier indice
                console.log(`Indice utilisé : ${usedClue}`);
            }
        }
        
        // Si tous les indices ont été utilisés sans succès
        console.log(`\nDommage ! Le mot mystère était : ${this.mysteryWord}`);
        this.discarded.push(this.currentCard);
    }

    async collectClues(activePlayer) {
        const clues = [];
        for(const player of this.players.filter(p => p !== activePlayer)) {
            const clue = await this.askQuestion(`${player}, donnez votre indice : `);
            clues.push({ player, clue: clue.toLowerCase() });
        }
        return clues;
    }

    validateClues(clues) {
        const invalid = new Set();
        const duplicates = new Set();
        const clueCounts = {};

        // Détection des doublons et validation
        clues.forEach(({clue}) => {
            clueCounts[clue] = (clueCounts[clue] || 0) + 1;
            
            // Vérification des règles d'invalidité
            if(this.isInvalidClue(clue)) invalid.add(clue);
        });

        // Marquer les doublons
        Object.entries(clueCounts).forEach(([clue, count]) => {
            if(count > 1) duplicates.add(clue);
        });

        const validClues = clues
            .filter(({clue}) => 
                !invalid.has(clue) && 
                !duplicates.has(clue) &&
                !this.isSameWordFamily(clue)
            )
            .map(c => c.clue);

        return { validClues, invalidClues: [...invalid, ...duplicates] };
    }

    isInvalidClue(clue) {
        const lowerMystery = this.mysteryWord.toLowerCase();
        return (
            clue === lowerMystery
        );
    }

    isSameWordFamily(clue) {
        const roots = ["prince", "roi", "chat", "jou"];
        return roots.some(root => clue.includes(root));
    }

    async nextRound() {
        this.activePlayerIndex = (this.activePlayerIndex + 1) % this.players.length;
        await this.playRound();
    }

    async endGame() {
        console.log("\n=== FIN DE LA PARTIE ===");
        console.log(`Score final: ${this.score} cartes réussies`);
        console.log(this.getScoreMessage());
    
        const replay = await this.askQuestion("Voulez-vous rejouer ? (Oui/Non) : ");
        if (replay.toLowerCase() === "oui") {
            this.score = 0; // Réinitialiser le score
            this.deck = Array.from({length: 13}, () => ({
                words: ["Europe", "Cirque", "Virus", "Crocodile", "Moutarde"]
            }));
            this.activePlayerIndex = 0; // Réinitialiser le joueur actif
            await this.playRound(); // Relancer une nouvelle partie
        } else {
            rl.close(); // Fermer le jeu
        }
    }

    getScoreMessage() {
        const messages = {
            13: "Score parfait !",
            12: "Incroyable !",
            11: "Génial !",
            10: "Waouh !",
            9: "Pas mal !",
            8: "Dans la moyenne",
            7: "Peut mieux faire",
            6: "C'est un bon début",
        };
        return messages[this.score] || "Essayez encore !";
    }

    saveRoundData(data) {
        fs.appendFileSync('game_log.json', JSON.stringify(data, null, 2) + ',\n');
    }

    async askQuestion(question) {
        return new Promise(resolve => rl.question(question, resolve));
    }
}

// Lancer le jeu
new JustOneGame().initialize();