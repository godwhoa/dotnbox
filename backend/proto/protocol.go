package proto

import (
	"dotnbox/dotnbox"
	"encoding/json"
)

// Very first packet sent, on connect
type GameConfig struct {
	M      int           `json:"m"`
	N      int           `json:"n"`
	Player dotnbox.Owner `json:"player"`
}

func (gc GameConfig) Type() string {
	return "GAMECONFIG"
}

type State struct {
	Grid   map[string]dotnbox.Owner `json:"grid"`
	Boxes  map[string]dotnbox.Owner `json:"boxes"`
	State  dotnbox.State            `json:"state"`
	Scores map[dotnbox.Owner]int    `json:"scores"`
	Turn   dotnbox.Owner            `json:"turn"`
}

func (s State) Type() string {
	return "STATE"
}

func FromGame(game *dotnbox.Game) State {
	return State{
		Grid:   game.Grid(),
		Boxes:  game.Boxes(),
		State:  game.State(),
		Scores: game.Scores(),
		Turn:   dotnbox.StateToTurn(game.State()),
	}
}

type Error struct {
	Error string `json:"error"`
}

func (e Error) Type() string {
	return "ERROR"
}

type Payload struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ProtocolPayload interface {
	Type() string
}

func ToPayload(payload ProtocolPayload) Payload {
	data, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return Payload{
		Type:    payload.Type(),
		Payload: data,
	}
}
