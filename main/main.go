package main

import (
	"hangman"
	"os"

	"github.com/nsf/termbox-go"
)

func main() {
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "--startWith" && os.Args[i+1] == "save.txt" {
			err := termbox.Init()
			if err != nil {
				panic(err)
			}
			defer termbox.Close()
			Print_termbox_game(hangman.LoadGame())
		} else if i+1 == len(os.Args) {
			Hangman(hangman.RandomPickLine())
		}
	}
}

func Hangman(game hangman.HangManData) {
	T := hangman.RandomPickLetter(game)
	game.Attempts = 10 // définit le nombre d'essai total avant le début du jeu
	for i := 0; i < len(game.ToFind); i++ {
		game.Word = append(game.Word, '_')
	}
	for j := 0; j < len(T); j++ {
		game.Word[T[j]] = game.ToFind[T[j]] // permet de
	}
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	Print_termbox_description(game)
}

func Print_termbox_game(game hangman.HangManData) {
	var WordClic []rune
	var TestingString []rune
	var histobool bool = false
	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	hangman.Print_tb_game_boxes(game)
	termbox.Flush()
	for {
		count := 0
		for g := 0; g < len(game.Word); g++ { // condition de victoire
			if game.Word[g] != '_' {
				count++
			}
			if count == len(game.Word) || string(TestingString) == string(game.ToFind) {
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				termbox.Flush()
				hangman.PrintWin(game)
				termbox.Flush()
				RestartWin(game)
				termbox.Close()
			}
		}
		count = 0
		if game.Attempts <= 0 { // condition de défaite
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			termbox.Flush()
			hangman.PrintLose(game)
			termbox.Flush()
			RestartLose(game)
			termbox.Close()
		} else {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				if ev.Key == termbox.KeyEsc { // Permet de sauvegarder la partie si la touche ESC est pressée
					hangman.SaveGame(game)
					termbox.Close()
					hangman.Clear()
					os.Exit(0)
				}
				if ev.Key == termbox.KeyArrowLeft {
					WordClic = nil
					Print_termbox_Help(game)
				}
				if ev.Key == termbox.KeyArrowRight {
					WordClic = nil
					Print_termbox_description(game)
				} else {
					if ev.Key == termbox.KeyEnter {
						if len(WordClic) == 0 {
							WordClic = nil
							break
						}
						TestingString = WordClic
						WordClic = nil
						for j := 0; j < len(game.ToFind); j++ {
							ct := 0
							if rune(TestingString[0]) == game.ToFind[j] && len(TestingString) == 1 {
								hangman.ModifyGameWord(game, string(TestingString)) // func qui transforme la string Word en mettant en place la lettre
								break
							} else if len(TestingString) > 1 && string(TestingString) != string(game.ToFind) && j == len(game.ToFind)-1 && ct < 1 {
								game.History += string(TestingString[0])
								game.Attempts -= 2
								ct++
								break
							} else if j == len(game.ToFind)-1 { // si la lettre testée est fausse alors on l'ajoute à la string Historique pour avoir un suivi et ne pas réessayer les même lettres
								for i := 0; i < len(game.History); i++ { // Permet de ne pas avoir une lettre qui apparaît plusieurs fois si elle est testée plusieurs fois
									if string(TestingString[0]) == string(game.History[i]) {
										histobool = true
										game.Attempts--
									}
								}
								if !histobool {
									game.History += string(TestingString[0])
									game.Attempts--
								}
								histobool = false
							}
						}
						termbox.Flush()
					}
					if ev.Ch == 0 && len(WordClic) > 0 {
						WordClic = WordClic[0 : len(WordClic)-1]
					} else if ev.Key == termbox.KeyBackspace {
						WordClic = WordClic[:len(WordClic)-1]
					} else if (ev.Ch >= 97 && ev.Ch <= 122) || (ev.Ch >= 65 && ev.Ch <= 90) {
						WordClic = append(WordClic, ev.Ch)
						WordClic = hangman.ToLower(WordClic)
					}
					termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
					hangman.Print_tb_game_boxes(game)
					hangman.Printf_tb(4, 20, termbox.ColorWhite, termbox.ColorDefault, hangman.ToHigher(string(WordClic)))
					termbox.Flush()
				}
			case termbox.EventResize:
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				Print_termbox_game(game)
				termbox.Flush()
			case termbox.EventError:
				panic(ev.Err)
			}
		}
	}
}

