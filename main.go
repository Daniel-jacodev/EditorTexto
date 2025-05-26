package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	lines  *List[*List[rune]]
	line   *Node[*List[rune]]
	cursor *Node[rune]
	screen tcell.Screen
	style  tcell.Style
}

func (e *Editor) InsertChar(r rune) {
	e.cursor = e.line.Value.Insert(e.cursor, r)
	e.cursor = e.cursor.Next()
}

func (e *Editor) KeyLeft() {
	if e.cursor != e.line.Value.Front() { // Se o cursor não está no início da linha
		e.cursor = e.cursor.Prev() // Move o cursor para a esquerda
		return
	}
	// Estamos no início da linha
	if e.line != e.lines.Front() { // Se não está na primeira linha
		e.line = e.line.Prev()        // Move para a linha anterior
		e.cursor = e.line.Value.End() // Move o cursor para o final da linha
	}
}

func (e *Editor) KeyEnter() {
	if e.cursor == e.line.Value.End() {
		e.lines.Insert(e.line.Next(), NewList[rune]())
		e.line = e.line.Next()
		e.cursor = e.line.Value.Front()
		return
	} else {
		novalinha := NewList[rune]()
		e.lines.Insert(e.line.Next(), novalinha)
		atual := e.cursor
		linhaMax := e.line.Value.End()
		// e.line = e.line.Next()
		// e.cursor = e.line.Value.Front()
		for atual != linhaMax {
			novalinha.Insert(novalinha.End(), atual.Value)
			e.line.Value.Remove(atual)
			atual = atual.Next()
		}
		e.line = e.line.Next()
		e.cursor = e.line.Value.Front()
		// e.line.Value.Remove(e.cursor.Next())
		// e.cursor = e.cursor.Next()
	}
}

func (e *Editor) KeyRight() {
	if e.cursor != e.line.Value.End() {
		e.cursor = e.cursor.Next()
		return
	}
	if e.line != e.lines.End() {
		e.line = e.line.Next()
		e.cursor = e.line.Value.Front()
	}
}

func (e *Editor) KeyUp() {
	if e.line == e.lines.Front() {
		e.line = e.lines.End()
		e.cursor = e.line.Value.Front()
		return
	}
	e.line = e.line.Prev()
	e.cursor = e.line.Value.Front()
}

func (e *Editor) KeyDown() {
	if e.line == e.lines.End() {
		e.line = e.lines.Front()
		e.cursor = e.line.Value.Front()
	}
	e.line = e.line.Next()
	e.cursor = e.line.Value.Front()
}
func (e *Editor) KeyBackspace() {

	if e.cursor != e.line.Value.Front() {
		e.cursor = e.cursor.Prev()
		e.line.Value.Remove(e.cursor)
		return
	}

	if e.line != e.lines.Front() {
		linhaAtual := e.line
		e.line = e.line.Prev()
		cursorAntigo := e.line.Value.Back()

		for node := linhaAtual.Value.Front(); node != linhaAtual.Value.End(); {
			next := node.Next()
			e.line.Value.Insert(e.line.Value.End(), node.Value)
			linhaAtual.Value.Remove(node)
			node = next
		}
		e.lines.Remove(linhaAtual)
		e.cursor = cursorAntigo.Next()
	}
}

func (e *Editor) KeyDelete() {
	if e.cursor != e.line.Value.Back() {
		cursoratual := e.cursor
		e.line.Value.Remove(e.cursor.Next())
		e.cursor = cursoratual.Next()
		return
	}
	if e.line != e.lines.End() {
		linhaAtual := e.line
		e.line = e.line.Next()
		// cursorMexido := e.line.Value.Back()
		for node := e.line.Value.Front(); node != e.line.Value.Back(); {
			linhaAtual.Value.Insert(linhaAtual.Value.End(), node.Value)
			e.line.Value.Remove(node)
			node = node.Next()
		}
		e.lines.Remove(e.line)
		e.line = linhaAtual
		e.cursor = e.line.Value.Back()

	}
}

func main() {
	// Texto inicial e posição do cursor
	editor := NewEditor()
	editor.Draw()
	editor.MainLoop()
	defer editor.screen.Fini() // Encerra a tela ao sair
}

func (e *Editor) MainLoop() {
	for {
		ev := e.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEsc, tcell.KeyCtrlC:
				return
			case tcell.KeyEnter:
				e.KeyEnter()
			case tcell.KeyLeft:
				e.KeyLeft()
			case tcell.KeyRight:
				e.KeyRight()
			case tcell.KeyUp:
				e.KeyUp()
			case tcell.KeyDown:
				e.KeyDown()
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				e.KeyBackspace()
			case tcell.KeyDelete:
				e.KeyDelete()
			default:
				if ev.Rune() != 0 {
					e.InsertChar(ev.Rune())
				}
			}
			e.Draw()
		case *tcell.EventResize:
			e.screen.Sync()
			e.Draw()
		}
	}
}

func NewEditor() *Editor {
	e := &Editor{}
	// Inicializa a tela
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("erro ao criar a tela: %v", err)
	}
	if err := screen.Init(); err != nil {
		fmt.Printf("erro ao iniciar a tela: %v", err)
	}
	e.screen = screen
	e.lines = NewList[*List[rune]]()
	e.lines.PushBack(NewList[rune]())
	e.line = e.lines.Front()
	e.cursor = e.line.Value.Back()
	// Define o estilo do texto (branco com fundo preto)
	e.style = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	// Limpa a tela e define o estilo base
	e.screen.SetStyle(e.style)
	e.screen.Clear()
	return e
}

func (e *Editor) Draw() {
	e.screen.Clear()
	x := 0
	y := 0
	for line := e.lines.Front(); line != e.lines.End(); line = line.Next() {
		for char := line.Value.Front(); ; char = char.Next() {
			data := char.Value
			if char == line.Value.End() {
				data = '⤶'
			}
			if data == ' ' {
				data = '·'
			}
			if char == e.cursor {
				e.screen.SetContent(x, y, data, nil, e.style.Reverse(true))
			} else {
				e.screen.SetContent(x, y, data, nil, e.style)
			}
			x++
			if char == line.Value.End() {
				break
			}
		}
		y++
		x = 0
	}
	e.screen.Show()
}
