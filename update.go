package main

import (
	"fmt"
	"github.com/pandasoli/goterm"
)


func update(tasks *[]Task, initial_y, selected int) {
  task := &(*tasks)[selected]
  y := initial_y + selected + 1 // +1 for top margin

  // I want to edit checked items without that line-through
  goterm.GoToXY(5, y)
  fmt.Print(task.Title)
  goterm.GoToXY(5, y)

  EditText(
    &task.Title,
    y,
    5,
    func(x int) {
      render(*tasks, initial_y, selected)

      goterm.GoToXY(5, y)
      fmt.Print(task.Title)
      goterm.GoToXY(5 + x, y)
    },
  )
}
