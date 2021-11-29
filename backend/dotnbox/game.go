package dotnbox

import (
	"errors"
	"sync"
	"time"
)

// Owner represents owner of a line on the grid
type Owner int

const (
	PlayerNone Owner = iota
	PlayerOne
	PlayerTwo
)

type State int

const (
	Waiting State = iota
	PlayerOneTurn
	PlayerTwoTurn
	Paused
	GameOver
)

func StateToTurn(state State) Owner {
	switch state {
	case PlayerOneTurn:
		return PlayerOne
	case PlayerTwoTurn:
		return PlayerTwo
	default:
		return PlayerNone
	}
}

// Ownership contains both the owner of a line and when they took ownership
// ts is important to determine final owner
type Ownership struct {
	Owner     Owner     `json:"owner"`
	Timestamp time.Time `json:"ts"`
}

type Game struct {
	M           int
	N           int
	state       State
	beforePause State
	grid        map[Line]Ownership
	boxes       map[Point]Owner
	scores      map[Owner]int
	mutex       sync.RWMutex
}

func NewGame(M, N int) *Game {
	return &Game{
		M:      M,
		N:      N,
		grid:   make(map[Line]Ownership),
		boxes:  make(map[Point]Owner),
		state:  Waiting,
		scores: map[Owner]int{PlayerOne: 0, PlayerTwo: 0, PlayerNone: 0},
	}
}

var ErrInvalidMove = errors.New("Invalid move")
var ErrAlreadyTaken = errors.New("Line has already been taken")
var ErrGamePaused = errors.New("Game has been paused")
var ErrGameNotOver = errors.New("Game is not over")
var ErrNotYourTurn = errors.New("It is not your turn")

func (g *Game) Evaluate() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.state == Paused {
		return
	}
	scores := map[Owner]int{PlayerOne: 0, PlayerTwo: 0, PlayerNone: 0}
	// For each box, we check if all 4 lines are placed
	// and update box ownership based on who placed the last line on the box
	for _, origin := range Boxes(g.M, g.N) {
		placedLines := 0
		var latest Ownership
		for _, line := range FindEdges(origin) {
			if ownership, exists := g.grid[line]; exists {
				placedLines++
				if ownership.Timestamp.After(latest.Timestamp) {
					latest = ownership
				}
			}
		}
		if placedLines == 4 {
			g.boxes[origin] = latest.Owner
			scores[latest.Owner]++
		}
	}

	oldscore := g.scores
	g.scores = scores

	if len(g.boxes) == g.M*g.N {
		g.state = GameOver
		return
	}

	turn := StateToTurn(g.state)
	if scores[turn] > oldscore[turn] && turn != PlayerNone {
		return
	}

	switch g.state {
	case Waiting:
		g.state = PlayerOneTurn
	case PlayerOneTurn:
		g.state = PlayerTwoTurn
	case PlayerTwoTurn:
		g.state = PlayerOneTurn
	}
}

func (g *Game) Pause() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	switch g.state {
	case PlayerOneTurn, PlayerTwoTurn:
		g.beforePause = g.state
		g.state = Paused
	default:
		return
	}
}

func (g *Game) Resume() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.state == Paused {
		g.state = g.beforePause
	}
}

func (g *Game) Rematch() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.state != GameOver {
		return ErrGameNotOver
	}
	g.grid = make(map[Line]Ownership)
	g.boxes = make(map[Point]Owner)
	g.scores = map[Owner]int{PlayerOne: 0, PlayerTwo: 0, PlayerNone: 0}
	g.state = PlayerOneTurn
	return nil
}

func (g *Game) Winner() (higestOwner Owner, highestScore int) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	for owner, score := range g.scores {
		if score > highestScore {
			higestOwner, highestScore = owner, score
		}
	}
	return
}

func (g *Game) Place(line Line, owner Owner) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.state == Paused {
		return ErrGamePaused
	}
	if StateToTurn(g.state) != owner {
		return ErrNotYourTurn
	}
	if !line.IsValid(g.M, g.N) {
		return ErrInvalidMove
	}

	if _, exists := g.grid[line]; exists {
		return ErrAlreadyTaken
	}
	g.grid[line] = Ownership{
		Owner:     owner,
		Timestamp: time.Now().UTC(),
	}
	return nil
}

func (g *Game) State() State {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.state
}

func (g *Game) Scores() map[Owner]int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.scores
}

func (g *Game) Boxes() map[string]Owner {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	boxes := make(map[string]Owner)
	for point, owner := range g.boxes {
		boxes[point.String()] = owner
	}
	return boxes
}

func (g *Game) Grid() map[string]Owner {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	grid := make(map[string]Owner)
	for line, ownership := range g.grid {
		grid[line.String()] = ownership.Owner
	}
	return grid
}
