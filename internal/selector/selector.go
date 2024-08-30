package selector

import (
	"fmt"
	"github.com/REKA-DEV/runtime-manager/internal/ascii"
	"github.com/REKA-DEV/runtime-manager/internal/terminal"
	"github.com/containerd/console"
	"unsafe"
)

var (
	UP    = string([]byte{ascii.ESC, ascii.CHAR_lsqb, ascii.CHAR_A})
	DOWN  = string([]byte{ascii.ESC, ascii.CHAR_lsqb, ascii.CHAR_B})
	RIGHT = string([]byte{ascii.ESC, ascii.CHAR_lsqb, ascii.CHAR_C})
	LEFT  = string([]byte{ascii.ESC, ascii.CHAR_lsqb, ascii.CHAR_D})
	CLEAR = string([]byte{ascii.ESC, ascii.CHAR_lsqb, ascii.CHAR_K})
	HIDE  = string([]byte{ascii.ESC, ascii.CHAR_lsqb, ascii.CHAR_quest, ascii.CHAR_2, ascii.CHAR_5, ascii.CHAR_l})
	SHOW  = string([]byte{ascii.ESC, ascii.CHAR_lsqb, ascii.CHAR_quest, ascii.CHAR_2, ascii.CHAR_5, ascii.CHAR_h})
)

type Selector[T any] struct {
	prompt string
	cursor int
	items  []Item[T]
}

type Item[T any] struct {
	text string
	data *T
}

func New[T any](prompt string) Selector[T] {
	return Selector[T]{
		prompt: prompt,
		cursor: 0,
		items:  make([]Item[T], 0),
	}
}

func (selector *Selector[T]) Add(text string, data *T) {
	selector.items = append(selector.items, Item[T]{
		text: text,
		data: data,
	})
}

func (selector *Selector[T]) drawItem(i int) {
	if i == selector.cursor {
		fmt.Printf("%s> %s", terminal.Color(terminal.COLOR_YELLOW), selector.items[i].text)
	} else {
		fmt.Printf("%s  %s", terminal.Color(terminal.COLOR_WHITE), selector.items[i].text)
	}
}

func (selector *Selector[T]) Move(n int) {
	if n == 0 {
		return
	}

	if n < 0 && !(selector.cursor+n < 0) {
		before := selector.cursor
		selector.cursor = selector.cursor + n
		for i := before; selector.cursor < i; i = i - 1 {
			selector.drawItem(i)
			fmt.Printf("%s%c", UP, ascii.CR)
		}
	}

	if 0 < n && !(len(selector.items) <= selector.cursor+n) {
		before := selector.cursor
		selector.cursor = selector.cursor + n
		for i := before; i < selector.cursor; i = i + 1 {
			selector.drawItem(i)
			fmt.Printf("%s%c", DOWN, ascii.CR)
		}
	}

	selector.drawItem(selector.cursor)
	fmt.Printf("%c", ascii.CR)
}

func (selector *Selector[T]) Run() (*T, error) {
	err := console.Current().SetRaw()
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s%c", selector.prompt, ascii.CR)

	for i := 0; i < len(selector.items); i++ {
		fmt.Printf("%c", ascii.LF)
	}

	for i := len(selector.items) - 1; 0 <= i; i = i - 1 {
		selector.drawItem(i)
		fmt.Printf("%s%c", UP, ascii.CR)
	}

	fmt.Printf("%s", DOWN)

	fmt.Printf("%s", HIDE)

	for {
		read, err2 := terminal.Read(console.Current())
		if err2 != nil {
			err = err2
			break
		}

		if len(read) == 1 && read[0] == ascii.ETX {
			err = fmt.Errorf("canceled")
			break
		}

		if len(read) == 1 && read[0] == ascii.CR {
			break
		}

		rstr := unsafe.String(&read[0], len(read))

		if rstr == UP {
			selector.Move(-1)
			continue
		}

		if rstr == DOWN {
			selector.Move(1)
			continue
		}
	}

	fmt.Printf("%s", SHOW)

	if err != nil {
		return nil, err
	}

	for i := selector.cursor; i <= len(selector.items); i = i + 1 {
		fmt.Printf("%s", DOWN)
	}

	return selector.items[selector.cursor].data, nil
}
