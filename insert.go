package main
import (
  "tasks/term"
)


func insert(tasks *[]Task, initial_y int, selected *int) {
  new_tasks := make([]Task, 0, len(*tasks) + 1)
  new_tasks = append(new_tasks, (*tasks)[:*selected + 1]...)
  new_tasks = append(new_tasks, Task {})
  new_tasks = append(new_tasks, (*tasks)[*selected + 1:]...)

  *tasks = new_tasks
  task := &(*tasks)[*selected + 1]

  makeSpace(*tasks, &initial_y)
  render(*tasks, initial_y, *selected)

  term.GoToXY(5, initial_y + *selected + 1)

  EditText(
    &task.Title,
    0,
    func(x int) {
      render(*tasks, initial_y, *selected + 1)
      term.GoToXY(5 + x, initial_y + *selected + 1)
    },
  )

  if len(task.Title) == 0 {
    new_tasks := make([]Task, 0, len(*tasks) - 1)
    new_tasks = append(new_tasks, (*tasks)[:*selected + 1]...)
    new_tasks = append(new_tasks, (*tasks)[*selected + 2:]...)

    *tasks = new_tasks
  } else {
    *selected++
  }
}
