package main
import (
  "fmt"
  "strings"

  "github.com/pandasoli/goterm"
)


var help_lines = []string {
  "[i] insert item",
  "[d] delete item",
  "[z] restore deleted",
  "[u] edit item",
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

func makeSpace(tasks []Task, initial_y *int) error {
  w, h, err := goterm.GetWinSize()
  if err != nil { return err }

  usable_h := h - *initial_y
  needed_h := getNeededSpace(tasks)

  if usable_h < needed_h {
    goterm.GoToXY(0, *initial_y)

    for range make([]int, needed_h - usable_h) {
      fmt.Println(
        strings.Repeat(" ", w),
      )

      *initial_y--
      usable_h++
    }

    //*initial_y -= needed_h - usable_h
    //usable_h = needed_h
  }

  return nil
}

func render(tasks []Task, initial_y int, selected int) error {
  w, _, err := goterm.GetWinSize()
  if err != nil { return err }

  // Calculate help instructions stuff
  larger_help_line := 0

  for _, line := range help_lines {
    if len(line) > larger_help_line {
      larger_help_line = len(line)
    }
  }

  // Clear
  h := getNeededSpace(tasks)

  for i := range make([]int, h) {
    goterm.GoToXY(0, initial_y + i)
    fmt.Println(
      strings.Repeat(" ", w),
    )
  }

  // Print
  for i, task := range tasks {
    goterm.GoToXY(1, initial_y + i)
    title := task.Title

    if task.Done {
      title = "\033[9m" + title
      fmt.Print("\033[1;34m[x]\033[0m")
    } else {
      fmt.Print("[ ]")
    }

    // The space after the title is for when I delete a char when I'm creating a item
    // thus I remove the char I just deleted from the screen.
    fmt.Printf(" %s\033[0m ", title)
  }

  // Clear line after the last item
  goterm.GoToXY(1, initial_y + len(tasks))
  fmt.Print(strings.Repeat(" ", w))

  // Show help instructions
  for i, line := range help_lines {
    goterm.GoToXY(w - larger_help_line - 1, initial_y + i)
    fmt.Printf("\033[90m%s\033[0m", line)
  }

  // Go to the selected item
  goterm.GoToXY(2, initial_y + selected)

  return nil
}
