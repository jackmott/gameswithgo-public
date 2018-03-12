// TODO On Stream  ** Starting Soon **
// - Make a simple Event display UI, remove old events, put more stuff in the events
// - Fix quick keypress bug that occurs with many monsters
// - Make monsters route around other monsters

package game

import (
	"bufio"
	"fmt"
	"math"
	"os"
	_ "strconv"
)

type Game struct {
	LevelChans []chan *Level
	InputChan  chan *Input
	Level      *Level
}

func NewGame(numWindows int, levelPath string) *Game {
	levelChans := make([]chan *Level, numWindows)
	for i := range levelChans {
		levelChans[i] = make(chan *Level)
	}
	inputChan := make(chan *Input)

	return &Game{levelChans, inputChan, loadLevelFromFile(levelPath)}

}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	QuitGame
	CloseWindow
	Search //temporary
)

type Input struct {
	Typ          InputType
	LevelChannel chan *Level
}

type Tile rune

const (
	StoneWall  Tile = '#'
	DirtFloor  Tile = '.'
	ClosedDoor Tile = '|'
	OpenDoor   Tile = '/'
	Blank      Tile = 0
	Pending    Tile = -1
)

type Pos struct {
	X, Y int
}

type Entity struct {
	Pos
	Name string
	Rune rune
}

type Character struct {
	Entity
	Hitpoints    int
	Strength     int
	Speed        float64
	ActionPoints float64
}

type Player struct {
	Character
}

type Level struct {
	Map      [][]Tile
	Player   *Player
	Monsters map[Pos]*Monster
	Events   []string
	EventPos int
	Debug    map[Pos]bool
}

type Attackable interface {
	GetActionPoints() float64
	SetActionPoints(float64)
	GetHitpoints() int
	SetHitpoints(int)
	GetAttackPower() int
}

func (c *Character) GetActionPoints() float64 {
	return c.ActionPoints
}

func (c *Character) SetActionPoints(ap float64) {
	c.ActionPoints = ap
}

func (c *Character) GetHitpoints() int {
	return c.Hitpoints
}

func (c *Character) SetHitpoints(h int) {
	c.Hitpoints = h
}

func (c *Character) GetAttackPower() int {
	return c.Strength
}

func Attack(a1, a2 Attackable) {
	a1.SetActionPoints(a1.GetActionPoints() - 1)
	a2.SetHitpoints(a2.GetHitpoints() - a1.GetAttackPower())
}

func (level *Level) AddEvent(event string) {
	level.Events[level.EventPos] = event

	level.EventPos++
	if level.EventPos == len(level.Events) {
		level.EventPos = 0
	}
}

func loadLevelFromFile(filename string) *Level {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	levelLines := make([]string, 0)
	longestRow := 0
	index := 0
	for scanner.Scan() {
		levelLines = append(levelLines, scanner.Text())
		if len(levelLines[index]) > longestRow {
			longestRow = len(levelLines[index])
		}
		index++
	}
	level := &Level{}
	level.Events = make([]string, 10)
	level.Player = &Player{}
	// TODO where should we initialize the player?
	level.Player.Strength = 5
	level.Player.Hitpoints = 20
	level.Player.Name = "GoMan"
	level.Player.Rune = '@'
	level.Player.Speed = 1.0
	level.Player.ActionPoints = 0

	level.Map = make([][]Tile, len(levelLines))
	level.Monsters = make(map[Pos]*Monster)

	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}

	for y := 0; y < len(level.Map); y++ {
		line := levelLines[y]
		for x, c := range line {
			var t Tile
			switch c {
			case ' ', '\t', '\n', '\r':
				t = Blank
			case '#':
				t = StoneWall
			case '|':
				t = ClosedDoor
			case '/':
				t = OpenDoor
			case '.':
				t = DirtFloor
			case '@':
				level.Player.X = x
				level.Player.Y = y
				t = Pending
			case 'R':
				level.Monsters[Pos{x, y}] = NewRat(Pos{x, y})
				t = Pending
			case 'S':
				level.Monsters[Pos{x, y}] = NewSpider(Pos{x, y})
				t = Pending
			default:
				panic("Invalid character in map")
			}
			level.Map[y][x] = t

		}
	}

	for y, row := range level.Map {
		for x, tile := range row {
			if tile == Pending {
				level.Map[y][x] = level.bfsFloor(Pos{x, y})
			}
		}
	}

	return level
}

func inRange(level *Level, pos Pos) bool {
	return pos.X < len(level.Map[0]) && pos.Y < len(level.Map) && pos.X >= 0 && pos.Y >= 0
}

func canWalk(level *Level, pos Pos) bool {
	if inRange(level, pos) {
		t := level.Map[pos.Y][pos.X]
		switch t {
		case StoneWall, ClosedDoor, Blank:
			return false
		default:
			return true
		}
	}
	return false
}

