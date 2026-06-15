package main

import (
	"fmt"
	"math/rand"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

type Cell struct {
	hidden bool
	isBomb bool
	number int
	flag   bool
}

type Cursor struct {
	x, y int
}

type model struct {
	grid      [][]Cell
	cursor    Cursor
	n         int
	gameover  bool
	flagcount int
	win       bool
	bombcount int
	mode      Difficulty
	inMenu    bool
}

func GRID(mode Difficulty) model {
	var n, bombcount int
	switch mode {
	case Easy:
		n, bombcount = 5, 5
	case Medium:
		n, bombcount = 9, 10
	case Hard:
		n, bombcount = 16, 40
	}
	return model{
		n: n,

		grid: buildGrid(n, bombcount),

		cursor: Cursor{0, 0},

		gameover: false,

		flagcount: bombcount,

		win: false,

		bombcount: bombcount,

		mode: mode,

		inMenu: true,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.inMenu {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "right":
				if m.mode != Hard {
					m.mode++
					m = GRID(m.mode)
				}
			case "left":
				if m.mode != Easy {
					m.mode--
					m = GRID(m.mode)
				}
			case "enter":
				m.inMenu = false
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "r":
				m = GRID(m.mode)
			case "up":
				if m.cursor.y > 0 {
					m.cursor.y--
				}
			case "down":
				if m.cursor.y < m.n-1 {
					m.cursor.y++
				}
			case "left":
				if m.cursor.x > 0 {
					m.cursor.x--
				}
			case "right":
				if m.cursor.x < m.n-1 {
					m.cursor.x++
				}
			case "space", "enter":
				if !m.gameover && !m.win {
					m.grid = floodFill(m.grid, m.cursor.x, m.cursor.y, m.n)
					if m.grid[m.cursor.x][m.cursor.y].isBomb && !m.grid[m.cursor.x][m.cursor.y].flag {
						m.gameover = true
						for i := range m.n {
							for j := range m.n {
								m.grid[i][j].hidden = false
							}
						}
					}
				}
			case "f", "p":
				if !m.gameover && !m.win {
					if m.grid[m.cursor.x][m.cursor.y].flag == true {
						m.grid[m.cursor.x][m.cursor.y].flag = false
						m.grid[m.cursor.x][m.cursor.y].hidden = true
						m.flagcount++
					} else {
						if m.flagcount != 0 {
							m.grid[m.cursor.x][m.cursor.y].flag = true
							m.grid[m.cursor.x][m.cursor.y].hidden = false
							m.flagcount--
						}
					}
				}
			}
		}
	}
	m.win = m.checkWin()
	return m, nil
}

func (m model) View() string {
	status := "##MINESWEEPER##\n\n"
	s := ""
	status += fmt.Sprintf("#Flag remaining: %d\n#Difficulty: %v\n\n", m.flagcount, m.mode)
	if m.inMenu {
		s += " SELECT MODE:"
		switch m.mode {
		case Easy:
			s += "[Easy] Medium  Hard "
		case Medium:
			s += " Easy [Medium] Hard "
		case Hard:
			s += " Easy  Medium [Hard]"
		}
	} else {
		if m.win {
			s += " Congratulations Nerd "
		} else {
			if m.gameover {
				status += "GAME OVER\nSKILL ISSUE\n\n"
			}
			for j := range m.n {
				for i := range m.n {
					if m.cursor.x == i && m.cursor.y == j {
						if m.grid[i][j].flag {
							s += lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render("[P]")
						} else if m.grid[i][j].hidden {
							s += lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render("[■]")
						} else {
							if m.grid[i][j].isBomb {
								s += lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("[Q]")
							} else {
								switch m.grid[i][j].number {
								case 1:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("21")).Render("[1]")
								case 2:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("28")).Render("[2]")
								case 3:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("[3]")
								case 4:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("19")).Render("[4]")
								case 5:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("88")).Render("[5]")
								case 6:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("37")).Render("[6]")
								case 7:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Render("[7]")
								case 8:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("[8]")
								default:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render("[□]")
								}
							}
						}
					} else {
						if m.grid[i][j].flag {
							s += lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render(" P ")
						} else if m.grid[i][j].hidden {
							s += lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render(" ■ ")
						} else {
							if m.grid[i][j].isBomb {
								s += lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(" Q ")
							} else {
								switch m.grid[i][j].number {
								case 1:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("21")).Render(" 1 ")
								case 2:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("28")).Render(" 2 ")
								case 3:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(" 3 ")
								case 4:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("19")).Render(" 4 ")
								case 5:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("88")).Render(" 5 ")
								case 6:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("37")).Render(" 6 ")
								case 7:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Render(" 7 ")
								case 8:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(" 8 ")
								default:
									s += lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render(" □ ")
								}
							}
						}
					}
				}
				if j != m.n-1 {
					s += "\n"
				}
			}
		}
	}
	status += lipgloss.NewStyle().BorderForeground(lipgloss.Color("67")).BorderStyle(lipgloss.RoundedBorder()).Render(s)
	status += "\nUse arror keys for navigation.\nPress p or f to flag.\nPress r to reset.\nPress q to quit.\n"

	return status
}
func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	}
	return "Unknown"
}
func (m model) checkWin() bool {
	if m.gameover {
		return false
	}
	for i := range m.n {
		for j := range m.n {
			if !m.grid[i][j].isBomb && m.grid[i][j].hidden {
				return false
			}
		}
	}
	return true
}

