package hangman

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/nsf/termbox-go"
)

type HangManData struct {
	Word             []rune     // Word composed of '_', ex: H_ll_
	ToFind           []rune     // Final word chosen by the program at the beginning. It is the word to find
	Attempts         int        // Number of attempts left
	HangmanPositions [10]string // It can be the array where the positions parsed in "hangman.txt" are stored
	History          string
	Difficulty       string
}

func RandomPickLine() HangManData { // permet de choisir un mot aléatoire à trouver dans les différents fichiers words.txt
	var diff string
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	var game HangManData
	Clear()
	fmt.Print("Difficulté : facile/normal/difficile (si autre chose, tous les mots des 3 difficultés seront utilisés)\n\n")
	fmt.Scan(&diff)
	if diff == "facile" { // permet de sélectionner la difficulté
		game.Difficulty = diff
		diff = "Ressources/words1.txt"
	} else if diff == "normal" {
		game.Difficulty = diff
		diff = "Ressources/words2.txt"
	} else if diff == "difficile" {
		game.Difficulty = diff
		diff = "Ressources/words3.txt"
	} else { // si pas de difficulté particulière sélectionnée, joue avec tous les mots dee tous les fichiers
		game.Difficulty = "Tous les mots"
		diff = "Ressources/words.txt"
	}
	Clear()
	// permet de faire un message variable selon la difficulté
	Printf_tb(20, 7, termbox.ColorWhite, termbox.ColorDefault, "Vous avez choisi la difficulté :")
	Printf_tb(35, 8, termbox.ColorWhite, termbox.ColorDefault, "%s", game.Difficulty)
	file, err := os.Open(diff)
	var lines []string
	if err != nil {
		log.Fatal(err)
	}
	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanWords)
	for Scanner.Scan() { // on lit et on append dans un tableau de string à chaque ligne du fichier
		lines = append(lines, Scanner.Text())
	}
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	randomnum := r1.Intn(len(lines))
	game.ToFind = []rune(lines[randomnum]) // on prend une string dans le tableau de string de manière random et on dit que c'est le mot qu'on doit trouver
	return game
}

func PrintHanged(game HangManData) { // permet de print les différentes postions du pendu
	file, err := os.Open("Ressources/hangman.txt")
	var lines []string
	if err != nil {
		log.Fatal(err)
	}
	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanLines)
	for Scanner.Scan() {
		lines = append(lines, Scanner.Text())
	}
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	nbfautes := 10 - game.Attempts // CONTINUER A TRAVAILLER
	x := nbfautes*8 - 8
	y := x + 6
	for i := x; i <= y && y <= 80; i++ {
		if nbfautes < 1 {
			break
		}
		Printf_tb(66, 14+i-x, termbox.ColorWhite, termbox.ColorDefault, lines[i])
	}
}

func PrintGameWord(game HangManData) { // permet de print dans Termbox (avec espace entre les lettres) game.Word (mot qu'on est en train d'essayer d'afficher en entier pour gagner)
	var modtab []rune
	for i := 0; i < len(game.Word); i++ {
		modtab = append(modtab, game.Word[i], ' ')
	}
	for i := 0; i < len(modtab); i++ {
		Printf_tb(3+i, 14, termbox.ColorWhite, termbox.ColorDefault, "%s", ToHigher(string(modtab[i])))
	}
}

func ModifyGameWord(game HangManData, letter string) HangManData { // permet de changer les _ du game.Word en lettres de toFind correspondantesS
	var indextofind []int
	for i := 0; i < len(game.ToFind); i++ {
		if rune(letter[0]) == game.ToFind[i] { // tableau qui stock les index où apparaît letter dans ToFind
			indextofind = append(indextofind, i)
		}
	}
	for j := 0; j < len(indextofind); j++ {
		game.Word[indextofind[j]] = game.ToFind[indextofind[j]]
	}
	return game
}

var clear map[string]func() //create a map for storing clear funcs

func init() { // permet d'initialiser la fonction clear
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func Clear() { // permet de clear le terminal
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Le terminal n'est pas supporté.")
	}
}

func RandomPickLetter(game HangManData) []int { // permet d'avoir 2 index afin de choisir randP index aléatoires à remplacer par les bonnes lettres dans game.Word
	var randP []int
	var blacklist []int
	n := len(game.ToFind)/2 - 1
	for i := 0; i < n; i++ {
		nbUtils := RandomBlacklist(len(game.ToFind), blacklist)
		blacklist = append(blacklist, nbUtils)
		randP = append(randP, nbUtils)
	}
	return randP
}

