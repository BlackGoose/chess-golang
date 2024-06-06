package main

import (
	"fmt"
	"strings"
)

// Перенести передвижения в мапу
// Только конь может перепрыгивать через союзников
// Король может делать рокировку
// Победа, если умер король
// Заменить ошибки на переменные
// Вести счет матча
// Превращение фигуры (пешки) в любую другую, если она дошла до конца поля
// Писать, кто ходит сейчас, проверять, что двигается фигура соответствующего цвета

var defaultBoard [8][8]figure = [8][8]figure{
	{figure{"rook", "♖", false}, figure{"knight", "♘", false}, figure{"bishop", "♗", false}, figure{"king", "♔", false}, figure{"queen", "♕", false}, figure{"bishop", "♗", false}, figure{"knight", "♘", false}, figure{"rook", "♖", false}},
	{figure{"pawn", "♙", false}, figure{"pawn", "♙", false}, figure{"pawn", "♙", false}, figure{"pawn", "♙", false}, figure{"pawn", "♙", false}, figure{"pawn", "♙", false}, figure{"pawn", "♙", false}, figure{"pawn", "♙", false}},
	{figure{}, figure{}, figure{}, figure{}, figure{}, figure{}, figure{}, figure{}},
	{figure{}, figure{}, figure{}, figure{}, figure{}, figure{}, figure{}, figure{}},
	{figure{}, figure{}, figure{}, figure{}, figure{}, figure{}, figure{}, figure{}},
	{figure{}, figure{}, figure{}, figure{}, figure{}, figure{}, figure{}, figure{}},
	{figure{"pawn", "♟", true}, figure{"pawn", "♟", true}, figure{"pawn", "♟", true}, figure{"pawn", "♟", true}, figure{"pawn", "♟", true}, figure{"pawn", "♟", true}, figure{"pawn", "♟", true}, figure{"pawn", "♟", true}},
	{figure{"rook", "♜", true}, figure{"knight", "♞", true}, figure{"bishop", "♝", true}, figure{"king", "♚", true}, figure{"queen", "♛", true}, figure{"bishop", "♝", true}, figure{"knight", "♞", true}, figure{"rook", "♜", true}},
}

type board struct {
	field [8][8]figure
	score [2]int
}

type figure struct {
	name, icon string
	color      bool
}

type coordinates struct {
	x rune
	y int
}

func InitBoard() *board {
	return &board{
		field: defaultBoard,
		score: [2]int{0, 0},
	}
}

func (b board) GetStrField() string {
	result := strings.Builder{}
	for i := len(b.field) - 1; i >= 0; i-- {
		result.WriteString(fmt.Sprint(i+1, "| "))
		for j := range b.field[i] {
			if b.field[i][j].name != "" {
				result.WriteString(b.field[i][j].icon + " ")
			} else {
				if (i+j)%2 == 0 {
					result.WriteString("◻ ")
				} else {
					result.WriteString("◼️ ")
				}
			}
		}
		result.WriteByte('\n')
	}
	result.WriteString("-+----------------\n | a b c d e f g h")
	return result.String()
}

func (b *board) move(x1, y1, x2, y2 int) {
	b.field[y2][x2] = b.field[y1][x1]
	b.field[y1][x1] = figure{}
}

func (b *board) Move(cor1, cor2 coordinates) error {
	x1, y1, x2, y2 := int(cor1.x-97), cor1.y-1, int(cor2.x-97), cor2.y-1
	if x1 >= len(b.field) || x2 >= len(b.field) || y1 >= len(b.field) || y2 >= len(b.field) || x1 == x2 && y1 == y2 {
		return fmt.Errorf("coordinates should be in [a-h][1-8] and must be different, choose another squares")
	}
	if b.field[y1][x1].name == "" {
		return fmt.Errorf("no figure at the first coordinates, choose another square")
	}
	f := b.field[y1][x1]

	if b.field[y2][x2].name != "" && b.field[y2][x2].color == f.color {
		return fmt.Errorf("allias figure at the second coordinates, choose another square")
	}

	switch f.name {
	case "pawn":
		if b.field[y2][x2].name != "" {
			if (x1-x2)*(x1-x2) == 1 && (y2-y1 == -1 && f.color || y2-y1 == 1 && !f.color) {
				b.move(x1, y1, x2, y2)
			} else {
				return fmt.Errorf("wrong move, %v can't move like that", f.name)
			}
		} else {
			if (y1-y2 == 1 || y1-y2 == 2 && y1 == 6) && f.color || (y2-y1 == 1 || (y2-y1 == 2 && y1 == 1)) && !f.color {
				b.move(x1, y1, x2, y2)
			} else {
				return fmt.Errorf("wrong move, %v can't move like that", f.name)
			}
		}
	case "rook":
		if x1 == x2 && y1 != y2 || x1 != x2 && y1 == y2 {
			b.move(x1, y1, x2, y2)
		} else {
			return fmt.Errorf("wrong move, %v can't move like that", f.name)
		}
	case "bishop":
		if (x1-x2)*(x1-x2) == (y1-y2)*(y1-y2) {
			b.move(x1, y1, x2, y2)
		} else {
			return fmt.Errorf("wrong move, %v can't move like that", f.name)
		}
	case "knight":
		if (x1-x2)*(x1-x2) == 1 && (y1-y2)*(y1-y2) == 4 || (x1-x2)*(x1-x2) == 4 && (y1-y2)*(y1-y2) == 1 {
			b.move(x1, y1, x2, y2)
		} else {
			return fmt.Errorf("wrong move, %v can't move like that", f.name)
		}
	case "king":
		if (x1-x2)*(x1-x2)+(y1-y2)*(y1-y2) <= 2 {
			b.move(x1, y1, x2, y2)
		} else {
			return fmt.Errorf("wrong move, %v can't move like that", f.name)
		}
	case "queen":
		if x1 == x2 && y1 != y2 || x1 != x2 && y1 == y2 || (x1-x2)*(x1-x2) == (y1-y2)*(y1-y2) {
			b.move(x1, y1, x2, y2)
		} else {
			return fmt.Errorf("wrong move, %v can't move like that", f.name)
		}
	default:
		return fmt.Errorf("internal error, unknown figure")
	}
	return nil
}

func main() {
	Board := InitBoard()
	fmt.Println(Board.GetStrField())

	for {
		var x1, x2 rune
		var y1, y2 int
		_, err := fmt.Scanf("%1c%1d%1c%1d\n", &x1, &y1, &x2, &y2)
		if err != nil {
			fmt.Println(err)
		}
		err = Board.Move(coordinates{x1, y1}, coordinates{x2, y2})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(Board.GetStrField())
	}
}