func floodFill(grid [][]Cell, x, y, n int) [][]Cell {
	if x >= n || y >= n || x < 0 || y < 0 {
		return grid
	}
	if !grid[x][y].hidden || grid[x][y].flag {
		return grid
	}
	if grid[x][y].hidden && grid[x][y].isBomb {
		return grid
	}
	if grid[x][y].hidden && grid[x][y].number != 0 {
		grid[x][y].hidden = false
		return grid
	}
	grid[x][y].hidden = false
	grid = floodFill(grid, x+1, y, n)
	grid = floodFill(grid, x+1, y+1, n)
	grid = floodFill(grid, x+1, y-1, n)
	grid = floodFill(grid, x, y+1, n)
	grid = floodFill(grid, x, y-1, n)
	grid = floodFill(grid, x-1, y, n)
	grid = floodFill(grid, x-1, y+1, n)
	grid = floodFill(grid, x-1, y-1, n)
	return grid
}

func placeBomb(n int, bomb_count int, grid [][]Cell) [][]Cell {
	for i := 0; i < bomb_count; {
		x := rand.Intn(n)
		y := rand.Intn(n)
		if !grid[x][y].isBomb {
			grid[x][y].isBomb = true
			i++
		}
	}
	return grid
}

func placeNum(n int, grid [][]Cell) [][]Cell {
	var number int = 0
	for i := range n {
		for j := range n {
			if grid[i][j].isBomb != true {
				if i > 0 && grid[i-1][j].isBomb == true {
					number++
				}
				if i > 0 && j > 0 && grid[i-1][j-1].isBomb == true {
					number++
				}
				if j > 0 && grid[i][j-1].isBomb == true {
					number++
				}
				if i < n-1 && j > 0 && grid[i+1][j-1].isBomb == true {
					number++
				}
				if i < n-1 && grid[i+1][j].isBomb == true {
					number++
				}
				if i < n-1 && j < n-1 && grid[i+1][j+1].isBomb == true {
					number++
				}
				if j < n-1 && grid[i][j+1].isBomb == true {
					number++
				}
				if i > 0 && j < n-1 && grid[i-1][j+1].isBomb == true {
					number++
				}
				grid[i][j].number = number
				number = 0
			}
		}
	}
	return grid
}

func buildGrid(n int, bomb_count int) [][]Cell {
	grid := make([][]Cell, n)
	for i := range grid {
		grid[i] = make([]Cell, n)
	}
	for i := range n {
		for j := range n {
			grid[i][j].hidden = true
		}
	}
	grid = placeBomb(n, bomb_count, grid)
	grid = placeNum(n, grid)
	return grid
}

func main() {
	mode := Easy
	p := tea.NewProgram(GRID(mode), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
