package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
  . "tasks/program"

	"github.com/pandasoli/goterm"
)


func git() (string, error) {
  fi, err := os.Stat(".git")
  if err != nil { return "", err }

  if fi.IsDir() {
    cmd := exec.Command("git", "rev-parse", "--show-toplevel")
    out, err := cmd.Output()
    if err != nil { return "", err }

    repoName := string(out)
    repoName = strings.TrimSpace(repoName)
    repoName = path.Base(repoName)

    return repoName, nil
  }

  return "", fmt.Errorf(".git is not a directory")
}

func main() {
  var repoName []string
  args := RenderArgs()

  if len(args) > 0 {
    if args[0] == "help" {
      lines := []string {
        "",
        "     \033[1mTasks\033[0m",
        "",
        "[i] insert new item",
        "[u] update item",
        "[d] delete item",
        "[z] undo",
        "[enter] check item",
        "[arrow up] select top item",
        "[arrow down] select item and bass",
        "[q] save and exit",
        "",
      }

      for _, line := range lines {
        fmt.Println(line)
      }

      return
    } else {
      repoName = args
    }
  } else {
    repo, _ := git()
    repoName = []string { repo }
  }

  escopes, err := ReadScopes(repoName)
  if err != nil { panic(err) }

  // Set raw mode
  termios, err := goterm.SetRawMode()
  if err != nil { panic(err) }

  defer func() {
    err := goterm.RestoreMode(termios)
    if err != nil { panic(err) }
  }()

  // Some vars
  deleteds := []Escope {}
  selected := Selection {}

  // Render
  initial_y, _, err := goterm.WhereXY()
  if err != nil { panic(err) }

  goterm.Blinking_Block_cursor()

  MakeSpace(escopes, &initial_y)
  Render(escopes, initial_y, selected)

  // Features
  Features := map[string]func() {
    "GoUp": func() {
      if selected.Task > 0 {
        selected.Task--
      } else if selected.Escope > 0 {
        selected.Escope--
        selected.Task = len(escopes[selected.Escope].Tasks) - 1
      }
    },
    "GoDown": func() {
      if selected.Task < len(escopes[selected.Escope].Tasks) - 1 {
        selected.Task++
      } else if selected.Escope < len(escopes) - 1 {
        selected.Escope++
        selected.Task = 0
      }
    },
    "CheckTask": func() {
      if len(escopes) == 0 { return }

      task := &escopes[selected.Escope].Tasks[selected.Task]
      task.Done = !task.Done
    },
    "DeleteTask": func() {
      if len(escopes) == 0 { return }

      task := escopes[selected.Escope].Tasks[selected.Task]
      task_y := GetTaskY(escopes, initial_y, selected.Escope, selected.Task)

      // Animation
      animation_time := time.Duration(30)
      goterm.HideCursor()

      cb := "[ ]"
      if task.Done { cb = "[x]" }
      if len(task.Title) >= 30 { animation_time = time.Duration(10) }

      for i := range cb {
        goterm.GoToXY(5 + i, task_y)
        fmt.Printf("\033[31m%c\033[0m", cb[i])
        time.Sleep(time.Millisecond * animation_time)
      }

      time.Sleep(time.Millisecond * animation_time) // for the space between

      for i := range task.Title {
        style := "\033[31m"

        if task.Done {
          style += "\033[9m"
        }

        goterm.GoToXY(9 + i, task_y)
        fmt.Printf("%s%c\033[0m", style, task.Title[i])
        time.Sleep(time.Millisecond * animation_time)
      }

      // Do
      escope := escopes[selected.Escope]
      escope.Tasks = []Task { task }

      deleteds = append(deleteds, escope)
      escopes[selected.Escope].Tasks = append(
        escopes[selected.Escope].Tasks[:selected.Task],
        escopes[selected.Escope].Tasks[selected.Task + 1:]...,
      )

      if selected.Task == len(escopes[selected.Escope].Tasks) && len(escopes) > 0 {
        selected.Task--
      }

      // Animation
      for i := range task.Title {
        goterm.GoToXY(8 + len(task.Title) - i, task_y)
        fmt.Print(" ")
        time.Sleep(time.Millisecond * animation_time)
      }

      time.Sleep(time.Millisecond * animation_time) // for the space between

      for i := range cb {
        goterm.GoToXY(5 + len(cb) - i, task_y)
        fmt.Print(" ")
        time.Sleep(time.Millisecond * animation_time)
      }

      goterm.ShowCursor()
    },
    "RestoreTask": func() {
      if len(deleteds) == 0 { return }

      last_escope := deleteds[len(deleteds) - 1]
      task := last_escope.Tasks[0]
      escope_i := -1

      for i, escope := range escopes {
        if escope.Title == last_escope.Title {
          escope_i = i
          break
        }
      }

      if escope_i == -1 {
        escopes = append(escopes, last_escope)
      } else {
        escopes[escope_i].Tasks = append(escopes[escope_i].Tasks, task)
      }

      new_deletes := make([]Escope, 0, len(deleteds) - 1)
      new_deletes = append(new_deletes, deleteds[:len(deleteds) - 1]...)

      deleteds = new_deletes

      // Animation
      task_y := GetTaskY(escopes, initial_y, selected.Escope, len(escopes[escope_i].Tasks) - 1)
      animation_time := time.Duration(30)

      goterm.HideCursor()

      cb := "[ ]"
      if task.Done { cb = "[x]" }
      if len(task.Title) >= 30 { animation_time = time.Duration(10) }

      for i := range cb {
        goterm.GoToXY(5 + i, task_y)
        fmt.Printf("\033[34m%c\033[0m", cb[i])
        time.Sleep(time.Millisecond * animation_time)
      }

      time.Sleep(time.Millisecond * animation_time) // for the space between

      for i := range task.Title {
        style := "\033[37m"

        if task.Done {
          style += "\033[9m"
        }

        goterm.GoToXY(9 + i, task_y)
        fmt.Printf("%s%c\033[0m", style, task.Title[i])
        time.Sleep(time.Millisecond * animation_time)
      }

      goterm.ShowCursor()
    },
  }

  // Main loop
  quit := false

  for !quit {
    str, err := goterm.Getch()
    if err != nil { panic(err) }

    switch str {
      case "\033[A" /* Up arrow */: Features["GoUp"]()
      case "\033[B" /* Down arrow */: Features["GoDown"]()

      case "q": quit = true
      case "d": Features["DeleteTask"]()
      case "z": Features["RestoreTask"]()

      case "i": Insert(&escopes, &initial_y, &selected)
      case "u": Update(&escopes, &initial_y, selected)

      case "\n": Features["CheckTask"]()
    }

    err = Render(escopes, initial_y, selected)
    if err != nil { panic(err) }
  }

  // Save actions
  err = Write(escopes)
  if err != nil { panic(err) }

  // Go to the end
  goterm.GoToXY(0, initial_y + GetNeededSpace(escopes))
  fmt.Println()
}