func RandomBlacklist(max int, blacklisted []int) int { // permet de ne pas avoir 2 fois le même index renvoyé dans le tableau d'int de RandomPickLetter
	excluded := map[int]bool{}
	for _, x := range blacklisted {
		excluded[x] = true
	}
	for {
		n := rand.Intn(max)
		if !excluded[n] {
			return n
		}
	}
}

func PrintHistory_tb(game HangManData) { // Permet de Print les lettres testées fausses de termbox et de sauter une ligne s'il y en a trop de fausses (même si on a que 10 vies je sais)
	for i := 0; i < len(game.History); i++ {
		if i+32 >= 58 {
			if i+6 >= 58 {
				break
			} else {
				Printf_tb(i+6, 15, termbox.ColorWhite, termbox.ColorDefault, ToHigher(string(game.History[i])))
			}
		} else {
			Printf_tb(32+i, 14, termbox.ColorWhite, termbox.ColorDefault, ToHigher(string(game.History[i])))
		}
	}
}

func SaveGame(game HangManData) { // Permet de sauvegarder les valeurs de la structure dans un save.txt encrypté en json
	file, _ := json.Marshal(game)
	_ = ioutil.WriteFile("save.txt", file, 0644)
}

func LoadGame() HangManData { // Permet de return un HangmanData qui contient les valeurs encryptées dans save.txt
	file, err := os.Open("save.txt")
	var lines []string
	if err != nil {
		log.Fatal(err)
	}
	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanLines)
	for Scanner.Scan() {
		lines = append(lines, Scanner.Text())
	}
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	str := lines[0]
	game := HangManData{}
	json.Unmarshal([]byte(str), &game)
	return game
}

