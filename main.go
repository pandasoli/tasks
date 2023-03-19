package main

import (
	"fmt"
	"os"
	"strings"

	"tasks/term"

	"gopkg.in/yaml.v3"
)


type Task struct {
  Title string `yaml:title`
  Done bool    `yaml:done`
}

func Read() ([]Task, error) {
  file, err := os.Open("tasks.yml")
  if err != nil { return nil, err }
  defer file.Close()

  content := make([]byte, 1024)
  count, err := file.Read(content)
  if err != nil { return nil, err }

  data := []Task {}
  err = yaml.Unmarshal(content[:count], &data)

  return data, err
}

func Write(data []Task) error {
  res, err := yaml.Marshal(data)
  if err != nil { return err }

  file, err := os.Create("tasks.yml")
  if err != nil { return err }
  defer file.Close()

  _, err = file.Write(res)
  if err != nil { return err }

  return nil
}

var help_lines = []string {
  "[i] insert item",
  "[d] delete item",
  "[z] restore deleted",
  "[enter] check item",
  "[q] save and quit",
}

func getNeededSpace(tasks []Task) int {
  // +1 because I clean the line after the last item
  needed_h := len(tasks) + 1

  if needed_h < len(help_lines) {
    needed_h = len(help_lines)
  }

  return needed_h
}

func reload(tasks []Task, initial_y *int, selected int) error {
  w, h, err := term.GetWinSize()
  if err != nil { return err }

  // Calculate help instructions stuff
  larger_help_line := 0

  for _, line := range help_lines {
    if len(line) > larger_help_line {
      larger_help_line = len(line)
    }
  }

  // Make space
  usable_h := h - *initial_y
  needed_h := getNeededSpace(tasks)

  if usable_h < needed_h {
    term.GoToXY(0, *initial_y)

    for range make([]int, needed_h) {
      fmt.Println(
        strings.Repeat(" ", w),
      )
    }

    *initial_y -= needed_h - usable_h
    usable_h += needed_h - usable_h
  }

  // Print
  for i, task := range tasks {
    bgcl := 0
    title := task.Title

    if task.Done {
      bgcl = 44
      title = "\033[9m" + title
    }

    // The space after the title is for when I delete a char when I'm creating a item
    // thus I remove the char I just deleted from the screen.
    term.GoToXY(1, *initial_y + i)
    fmt.Printf("[\033[%dm \033[0m] %s\033[0m \n", bgcl, title)
  }

  // Clear line after the last item
  term.GoToXY(1, *initial_y + len(tasks))
  fmt.Print(strings.Repeat(" ", w))

  // Show help instructions
  for i, line := range help_lines {
    term.GoToXY(w - larger_help_line - 1, *initial_y + i)
    fmt.Printf("\033[90m%s\033[0m", line)
  }

  // Go to the selected item
  term.GoToXY(2, *initial_y + selected)

  return nil
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

  initial_y++ // Give a space on the start

  term.Block_blinking_cursor()
  reload(tasks, &initial_y, 0)

  // Main loop
  quit := false
  deleteds := []Task {}
  selected := 0

  for !quit {
    key := 0

    for key == 0 || key == 0xe0 {
      key, err = term.Getch()
      if err != nil { panic(err) }
    }

    switch key {
      case 'q': quit = true
      case 'A': if selected > 0 { selected-- }
      case 'B': if selected < len(tasks) - 1 { selected++ }
      case 'd':
        deleteds = append(deleteds, tasks[selected])
        tasks = append(tasks[:selected], tasks[selected + 1:]...)

        if selected == len(tasks) { selected-- }
      case 'z':
        if len(deleteds) > 0 {
          last_task := deleteds[len(deleteds) - 1]
          tasks = append(tasks, last_task)
          deleteds = append(deleteds[:len(deleteds) - 1])
        }
      case 'i':
        tasks = append(tasks, Task {})
        task := &tasks[len(tasks) - 1]
        reload(tasks, &initial_y, selected)

        term.GoToXY(4, initial_y + len(tasks) - 1)

        key := 0

        for {
          key = 0

          for key == 0 || key == 0xe0 {
            key, err = term.Getch()
            if err != nil { panic(err) }
          }

          if key == 10 { // Enter
            break
          } else if key == 127 { // Backspace
            if len(task.Title) > 0 {
              task.Title = task.Title[:len(task.Title) - 1]
            }
          } else {
            task.Title += string(byte(key))
          }

          reload(tasks, &initial_y, selected)
          term.GoToXY(4 + len(task.Title), initial_y + len(tasks) - 1)
        }

        if len(task.Title) == 0 {
          tasks = append(tasks[:len(tasks) - 1])
        } else {
          selected = len(tasks) - 1
        }
      case 10:
        tasks[selected].Done = !tasks[selected].Done
    }

    err := reload(tasks, &initial_y, selected)
    if err != nil { panic(err) }
  }

  // Save actions
  err = Write(tasks)
  if err != nil { panic(err) }

  // Go to the end
  _, h, err := term.GetWinSize()
  y := initial_y + getNeededSpace(tasks) + 1

  for y > h - 3 {
    fmt.Println()
    h++
  }

  term.GoToXY(0, initial_y + getNeededSpace(tasks) + 1)
}