func checkDoor(level *Level, pos Pos) {
	t := level.Map[pos.Y][pos.X]
	if t == ClosedDoor {
		level.Map[pos.Y][pos.X] = OpenDoor
	}
}

func (player *Player) Move(to Pos, level *Level) {
	monster, exists := level.Monsters[to]
	if !exists {
		player.Pos = to
		return
	}
	Attack(level.Player, monster)
	fmt.Println("Player Attacked Monster")
	level.AddEvent("Player Attacked Monster")
	fmt.Println(level.Player.Hitpoints, monster.Hitpoints)
	if monster.Hitpoints <= 0 {
		level.AddEvent("Player Kills " + monster.Name)
		delete(level.Monsters, monster.Pos)
	}
	if level.Player.Hitpoints <= 0 {
		fmt.Println("YOU DIED")
		panic("You Died") // exit properly?

	}

}

func (game *Game) handleInput(input *Input) {
	level := game.Level
	p := level.Player
	switch input.Typ {
	case Up:
		newPos := Pos{p.X, p.Y - 1}
		if canWalk(level, newPos) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, Pos{p.X, p.Y - 1})
		}
	case Down:
		newPos := Pos{p.X, p.Y + 1}
		if canWalk(level, Pos{p.X, p.Y + 1}) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, Pos{p.X, p.Y + 1})
		}
	case Left:
		newPos := Pos{p.X - 1, p.Y}
		if canWalk(level, Pos{p.X - 1, p.Y}) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, Pos{p.X - 1, p.Y})
		}
	case Right:
		newPos := Pos{p.X + 1, p.Y}
		if canWalk(level, Pos{p.X + 1, p.Y}) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, Pos{p.X + 1, p.Y})
		}
	case Search:
		//bfs(ui, level, level.Player.Pos)
		level.astar(level.Player.Pos, Pos{3, 2})
	case CloseWindow:
		close(input.LevelChannel)
		chanIndex := 0
		for i, c := range game.LevelChans {
			if c == input.LevelChannel {
				chanIndex = i
				break
			}
		}
		game.LevelChans = append(game.LevelChans[:chanIndex], game.LevelChans[chanIndex+1:]...)
	}
}

func getNeighbors(level *Level, pos Pos) []Pos {
	neighbors := make([]Pos, 0, 4)
	left := Pos{pos.X - 1, pos.Y}
	right := Pos{pos.X + 1, pos.Y}
	up := Pos{pos.X, pos.Y - 1}
	down := Pos{pos.X, pos.Y + 1}

	if canWalk(level, right) {
		neighbors = append(neighbors, right)
	}
	if canWalk(level, left) {
		neighbors = append(neighbors, left)
	}
	if canWalk(level, up) {
		neighbors = append(neighbors, up)
	}
	if canWalk(level, down) {
		neighbors = append(neighbors, down)
	}

	return neighbors
}

func (level *Level) bfsFloor(start Pos) Tile {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool)
	visited[start] = true
	level.Debug = visited

	for len(frontier) > 0 {
		current := frontier[0]

		currentTile := level.Map[current.Y][current.X]
		switch currentTile {
		case DirtFloor:
			return DirtFloor
		default:
		}

		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true

			}
		}

	}
	return DirtFloor
}

func (level *Level) astar(start Pos, goal Pos) []Pos {
	frontier := make(pqueue, 0, 8)
	frontier = frontier.push(start, 1)
	cameFrom := make(map[Pos]Pos)
	cameFrom[start] = start
	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0

	var current Pos
	for len(frontier) > 0 {

		frontier, current = frontier.pop()

		if current == goal {
			path := make([]Pos, 0)
			p := current
			for p != start {
				path = append(path, p)
				p = cameFrom[p]
			}
			path = append(path, p)

			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			return path
		}

		for _, next := range getNeighbors(level, current) {
			newCost := costSoFar[current] + 1 // always 1 for now
			_, exists := costSoFar[next]
			if !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))
				priority := newCost + xDist + yDist
				frontier = frontier.push(next, priority)
				cameFrom[next] = current
			}

		}

	}

	return nil

}

func (game *Game) Run() {
	fmt.Println("Starting...")

	count := 0
	for _, lchan := range game.LevelChans {
		lchan <- game.Level
	}

	for input := range game.InputChan {
		if input.Typ == QuitGame {
			return
		}

		game.handleInput(input)

		//game.Level.AddEvent("move:" + strconv.Itoa(count))
		count++

		for _, monster := range game.Level.Monsters {
			monster.Update(game.Level)
		}

		if len(game.LevelChans) == 0 {
			return
		}

		for _, lchan := range game.LevelChans {
			lchan <- game.Level
		}
	}

}