func Print_termbox_description(game hangman.HangManData) {
	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	hangman.Print_tb_desc_boxes(game)
	termbox.Flush()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc { // Permet de sauvegarder la partie si la touche ESC est pressée
				hangman.SaveGame(game)
				termbox.Close()
				hangman.Clear()
				os.Exit(0)
			}
			if ev.Key == termbox.KeyArrowRight {
				Print_termbox_Help(game)
			}
			if ev.Key == termbox.KeyArrowLeft {
				Print_termbox_game(game)
			}
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			Print_termbox_description(game)
			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func Print_termbox_Help(game hangman.HangManData) {
	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	hangman.Print_tb_help_boxes(game)
	termbox.Flush()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc { // Permet de sauvegarder la partie si la touche ESC est pressée
				hangman.SaveGame(game)
				termbox.Close()
				hangman.Clear()
				os.Exit(0)
			}
			if ev.Key == termbox.KeyArrowRight {
				Print_termbox_game(game)
			}
			if ev.Key == termbox.KeyArrowLeft {
				Print_termbox_description(game)
			}
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			Print_termbox_Help(game)
			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func RestartLose(game hangman.HangManData) { // permet la boucle du programme
	var WordClic []rune
	hangman.PrintLose(game)
	hangman.Print_termbox_square(0, 16, 80, 6, termbox.ColorWhite, termbox.ColorBlack, "", termbox.ColorRed, termbox.ColorDefault)
	hangman.Printf_tb(11, 18, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "RECOMMENCER ? (Ecrivez oui si vous souhaitez recommencer)")
	termbox.Flush()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc { // Permet de sauvegarder la partie si la touche ESC est pressée
				termbox.Close()
				hangman.Clear()
				os.Exit(0)
			} else {
				if ev.Key == termbox.KeyEnter {
					if string(WordClic) == "oui" {
						termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
						termbox.Close()
						Hangman(hangman.RandomPickLine())
					} else {
						hangman.Clear()
						termbox.Close()
						os.Exit(0)
					}
				}
				if ev.Ch == 0 && len(WordClic) > 0 {
					WordClic = WordClic[0 : len(WordClic)-1]
				} else if ev.Key == termbox.KeyBackspace {
					WordClic[len(WordClic)-1] = 0
				} else if (ev.Ch >= 97 && ev.Ch <= 122) || (ev.Ch >= 65 && ev.Ch <= 90) {
					WordClic = append(WordClic, ev.Ch)
					WordClic = hangman.ToLower(WordClic)
				}
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				hangman.PrintLose(game)
				hangman.Print_termbox_square(0, 16, 80, 6, termbox.ColorWhite, termbox.ColorBlack, "", termbox.ColorRed, termbox.ColorDefault)
				hangman.Printf_tb(11, 18, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "RECOMMENCER ? (Ecrivez oui si vous souhaitez recommencer)")
				hangman.Printf_tb(4, 20, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, hangman.ToHigher(string(WordClic)))
				termbox.Flush()
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func RestartWin(game hangman.HangManData) { // permet la boucle du programme
	var WordClic []rune
	hangman.PrintWin(game)
	hangman.Print_termbox_square(0, 16, 80, 6, termbox.ColorWhite, termbox.ColorBlack, "", termbox.ColorRed, termbox.ColorDefault)
	hangman.Printf_tb(11, 18, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "RECOMMENCER ? (Ecrivez oui si vous souhaitez recommencer)")
	termbox.Flush()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc { // Permet de sauvegarder la partie si la touche ESC est pressée
				termbox.Close()
				hangman.Clear()
				os.Exit(0)
			} else {
				if ev.Key == termbox.KeyEnter {
					if string(WordClic) == "oui" {
						termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
						termbox.Close()
						Hangman(hangman.RandomPickLine())
					} else {
						hangman.Clear()
						termbox.Close()
						os.Exit(0)
					}
				}
				if ev.Ch == 0 && len(WordClic) > 0 {
					WordClic = WordClic[0 : len(WordClic)-1]
				} else if ev.Key == termbox.KeyBackspace {
					WordClic[len(WordClic)-1] = 0
				} else if (ev.Ch >= 97 && ev.Ch <= 122) || (ev.Ch >= 65 && ev.Ch <= 90) {
					WordClic = append(WordClic, ev.Ch)
					WordClic = hangman.ToLower(WordClic)
				}
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				hangman.PrintWin(game)
				hangman.Print_termbox_square(0, 16, 80, 6, termbox.ColorWhite, termbox.ColorBlack, "", termbox.ColorRed, termbox.ColorDefault)
				hangman.Printf_tb(11, 18, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "RECOMMENCER ? (Ecrivez oui si vous souhaitez recommencer)")
				hangman.Printf_tb(4, 20, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, hangman.ToHigher(string(WordClic)))
				termbox.Flush()
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
