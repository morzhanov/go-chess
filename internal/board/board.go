package board

import (
	_ "embed"
	"errors"
	"log"
	"strings"
)

const (
	rockW   = "♜"
	knightW = "♞"
	bishopW = "♝"
	queenW  = "♛"
	kingW   = "♚"
	pawnW   = "♟"
	rockB   = "♖"
	knightB = "♘"
	bishopB = "♗"
	queenB  = "♕"
	kingB   = "♔"
	pawnB   = "♙"
)

var digitIndexMap = map[string]int{
	"1": 0,
	"2": 1,
	"3": 2,
	"4": 3,
	"5": 4,
	"6": 5,
	"7": 6,
	"8": 7,
}
var letterIndexMap = map[string]int{
	"a": 0,
	"b": 1,
	"c": 2,
	"d": 3,
	"e": 4,
	"f": 5,
	"g": 6,
	"h": 7,
}

type Board interface {
	Update(request string) error
	GetPieces() [][]string
	GetMove() bool
}

type board struct {
	Pieces [][]string
	Move   bool
}

func New() Board {
	pieces := make([][]string, 8)
	for j := 2; j < 6; j++ {
		pieces[j] = []string{" ", " ", " ", " ", " ", " ", " ", " "}
	}
	pieces[0] = []string{"♜", "♞", "♝", "♛", "♚", "♝", "♞", "♜"}
	pieces[1] = []string{"♟", "♟", "♟", "♟", "♟", "♟", "♟", "♟"}
	pieces[6] = []string{"♙", "♙", "♙", "♙", "♙", "♙", "♙", "♙"}
	pieces[7] = []string{"♖", "♘", "♗", "♕", "♔", "♗", "♘", "♖"}
	return &board{Pieces: pieces, Move: true}
}

func (b *board) Update(request string) error {
	request = strings.ToLower(request)
	request = strings.Replace(request, "\n", "", 1)
	if len(request) != 4 {
		return errors.New("wrong command provided")
	}
	chars := strings.SplitAfter(request, "")
	// TODO: make right swap: https://stackoverflow.com/questions/38297882/cant-swap-elements-of-2d-array-slice-using-golang
	x1s, y1s, x2s, y2s := chars[0], chars[1], chars[2], chars[3]
	x1 := letterIndexMap[x1s]
	y1 := digitIndexMap[y1s]
	x2 := letterIndexMap[x2s]
	y2 := digitIndexMap[y2s]
	// TODO: check wheter user can perform a move
	log.Println("board before update: ", b)
	log.Printf("x1 = %d, y1 = %d, x2 = %d, y2 = %d\n", x1, y1, x2, y2)
	b.Pieces[x1][y1], b.Pieces[x2][y2] = b.Pieces[x2][y2], b.Pieces[x1][y1]
	log.Println("board after update: ", b)
	b.Move = !b.Move
	return nil
}

func (b *board) GetPieces() [][]string {
	return b.Pieces
}

func (b *board) GetMove() bool {
	return b.Move
}