func Print_tb(x, y int, fg, bg termbox.Attribute, msg string) { // permet de print une box
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func Printf_tb(x, y int, fg, bg termbox.Attribute, format string, args ...interface{}) { // permet de print une box de manière plus poussée
	s := fmt.Sprintf(format, args...)
	Print_tb(x, y, fg, bg, s)
}

var borderTopLeft rune = 0x250C // Termbox utiles très utiles (caractères pour dessiner les boxx)
var borderTopRight rune = 0x2510
var borderBotomLeft rune = 0x2514
var borderBottomRight rune = 0x2518
var borderHorizontal rune = 0x2500
var borderVertical rune = 0x2502

// CETTE FONCTION PERMET DE DEFINIR UN CARRE TERMBOX : on demande l'abscisse du point en haut à gauche du rectangle, la coordonnée x à ne pas dépasser pour la longueur, la longueur en y, la couleur des caractères, le surlignage, le titre de la box, la couleur du titre et enfin le surlignage du titre
func Print_termbox_square(xupleft, yupleft, xright, lengthy int, fg termbox.Attribute, bg termbox.Attribute, s string, couleurS termbox.Attribute, bg2 termbox.Attribute) {
	termbox.SetCell(xupleft, yupleft, borderTopLeft, fg, bg) // Ces 4 lignes définissent les caractèresspéciaux à afficher aux coins du rectangle
	termbox.SetCell(xright, yupleft, borderTopRight, fg, bg)
	termbox.SetCell(xupleft, yupleft+lengthy, borderBotomLeft, fg, bg)
	termbox.SetCell(xright, yupleft+lengthy, borderBottomRight, fg, bg)
	for i := xupleft + 1; i < xright; i++ { // print les côtés horizontaux du rectangle termbox
		termbox.SetCell(i, yupleft, borderHorizontal, fg, bg)
		termbox.SetCell(i, lengthy+yupleft, borderHorizontal, fg, bg)
	}
	for i := yupleft + 1; i < lengthy+yupleft; i++ { // print les côtés verticaux du rectangle termbox
		termbox.SetCell(xupleft, i, borderVertical, fg, bg)
		termbox.SetCell(xright, i, borderVertical, fg, bg)
	}
	Printf_tb(xupleft+4, yupleft, couleurS|termbox.AttrBold, bg2, s)
}

func Print_tb_game_boxes(game HangManData) { // permet de print toute la structure du Termbox
	Print_termbox_square(0, 0, 80, 7, termbox.ColorWhite, termbox.ColorBlack, "", termbox.ColorMagenta, termbox.ColorBlack)
	Printf_tb(35, 2, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault, "HANGMAN !")
	Printf_tb(19, 3, termbox.ColorMagenta, termbox.ColorDefault, "(Appuyez sur ÉCHAP pour quitter la partie)")
	Printf_tb(20, 5, termbox.ColorLightRed|termbox.AttrBold, termbox.ColorDefault, "GAME")
	Printf_tb(34, 5, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault, "DESCRIPTION")
	Printf_tb(56, 5, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault, "AIDE")
	/* DIFFICULTÉ */ Print_termbox_square(0, 8, 80, 3, termbox.ColorRed, termbox.ColorBlack, "DIFFICULTÉ", termbox.ColorLightRed, termbox.ColorBlack)
	/* MOT À DEVINER */ Print_termbox_square(0, 12, 29, 4, termbox.ColorRed, termbox.ColorBlack, "MOT À DEVINER", termbox.ColorWhite, termbox.ColorBlack)
	/* LETTRES UTILISEES */ Print_termbox_square(30, 12, 60, 4, termbox.ColorRed, termbox.ColorBlack, "LETTRES UTILISÉES", termbox.ColorWhite, termbox.ColorBlack)
	/* PENDU */ Print_termbox_square(61, 12, 80, 10, termbox.ColorRed, termbox.ColorBlack, "PENDU", termbox.ColorWhite, termbox.ColorBlack)
	/* ENTREE UTILISATEUR */ Print_termbox_square(0, 17, 60, 5, termbox.ColorWhite, termbox.ColorBlack, "ENTRÉE UTILISATEUR", termbox.ColorLightRed, termbox.ColorBlack)
	Printf_tb(23, 9, termbox.ColorWhite, termbox.ColorDefault, "Vous avez choisi la difficulté :")
	Printf_tb(35, 10, termbox.ColorWhite, termbox.ColorDefault, "%s", game.Difficulty)
	PrintHanged(game)
	PrintGameWord(game)
	PrintHistory_tb(game)
}

func Print_tb_desc_boxes(game HangManData) { // permet de print toute la structure du Termbox
	Print_termbox_square(0, 0, 80, 7, termbox.ColorWhite, termbox.ColorBlack, "", termbox.ColorMagenta, termbox.ColorBlack)
	Printf_tb(35, 2, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault, "HANGMAN !")
	Printf_tb(19, 3, termbox.ColorMagenta, termbox.ColorDefault, "(Appuyez sur ÉCHAP pour quitter la partie)")
	Printf_tb(20, 5, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault, "GAME")
	Printf_tb(34, 5, termbox.ColorLightRed|termbox.AttrBold, termbox.ColorDefault, "DESCRIPTION")
	Printf_tb(56, 5, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault, "AIDE")
	Print_termbox_square(0, 8, 80, 14, termbox.ColorWhite, termbox.ColorBlack, "", termbox.ColorMagenta, termbox.ColorBlack)
	descprint()
}

func descprint() { // Print les string qui décrivent le jeu
	Printf_tb(12, 10, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "Bienvenue ! ")
	Printf_tb(24, 10, termbox.ColorWhite, termbox.ColorDefault, "Vous êtes sur le ")
	Printf_tb(41, 10, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault, "HANGMAN")
	Printf_tb(48, 10, termbox.ColorWhite, termbox.ColorDefault, ", notre jeu de pendu !")
	Printf_tb(12, 12, termbox.ColorWhite, termbox.ColorDefault, "Le but du jeu est de deviner le mot aléatoire. Pour cela,")
	Printf_tb(12, 13, termbox.ColorWhite, termbox.ColorDefault, "il vous faudra entrer les lettres unes par unes ou le mot.")
	Printf_tb(14, 15, termbox.ColorWhite, termbox.ColorDefault, "Mais ")
	Printf_tb(19, 15, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault, "/!\\ ATTENTION /!\\ ")
	Printf_tb(17+20, 15, termbox.ColorWhite, termbox.ColorDefault, "vous n'avez que 10 vies, et une")
	Printf_tb(22, 16, termbox.ColorWhite, termbox.ColorDefault, "fois que Jo est pendu, vous avez perdu.")
	Printf_tb(34, 18, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault, "BONNE CHANCE !")
	Printf_tb(12, 20, termbox.ColorMagenta, termbox.ColorDefault, "(Flèche de gauche/ droite pour naviguer dans les onglets)")
}

func Print_tb_help_boxes(game HangManData) { // permet de print toute la structure du Termbox
	Print_termbox_square(0, 0, 80, 7, termbox.ColorWhite, termbox.ColorBlack, "", termbox.ColorMagenta, termbox.ColorBlack)
	Printf_tb(35, 2, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault, "HANGMAN !")
	Printf_tb(19, 3, termbox.ColorMagenta, termbox.ColorDefault, "(Appuyez sur ÉCHAP pour quitter la partie)")
	Printf_tb(20, 5, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault, "GAME")
	Printf_tb(34, 5, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault, "DESCRIPTION")
	Printf_tb(56, 5, termbox.ColorLightRed|termbox.AttrBold, termbox.ColorDefault, "AIDE")
	Print_termbox_square(0, 8, 80, 14, termbox.ColorWhite, termbox.ColorBlack, "", termbox.ColorMagenta, termbox.ColorBlack)
	helpprint()
}

func helpprint() {
	Printf_tb(17, 10, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "Voici quelques commandes utiles pour le jeu :")
	Printf_tb(8, 13, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "Flèche de gauche < / Flèche de droite > :")
	Printf_tb(49, 13, termbox.ColorWhite, termbox.ColorDefault, " Changer d'onglet")
	Printf_tb(8, 15, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "ÉCHAP :")
	Printf_tb(15, 15, termbox.ColorWhite, termbox.ColorDefault, " Quitter la partie (fait une sauvegarde)")
	Printf_tb(23, 17, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "POUR REPRENDRE OÙ VOUS EN ÉTIEZ : ")
	Printf_tb(18, 18, termbox.ColorWhite, termbox.ColorDefault, " go run main/main.go --startWith save.txt")
	Printf_tb(4, 20, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "Pour tester une lettre, appuyez sur la touche du clavier puis sur ENTRÉE")
}

func PrintWin(game HangManData) { // Print l'écran de victoire
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	file, err := os.Open("Ressources/YouWin.txt")
	var lines []string
	if err != nil {
		log.Fatal(err)
	}
	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanLines)
	for Scanner.Scan() {
		lines = append(lines, Scanner.Text())
	}
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	Print_termbox_square(0, 0, 80, 15, termbox.ColorYellow, termbox.ColorDefault, "Fin de Partie", termbox.ColorLightYellow, termbox.ColorDefault)
	for i := 1; i < len(lines); i++ {
		Print_tb(11, 3+i, termbox.ColorLightYellow, termbox.ColorDefault, lines[i])
	}
	Printf_tb(27, 13, termbox.ColorWhite, termbox.ColorDefault, "LE MOT ETAIT : ")
	Printf_tb(42, 13, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "%s", string(game.ToFind))
}

