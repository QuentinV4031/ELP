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
        // Lire et parser le fichier JSON avec les cartes
        const data = fs.readFileSync('/Users/Quentin/Elm/JS/words.json', 'utf-8');
        this.deck = JSON.parse(data);

        if (!Array.isArray(this.deck) || this.deck.length === 0) {
            throw new Error("Le deck est vide ou mal formaté.");
        }
        
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
        
        // Phase 4: Deviner le mot
        console.log("\nIndices valides:", validClues.join(", "));
        const guess = await this.askQuestion(`${activePlayer}, quel est le mot mystère ? `);
        
        // Gestion des résultats
        this.handleGuessResult(guess.toLowerCase(), validClues);
        
        // Sauvegarde des données
        this.saveRoundData({
            mysteryWord: this.mysteryWord,
            clues,
            validClues,
            guess,
            result: guess.toLowerCase() === this.mysteryWord.toLowerCase()
        });
        
        await this.playRound();
    }

    async collectClues(activePlayer) {
        const clues = [];
        for(const player of this.players.filter(p => p !== activePlayer)) {
            const clue = await this.askQuestion(`${player}, donnez votre indice (utilisez seulement des lettres.) : `);
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
        // Implémentation simplifiée (à améliorer)
        const roots = ["prince", "roi", "chat", "jou"];
        return roots.some(root => clue.includes(root));
    }

    handleGuessResult(guess, validClues) {
        if(guess === this.mysteryWord.toLowerCase()) {
            this.score++;
            console.log(`Correct ! Score: ${this.score}`);
            this.discarded.push(this.currentCard);
        } else {
            console.log(`Incorrect. Le mot était: ${this.mysteryWord}`);
            // Retirer une carte supplémentaire en cas d'échec
            if(this.deck.length > 0) this.discarded.push(this.deck.pop());
        }
    }

    async nextRound() {
        this.activePlayerIndex = (this.activePlayerIndex + 1) % this.players.length;
        await this.playRound();
    }

    endGame() {
        console.log("\n=== FIN DE LA PARTIE ===");
        console.log(`Score final: ${this.score} cartes réussies`);
        console.log(this.getScoreMessage());
        rl.close();
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
        return new Promise(resolve => {
            const askAgain = () => {
                rl.question(question, (answer) => {
                    if (answer.includes(' ')) {
                        console.log("Erreur : votre réponse ne doit pas contenir d'espace. Essayez encore.");
                        askAgain();
                    } else {
                        resolve(answer);
                    }
                });
            };
            askAgain();
        });
    }    
}

// Lancer le jeu
new JustOneGame().initialize();