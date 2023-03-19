package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pandasoli/goterm"
)


type Task struct {
  Title string
  Done bool
}

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
  var repoName string
  args := renderArgs()

  repoName = args.Escope

  if repoName == "" {
    repoName, _ = git()
  }

  tasks, err := ReadScoped(repoName)
  if err != nil { panic(err) }

  // Set raw mode
  termios, err := goterm.SetRawMode()
  if err != nil { panic(err) }

  defer func() {
    err := goterm.RestoreMode(termios)
    if err != nil { panic(err) }
  }()

  // Render
  initial_y, _, err := goterm.WhereXY()
  if err != nil { panic(err) }
  initial_y++

  goterm.Blinking_Block_cursor()

  makeSpace(tasks, &initial_y)
  render(tasks, initial_y, 0)

  // Main loop
  quit := false
  deleteds := []Task {}
  selected := 0

  for !quit {
    str := ""

    for str == "" {
      str, err = goterm.Getch()
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
        if len(tasks) > 0 {
          tasks[selected].Done = !tasks[selected].Done
        }
    }

    err := render(tasks, initial_y, selected)
    if err != nil { panic(err) }
  }

  // Save actions
  err = Write(tasks, repoName)
  if err != nil { panic(err) }

  // Go to the end
  goterm.GoToXY(0, initial_y + getNeededSpace(tasks))
  fmt.Println()
}
