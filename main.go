package main

import (
	"bufio"
	c "color"
	. "fmt"
	r "math/rand"
	"os"
	"os/exec"
	"runtime"
	t "time"

	term "github.com/nsf/termbox-go"
)

//Title animation
func Animated(i int, s string) string {
	temp := ""
	for j := 0; j < 84; j++ {
		temp += string(s[j])
	}
	for j := 84; j < len(s); j++ {
		if j < i {
			temp += " "
			continue
		}
		if j != i {
			temp += string(s[j])
		} else {
			temp += c.YellowString("%s", "C")
		}
	}
	return temp
}

func MainMenu() error {
	Refresh()
	titlefile, _ := os.Open("title")
	scanner := bufio.NewScanner(titlefile)
	title := ""
	for scanner.Scan() {
		title += scanner.Text() + "\n"
	}
	var i int
	for i = 86; i < 125; i++ {
		Refresh()
		Println(Animated(i, title))
		t.Sleep(t.Second / 20)
	}
	Refresh()
	Print("\033[1m" + title + "\033[0m")
	Println("\t  by TRTNA\n\n")
	Println("Remember to maximize your terminal windows and set a proper zoom")
	Println("You can read the game guide in the README file")
	Print("\n\n\033[1m...Press enter to start the game...\033[0m")
	r := bufio.NewReader(os.Stdin)
	var line []byte
	line = append(line, 1)
	for len(line) != 0 {
		line, _, _ = r.ReadLine()
	}
	Refresh()
	return nil
}

func main() {
	//Creation of player and ghosts
	//Various settings
	runtime.GOMAXPROCS(4)
	r.Seed(int64(t.Now().Nanosecond()))
	l := CreateLab()
	player := CreatePlayer()
	ghosts := make([]Ghost, 5)
	for i := 0; i < 5; i++ {
		ghosts[i] = CreateGhost()
	}
	//Starting menu
	MainMenu()
	//Goroutines and music
	go Handling(&player, &l)
	go CheckDeath(&player, &ghosts)
	cmd := exec.Command("play", "theme.mp3", "repeat")
	cmd.Stdout = os.Stdout
	cmd.Start()
	var alldied bool
	//Game
	for {
		alldied = true
		l.Print()
		_Ghost(&ghosts, &l)
		if player.Alive == false {
			term.Close()
			break
		}
		for _, g := range ghosts {
			if !g.Alive {
				if l[g.Pos.X][g.Pos.Y] != 1 {
					l[g.Pos.X][g.Pos.Y] = 0
				} else {
					l[g.Pos.X][g.Pos.Y] = 1
				}
			} else {
				alldied = false
			}
		}
		if alldied {
			term.Close()
			break
		}
	}
	Refresh()
	if player.Alive {
		Println("ALL GHOSTS KILLED")
		Println("Thank you for playing HANEMAN by TRTNA")
	} else {
		Println("GAME OVER")
		Println("Thank you for playing HANEMAN by TRTNA")
	}
	//Turn off the music
	exit := exec.Command("killall", "play")
	exit.Stdout = os.Stdout
	exit.Run()
}
