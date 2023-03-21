package program

import (
	"github.com/pandasoli/goterm"
)


func GetTaskY(escopes []Escope, initial_y, escopei, taski int) int {
  res := initial_y + 1 // for top margin

  for escope_i, escope := range escopes {
    if escope_i == escopei {
      res += taski + 1 // for title
      break
    } else {
      res += len(escope.Tasks) + 2 // For the title and spacing
    }
  }

  return res
}

func Insert(escopes *[]Escope, initial_y *int, selected *Selection) {
  tasks := &(*escopes)[selected.Escope].Tasks

  new_tasks := make([]Task, 0, len(*tasks) + 1)
  new_tasks = append(new_tasks, (*tasks)[:selected.Task]...)
  new_tasks = append(new_tasks, Task {})
  selected.Task = len(new_tasks) - 1
  new_tasks = append(new_tasks, (*tasks)[selected.Task:]...)

  *tasks = new_tasks
  task := &(*tasks)[selected.Task]

  MakeSpace(*escopes, initial_y)
  Render(*escopes, *initial_y, *selected)

  tasky := GetTaskY(*escopes, *initial_y, selected.Escope, selected.Task)

  goterm.GoToXY(10, tasky)

  EditText(
    &task.Title,
    tasky,
    9,
    func(x int) {
      Render(*escopes, *initial_y, *selected)
      goterm.GoToXY(10 + x, tasky)
    },
  )

  if len(task.Title) == 0 {
    new_tasks := make([]Task, 0, len(*escopes) - 1)
    new_tasks = append(new_tasks, (*tasks)[:selected.Task]...)
    new_tasks = append(new_tasks, (*tasks)[selected.Task + 1:]...)

    *tasks = new_tasks
  }
}
