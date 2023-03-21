package main
import (
  "fmt"
  "strings"

  "github.com/pandasoli/goterm"
)

func getNeededSpace(tasks []Task) int {
  var needed_h int

  needed_h += len(tasks)
  needed_h += 2 // Y space

  return needed_h
}

func makeSpace(tasks []Task, initial_y *int) error {
  _, h, err := goterm.GetWinSize()
  if err != nil { return err }

  needed_h := getNeededSpace(tasks)

  if h < needed_h {
    return fmt.Errorf("There's not the needed height. Needed: %d, have: %d", needed_h, h)
  }

  had_space := h - *initial_y

  if needed_h > had_space {
    missing_space := needed_h - had_space

    goterm.GoToXY(0, *initial_y)
    for range make([]int, had_space) {
      fmt.Println()
    }

    for range make([]int, missing_space) {
      fmt.Println()
      *initial_y--
    }
  }

  return nil
}

func render(tasks []Task, initial_y int, selected int) error {
  w, _, err := goterm.GetWinSize()
  if err != nil { return err }

  // Clear
  h := getNeededSpace(tasks)
  goterm.GoToXY(0, initial_y)

  for range make([]int, h) {
    fmt.Println(
      strings.Repeat(" ", w),
    )
  }

  goterm.GoToXY(0, initial_y)

  // Add margin at top
  fmt.Println()

  // Print
  for _, task := range tasks {
    title := task.Title

    if task.Done {
      title = "\033[9m" + title
      fmt.Print(" \033[1;34m[x]\033[0m")
    } else {
      fmt.Print(" [ ]")
    }

    fmt.Printf(" %s\033[0m\n", title)
  }

  // Go to the selected item
  // +1 because of the margin
  goterm.GoToXY(2, initial_y + selected + 1)

  return nil
}
