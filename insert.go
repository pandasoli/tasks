package main

import (
	"github.com/pandasoli/goterm"
)


func insert(tasks *[]Task, initial_y, selected *int) {
  new_tasks := make([]Task, 0, len(*tasks) + 1)
  new_tasks = append(new_tasks, (*tasks)[:*selected]...)
  new_tasks = append(new_tasks, Task {})
  *selected = len(new_tasks) - 1
  new_tasks = append(new_tasks, (*tasks)[*selected:]...)

  *tasks = new_tasks
  task := &(*tasks)[*selected]

  makeSpace(*tasks, initial_y)
  render(*tasks, *initial_y, *selected)

  goterm.GoToXY(5, *initial_y + *selected + 1) // +1 for top margin

  EditText(
    &task.Title,
    *initial_y + *selected + 1,
    5,
    func(x int) {
      render(*tasks, *initial_y, *selected)
      goterm.GoToXY(5 + x, *initial_y + *selected + 1) // +1 for top margin
    },
  )

  if len(task.Title) == 0 {
    new_tasks := make([]Task, 0, len(*tasks) - 1)
    new_tasks = append(new_tasks, (*tasks)[:*selected]...)
    new_tasks = append(new_tasks, (*tasks)[*selected + 1:]...)

    *tasks = new_tasks
  }
}
