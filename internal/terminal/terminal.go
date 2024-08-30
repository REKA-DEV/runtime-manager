package terminal

import (
	"fmt"
	"github.com/REKA-DEV/runtime-manager/internal/ascii"
	"github.com/containerd/console"
)

type TColor = int

const (
	COLOR_BLACK   TColor = 30
	COLOR_RED     TColor = 31
	COLOR_GREEN   TColor = 32
	COLOR_YELLOW  TColor = 33
	COLOR_BLUE    TColor = 34
	COLOR_MAGENTA TColor = 35
	COLOR_CYAN    TColor = 36
	COLOR_WHITE   TColor = 37
)

type TCursorMove = string

const (
	CURSOR_MOVE_UP    TCursorMove = "A"
	CURSOR_MOVE_DOWN  TCursorMove = "B"
	CURSOR_MOVE_LEFT  TCursorMove = "C"
	CURSOR_MOVE_RIGHT TCursorMove = "D"
)

var sequenceMap = map[byte]func(b []byte, c console.Console) ([]byte, error){
	ascii.CHAR_lsqb: controlSequenceIntroducer,
}

func Color(c TColor) string {
	return fmt.Sprintf("\033[%dm", c)
}

func CursorMove(n int, m TCursorMove) string {
	return fmt.Sprintf("\033[%d%s", n, m)
}

func Read(c console.Console) ([]byte, error) {
	read := make([]byte, 0)
	b := make([]byte, 1)

	_, err := c.Read(b)
	if err != nil {
		return nil, err
	}

	read = append(read, b[0])

	if b[0] != ascii.ESC {
		return read, nil
	}

	_, err = c.Read(b)
	if err != nil {
		return nil, err
	}

	sequence, ok := sequenceMap[b[0]]
	if !ok {
		return nil, fmt.Errorf("invalid sequence: %d", b[0])
	}

	read = append(read, b[0])

	read, err = sequence(read, c)
	if err != nil {
		return nil, err
	}

	return read, nil
}

func controlSequenceIntroducer(read []byte, c console.Console) ([]byte, error) {
	b := make([]byte, 1)

	for {
		_, err := c.Read(b)
		if err != nil {
			return nil, err
		}

		if b[0] < ascii.CHAR_0 || ascii.CHAR_quest < b[0] {
			break
		}

		read = append(read, b[0])
	}

	if ascii.SP <= b[0] && b[0] <= ascii.CHAR_sol {
		read = append(read, b[0])
		for {
			_, err := c.Read(b)
			if err != nil {
				return nil, err
			}

			if b[0] < ascii.SP || ascii.CHAR_sol < b[0] {
				break
			}

			read = append(read, b[0])
		}
	}

	if b[0] < ascii.CHAR_commat || ascii.CHAR_tilde < b[0] {
		return nil, fmt.Errorf("invalid control sequence introducer sequence: %d", b[0])
	}

	read = append(read, b[0])

	return read, nil
}