func PrintLose(game HangManData) { // Print l'écran de défaite
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	file, err := os.Open("Ressources/YouLose.txt")
	var lines []string
	if err != nil {
		log.Fatal(err)
	}
	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanLines)
	for Scanner.Scan() {
		lines = append(lines, Scanner.Text())
	}
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	Print_termbox_square(0, 0, 80, 15, termbox.ColorLightRed, termbox.ColorBlack, "Fin de Partie", termbox.ColorLightRed, termbox.ColorDefault)
	for i := 1; i < len(lines); i++ {
		Print_tb(9, 2+i, termbox.ColorLightRed, termbox.ColorDefault, lines[i])
	}
	Printf_tb(27, 13, termbox.ColorWhite, termbox.ColorDefault, "LE MOT ETAIT : ")
	Printf_tb(42, 13, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "%s", string(game.ToFind))
}

func ToHigher(s string) string { // Permet de mettre en MAJUSCULE une string
	modstr := []rune(s)
	for j := 0; j < len(s); j++ {
		if modstr[j] >= 97 && s[j] <= 122 {
			modstr[j] = modstr[j] - 32
		} else {
			continue
		}
	}
	s = string(modstr)
	return s
}

func ToLower(modstr []rune) []rune {
	for j := 0; j < len(modstr); j++ {
		if modstr[j] >= 65 && modstr[j] <= 90 {
			modstr[j] = modstr[j] + 32
		} else {
			continue
		}
	}
	return modstr
}

/*
func SortRuneTable(game HangManData) []rune {
	table := []rune(game.History)
	for j := range table {
		for k := range table {
			if table[j] < table[k] {
				table[j], table[k] = table[k], table[j]
			}
		}
	}
	return table
}
*/
