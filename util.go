package main

import (
	"bufio"
	c "color"
	. "fmt"
	r "math/rand"
	"os"
	"os/exec"
)

//Creation & settings

func CreatePlayer() Player {
	var temp Player
	temp.Alive = true
	temp.Dir = 4
	temp.Pos.X = Dim - 2
	temp.Pos.Y = 1
	return temp
}

func CreateGhost() Ghost {
	var temp Ghost
	temp.Alive = true
	temp.Dir = 0
	temp.Pos.X = Dim / 2
	temp.Pos.Y = Dim / 2
	return temp
}

func CreateLab() Lab {
	var temp Lab
	file, err := os.Open("tab")
	defer file.Close()
	if err != nil {
		return temp
	}
	scanner := bufio.NewScanner(file)
	var j int
	for scanner.Scan() {
		var a [Dim]int
		line := scanner.Text()
		for i, c := range line {
			switch c {
			case '#':
				a[i] = -1
			case 'C':
				a[i] = 1
			case '*':
				a[i] = 0
			default:
				a[i] = 2
			}
		}
		temp[j] = a
		j++
	}
	temp[4][1] = 3
	temp[Dim-5][Dim-2] = 3
	return temp
}

//Output control

func Refresh() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Start()
	cmd.Wait()
}

func (l *Lab) Print() {
	if Special {
		c.Red("SPECIAL ACTIVATED")
	} else {
		c.Red("SPECIAL DEACTIVATED")
	}
	for i := 0; i < len(l); i++ {
		for j := 0; j < len(l[i]); j++ {
			if l[i][j] == -1 {
				Print("\033[1m##")
			} else if l[i][j] == 0 {
				Print("  ")
			} else if l[i][j] == 1 {
				Print(c.YellowString("%s", "C "))
			} else if l[i][j] == 2 {
				Print("° ")
			} else if l[i][j] == 3 {
				if r.Intn(2) == 1 {
					Print("0 ")
				} else {
					Print("* ")
				}
			} else {
				switch l[i][j] {
				case 4:
					Print(c.RedString("%s", "& "))
				case 5:
					Print(c.GreenString("%s", "& "))
				case 6:
					Print(c.BlueString("%s", "& "))
				case 7:
					Print(c.CyanString("%s", "& "))
				case 8:
					Print(c.MagentaString("%s\033[0m", "& "))
				}
			}
		}
		Println()
	}
	c.Yellow("Points: %d", Points)
	for i := 0; i < 11; i++ {
		Println()
	}
}
