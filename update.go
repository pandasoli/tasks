package main

import (
	"fmt"
	"tasks/term"
)


func update(tasks *[]Task, initial_y, selected int) {
  task := &(*tasks)[selected]
  y := initial_y + selected

  // I want to edit checked items without that line-through
  term.GoToXY(5, y)
  fmt.Print(task.Title)
  term.GoToXY(5, y)

  EditText(
    &task.Title,
    0,
    func(x int) {
      render(*tasks, initial_y, selected)

      term.GoToXY(5, y)
      fmt.Print(task.Title)
      term.GoToXY(5 + x, y)
    },
  )
}
