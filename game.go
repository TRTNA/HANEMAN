package main

import (
	r "math/rand"
	t "time"

	term "github.com/nsf/termbox-go"
)

//Player's functions
func Handling(player *Player, l *Lab) {
	err := term.Init()
	if err != nil {
		panic(err)
	}
	for {
		if ((player.Pos.X == 4 && player.Pos.Y == 1) || (player.Pos.X == Dim-5 && player.Pos.Y == Dim-2)) && Special == false {
			Special = true
			go Wait(l)
		}
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				term.Close()
			case term.KeyArrowUp:
				if CheckUp(&player.Pos, l) {
					Move(2, &player.Pos)
				}
				l.Update(player)
			case term.KeyArrowDown:
				if CheckDown(&player.Pos, l) {
					Move(0, &player.Pos)
				}
				l.Update(player)
			case term.KeyArrowLeft:
				if CheckLeft(&player.Pos, l) {
					Move(3, &player.Pos)
				}
				l.Update(player)
			case term.KeyArrowRight:
				if CheckRight(&player.Pos, l) {
					Move(1, &player.Pos)
				}
				l.Update(player)
			default:
				break
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
	term.Close()
}

//Update player's position on the table
func (l *Lab) Update(p *Player) {
	if l[p.Pos.X+1][p.Pos.Y] == 1 {
		l[p.Pos.X+1][p.Pos.Y] = 0
	}
	if l[p.Pos.X][p.Pos.Y+1] == 1 {
		l[p.Pos.X][p.Pos.Y+1] = 0
	}
	if l[p.Pos.X-1][p.Pos.Y] == 1 {
		l[p.Pos.X-1][p.Pos.Y] = 0
	}
	if l[p.Pos.X][p.Pos.Y-1] == 1 {
		l[p.Pos.X][p.Pos.Y-1] = 0
	}
	l[p.Pos.X][p.Pos.Y] = 1
}

//Generic movement and handling functions

func Move(dir byte, pos *Position) {
	switch dir {
	case 0:
		(*pos).X++
	case 1:
		(*pos).Y++
	case 2:
		(*pos).X--
	case 3:
		(*pos).Y--
	}
}

func CheckDown(pos *Position, t *Lab) bool {
	if pos.X != Dim-2 && (*t)[pos.X+1][pos.Y] != -1 {
		if (*t)[pos.X+1][pos.Y] == 2 {
			Points++
		}
		return true
	}
	return false
}

func CheckUp(pos *Position, t *Lab) bool {
	if pos.X != 1 && (*t)[pos.X-1][pos.Y] != -1 {
		if (*t)[pos.X-1][pos.Y] == 2 {
			Points++
		}
		return true
	}
	return false
}

func CheckLeft(pos *Position, t *Lab) bool {
	if pos.Y != 1 && (*t)[pos.X][pos.Y-1] != -1 {
		if (*t)[pos.X][pos.Y-1] == 2 {
			Points++
		}
		return true
	}
	return false
}

func CheckRight(pos *Position, t *Lab) bool {
	if pos.Y != Dim-2 && (*t)[pos.X][pos.Y+1] != -1 {
		if (*t)[pos.X][pos.Y+1] == 2 {
			Points++
		}
		return true
	}
	return false
}

func CheckGhostDir(pos *Position, t *Lab) []byte {
	a := make([]byte, 0)
	if pos.X != Dim-2 && (*t)[pos.X+1][pos.Y] != -1 {
		a = append(a, 0)
	}
	if pos.X != 1 && (*t)[pos.X-1][pos.Y] != -1 {
		a = append(a, 2)
	}
	if pos.Y != 1 && (*t)[pos.X][pos.Y-1] != -1 {
		a = append(a, 3)
	}
	if pos.Y != Dim-2 && (*t)[pos.X][pos.Y+1] != -1 {
		a = append(a, 1)
	}
	return a
}

//
//
//Ghosts' functions

func HasToChangeDir(l *Lab, ghost *Ghost) bool {
	switch ghost.Dir {
	case 0:
		if (*l)[ghost.Pos.X+1][ghost.Pos.Y] == -1 {
			return true
		}
	case 1:
		if (*l)[ghost.Pos.X][ghost.Pos.Y+1] == -1 {
			return true
		}
	case 2:
		if (*l)[ghost.Pos.X-1][ghost.Pos.Y] == -1 {
			return true
		}
	case 3:
		if (*l)[ghost.Pos.X][ghost.Pos.Y-1] == -1 {
			return true
		}
	}
	return false
}

func ChangeDir(l *Lab, ghost *Ghost) byte {
	var possdir []byte
	switch ghost.Dir {
	case 0, 2:
		if (*l)[ghost.Pos.X][ghost.Pos.Y+1] != -1 {
			possdir = append(possdir, 1)
		}
		if (*l)[ghost.Pos.X][ghost.Pos.Y-1] != -1 {
			possdir = append(possdir, 3)
		}
	case 1, 3:
		if (*l)[ghost.Pos.X+1][ghost.Pos.Y] != -1 {
			possdir = append(possdir, 0)
		}
		if (*l)[ghost.Pos.X-1][ghost.Pos.Y] != -1 {
			possdir = append(possdir, 2)
		}
	}
	if len(possdir) == 0 {
		return ghost.Dir
	}
	return possdir[r.Intn(len(possdir))]
}

//Handling function for the ghosts
func _Ghost(g *[]Ghost, l *Lab) {
	t.Sleep(t.Second / 10)
	for i := 0; i < len((*g)); i++ {
		if (*g)[i].Alive {
			if HasToChangeDir(l, &(*g)[i]) {
				listdir := CheckGhostDir(&(*g)[i].Pos, l)
				dir := listdir[r.Intn(len(listdir))]
				Move(dir, &((*g)[i]).Pos)
				(*g)[i].Dir = dir
			} else {
				dir := ChangeDir(l, &(*g)[i])
				if int64(t.Now().Second())%2 == 0 {
					Move((*g)[i].Dir, &((*g)[i].Pos))
				} else {
					Move(dir, &((*g)[i].Pos))
					(*g)[i].Dir = dir
				}
			}
			if (*l)[(*g)[i].Pos.X+1][(*g)[i].Pos.Y] == (i + 4) {
				(*l)[(*g)[i].Pos.X+1][(*g)[i].Pos.Y] = 2
			}
			if (*l)[(*g)[i].Pos.X][(*g)[i].Pos.Y+1] == (i + 4) {
				(*l)[(*g)[i].Pos.X][(*g)[i].Pos.Y+1] = 2
			}
			if (*l)[(*g)[i].Pos.X-1][(*g)[i].Pos.Y] == (i + 4) {
				(*l)[(*g)[i].Pos.X-1][(*g)[i].Pos.Y] = 2
			}
			if (*l)[(*g)[i].Pos.X][(*g)[i].Pos.Y-1] == (i + 4) {
				(*l)[(*g)[i].Pos.X][(*g)[i].Pos.Y-1] = 2
			}
			(*l)[(*g)[i].Pos.X][(*g)[i].Pos.Y] = (i + 4)
		}
	}
}

//Deactivate the special after 10 seconds
func Wait(l *Lab) {
	l[4][1] = 0
	l[Dim-5][Dim-2] = 0
	t.Sleep(10 * t.Second)
	Special = false
	l[4][1] = 3
	l[Dim-5][Dim-2] = 3
	return
}

//Constantly checks if the player or a ghost has died
func CheckDeath(player *Player, ghosts *[]Ghost) {
	for {
		for i := 0; i < len(*ghosts); i++ {
			if ((*ghosts)[i]).Alive {
				if !Special {
					if (*player).Pos.X == (*ghosts)[i].Pos.X && (*ghosts)[i].Pos.Y == (*player).Pos.Y {
						player.Alive = false
						return
					}
				} else {
					if (*player).Pos.X == (*ghosts)[i].Pos.X && (*ghosts)[i].Pos.Y == (*player).Pos.Y {
						(*ghosts)[i].Alive = false
					}
				}
			}
		}
	}
}
