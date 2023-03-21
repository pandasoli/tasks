package program

import (
	"fmt"
	"github.com/pandasoli/goterm"
)


func Update(escopes *[]Escope, initial_y *int, selected Selection) {
  task := &(*escopes)[selected.Escope].Tasks[selected.Task]
  y := GetTaskY(*escopes, *initial_y, selected.Escope, selected.Task)

  // I want to edit checked items without that line-through
  goterm.GoToXY(10, y)
  fmt.Print(task.Title)

  EditText(
    &task.Title,
    y,
    9,
    func(x int) {
      Render(*escopes, *initial_y, selected)

      goterm.GoToXY(10, y)
      fmt.Print(task.Title)
      goterm.GoToXY(10 + x, y)
    },
  )
}
