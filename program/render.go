package program

import (
  "fmt"
  "strings"

  "github.com/pandasoli/goterm"
)


func GetNeededSpace(escopes []Escope) int {
  var needed_h int

  needed_h += len(escopes) // Escopes' title
  needed_h += 2 // Y space

  for _, escope := range escopes {
    needed_h += len(escope.Tasks) + 1 // +1 for spacing
  }

  return needed_h
}

func MakeSpace(escopes []Escope, initial_y *int) error {
  _, h, err := goterm.GetWinSize()
  if err != nil { return err }

  needed_h := GetNeededSpace(escopes)

  if h < needed_h {
    return fmt.Errorf("There's not the needed height. Needed: %d, have: %d", needed_h, h)
  }

  had_space := h - *initial_y

  if needed_h > had_space {
    missing_space := needed_h - had_space

    // goterm.GoToXY(0, *initial_y)
    // for range make([]int, had_space) {
    //   fmt.Println()
    // }

    goterm.GoToXY(0, h)
    for range make([]int, missing_space) {
      fmt.Println()
    }

    *initial_y -= missing_space
  }

  return nil
}

func Render_task(task Task, cb_cl, tl_cl int) {
  title := task.Title
  cb := fmt.Sprintf("    \033[3%dm[ ]\033[0m", cb_cl)

  if task.Done {
    title = fmt.Sprintf("\033[3%d;9m%s", tl_cl, title)
    cb = fmt.Sprintf("    \033[3%dm[x]\033[0m", cb_cl)
  }

  fmt.Printf(" %s %s\033[0m\n", cb, title)
}

func Render(escopes []Escope, initial_y int, selected Selection) error {
  w, _, err := goterm.GetWinSize()
  if err != nil { return err }

  // Clear
  h := GetNeededSpace(escopes)
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
  for _, escope := range escopes {
    fmt.Printf("  \033[1;34m%s\033[0m\n", escope.Title)

    for _, task := range escope.Tasks {
      Render_task(task, 4, 7)
    }

    fmt.Println()
  }

  // Go to the selected item
  selection_y := GetTaskY(escopes, initial_y, selected.Escope, selected.Task)
  goterm.GoToXY(6, selection_y)

  return nil
}
