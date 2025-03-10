# ELP

Le fichier de code est " justOne_v5.js "  (nous n'avions pas encore bien compris l'utilité de Git  mais c'est le cas maintenant)

Partie javascript

lancement du jeu:
taper la commande
npm start

la partie se lance
le jeu commence automatiquement le premier tour, 
il présente les mots d'une carte et demande de choisir l'indice du mot voulu
il demande tour à tour les indices des autres joueurs puis demande au joueur sélectionné de deviner le mot
on passe à la manche suivante

format de sauvegarde des parties:

{
"mysteryWord": "voiture",
  "clues": [{"player": "Ivan","clue": "auto"}, {"joueur2":"","clue":""} ... ],
  "validClues": ["auto", "clue2", ... ],
  "guess": "voiture",
  "result": true
}
