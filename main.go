package main
import (
	"fmt"

	"tasks/term"
)


type Task struct {
  Title string
  Done bool
}

func main() {
  tasks, err := Read()
  if err != nil { panic(err) }

  // Set raw mode
  termios, err := term.SetRawMode()
  if err != nil { panic(err) }

  defer func() {
    err := term.RestoreMode(termios)
    if err != nil { panic(err) }
  }()

  // Render
  initial_y, _, err := term.WhereXY()
  if err != nil { panic(err) }
  initial_y++

  term.Block_blinking_cursor()

  makeSpace(tasks, &initial_y)
  render(tasks, initial_y, 0)

  // Main loop
  quit := false
  deleteds := []Task {}
  selected := 0

  for !quit {
    str := ""

    for str == "" {
      str, err = term.Getch()
      if err != nil { panic(err) }
    }

    switch str {
      case "q": quit = true
      case "\033[A": if selected > 0 { selected-- }
      case "\033[B": if selected < len(tasks) - 1 { selected++ }
      case "d":
        deleteds = append(deleteds, tasks[selected])
        tasks = append(tasks[:selected], tasks[selected + 1:]...)

        if selected == len(tasks) { selected-- }
      case "z":
        if len(deleteds) > 0 {
          last_task := deleteds[len(deleteds) - 1]
          tasks = append(tasks, last_task)
          deleteds = append(deleteds[:len(deleteds) - 1])
        }
      case "i": insert(&tasks, initial_y, &selected)
      case "u": update(&tasks, initial_y, selected)

      case "\n":
        tasks[selected].Done = !tasks[selected].Done
    }

    err := render(tasks, initial_y, selected)
    if err != nil { panic(err) }
  }

  // Save actions
  err = Write(tasks)
  if err != nil { panic(err) }

  // Go to the end
  term.GoToXY(0, initial_y + getNeededSpace(tasks))
  fmt.Println()
}
