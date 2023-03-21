package program

import (
	"fmt"
	"strings"

	"github.com/pandasoli/goterm"
	"golang.design/x/clipboard"
)


func lastSpaceIndex(str string) int {
  for i := len(str) - 1; i >= 0; i-- {
    if str[i] == ' ' {
      return i
    }
  }

  return 0
}

func ShowDebug(str string, w, left_x int) {
  goterm.GoToXY(left_x, 10)

  explicit_keys := strings.ReplaceAll(str, "\b", "\\b")
  explicit_keys = strings.ReplaceAll(explicit_keys, "\n", "\\n")
  explicit_keys = strings.ReplaceAll(explicit_keys, "\033", "\\033")

  final_str := fmt.Sprintf("\"%s\" %v (%d)", explicit_keys, []rune (str), len(str))

  fmt.Println(final_str)

  fmt.Println(strings.Repeat(" ", w))
  fmt.Println(strings.Repeat(" ", w))
}

func EditText(text_ *string, top, left_x int, end_callback func(x int)) {
  // Load some variables and constants
  w, _, err := goterm.GetWinSize()
  if err != nil { panic(err) }

  err = clipboard.Init()
  if err != nil { panic(err) }

  x := 0
  text := []rune(*text_)
  quit := false
  mode := "insert" // "replace" | "insert" | "select"
  selection_start := 0

  // Create events' functions
  Features := map[string]func() {
    "DeleteCurrentChar": func() {
      if x > 0 {
        x--

        copy(text[x:], text[x + 1:])
        text = text[:len(text) - 1]
      }
    },
    "DeleteNextChar": func() {
      if x < len(text) {
        copy(text[x:], text[x + 1:])
        text = text[:len(text) - 1]
      }
    },
    "DeleteCurrentWord": func() {
      i := lastSpaceIndex(string(text[:x]))

      left := text[:i]
      text = append(left, text[x:]...)

      x = i
    },
    "DeleteNextWord": func() {
      right := string(text[x:])
      trim := strings.TrimLeft(right, " ")
      spaces_count := len(right) - len(trim)

      i := strings.Index(trim, " ")

      if i == -1 {
        i = len(text[x:])
        spaces_count = 0
      }

      left := text[:x]
      text = append(left, text[x + spaces_count + i:]...)
    },
    "GoToLeft": func() {
      if x > 0 { x-- }
    },
    "GoToRight": func() {
      // These modes cannot pass the last character
      if
        (mode == "replace") &&
        x == len(text) - 1 {
        return
      }

      if x < len(text) {
        x++
      }
    },
    "GoToEnd": func() {
      x = len(text)

      // These modes cannot pass the last character
      if mode == "replace" {
        x--
      }
    },
    "GoToStart": func() {
      x = 0
    },
    "ComeBackOneWord": func() {
      x = lastSpaceIndex(
        strings.TrimRight(string(text[:x]), " "),
      )

      if x > 0 { x++ }
    },
    "GoOneWordForward": func() {
      right := string(text[x:])
      trim := strings.TrimLeft(right, " ")
      spaces_count := len(right) - len(trim)
      i := strings.Index(trim, " ")

      if i < 0 {
        x = len(text)
      } else {
        x += i + spaces_count + 1
      }
    },
  }


  ChangeMode := func(new_mode string) bool {
    if new_mode == "replace" && len(text) == 0 {
      return false
    }

    switch new_mode {
      case "replace":
        goterm.Underline_cursor()
        if x == len(text) && len(text) > 0 { x-- }

      case "insert": goterm.IBeam_cursor()
      case "select":
        selection_start = x
        goterm.Blinking_IBeam_cursor()
    }

    mode = new_mode
    return true
  }

  Type := func(str string) {
    switch mode {
      case "insert":
        text = append(
          text[:x],
          append([]rune(str), text[x:]...)...
        )

        x++

      case "replace":
        text = append(
          text[:x],
          append([]rune(str), text[x + 1:]...)...
        )
    }
  }

  TryChangeMode := func(new_mode, str string) {
    if !ChangeMode(new_mode) {
      Type(str)
    }
  }

  Copy := func(str string) {
    // Ensure the selection end is bigger than the start
    start := selection_start
    end := x

    if end < start {
      temp := end

      end = start
      start = temp
    }

    // Move to clipboard
    selected_text := text[start:end]
    clipboard.Write(
      clipboard.FmtText,
      []byte (string(selected_text)),
    )
  }

  HandleModeEvents := func(input string) {
    new_mode := "insert"
    after := func() {}

    switch input {
      case "\033" /* Esc */: ChangeMode("insert")
      case "\n" /* Enter */: quit = true

      case "\033[A" /* Up arrow */: Features["GoToStart"]()
      case "\033[B" /* Down arrow */: Features["GoToEnd"]()
      case "\033[C" /* Right arrow */: Features["GoToRight"]()
      case "\033[D" /* Left arrow */: Features["GoToLeft"]()

      case "\033[1;2C" /* Shift + Right arrow */:
        new_mode = "select"
        old_selection_start := selection_start
        already_in_mode := mode == "select"

        after = func() {
          Features["GoToRight"]()

          if already_in_mode {
            selection_start = old_selection_start
          }
        }

      case "\033[1;2D" /* Shift + Left arrow */:
        new_mode = "select"
        old_selection_start := selection_start
        already_in_mode := mode == "select"

        after = func() {
          Features["GoToLeft"]()

          if already_in_mode {
            selection_start = old_selection_start
          }
        }

      case "\033[H" /* Home */: Features["GoToStart"]()
      case "\033[F" /* End */: Features["GoToEnd"]()

      case "\033[2~" /* Insert */: TryChangeMode("replace", input)

      default:
        switch mode {
          case "insert":
            switch input {
              case "\x7f" /* Backspace */: Features["DeleteCurrentChar"]()
              case "\b" /* Ctrl + Backspace */: Features["DeleteCurrentWord"]()
              case "\x17" /* Ctrl + w */: Features["DeleteCurrentWord"]()

              case "\033[1;5C" /* Ctrl + Right arrow */: Features["GoOneWordForward"]()
              case "\033[1;5D" /* Ctrl + Left arrow */: Features["ComeBackOneWord"]()

              case "\033[3~" /* Delete */: Features["DeleteNextChar"]()
              case "\033[3;5~" /* Ctrl + Delete */: Features["DeleteNextWord"]()

              default:
                Type(input)
            }

          case "select":
            switch input {
              case "y":
                Copy(input)
                ChangeMode("normal")
            }

          default:
            Type(input)
        }
    }

    if new_mode != mode {
      ChangeMode(new_mode)
    }

    after()
  }

  // Initialize somethings
  goterm.ShowCursor()
  ChangeMode("insert")

  // The main loop
  for !quit {
    // Change the background color if in select mode
    final_text := make([]rune, len(text))
    copy(final_text, text)

    if mode == "select" {
      // Ensure the selection end is bigger than the start
      start := selection_start
      end := x

      if end < start {
        temp := end

        end = start
        start = temp
      }

      // Put colors
      final_text = append(
        final_text[:end],
        append([]rune("\033[0m"), final_text[end:]...)...
      )

      final_text = append(
        final_text[:start],
        append([]rune("\033[40m"), final_text[start:]...)...
      )
    }

    // Print
    goterm.GoToXY(left_x, top)

    fmt.Print(
      string(final_text) +
      strings.Repeat(" ", w - left_x - len(text)), // Remove the deleted character from the screen
    )

    // Come back to the cursor position
    goterm.GoToXY(left_x + x, top)

    // Get char and handle it
    str, err := goterm.Getch()
    if err != nil { panic(err) }

    HandleModeEvents(str)
    *text_ = string(text)
  }

  // Make stuff like they were before
  goterm.Block_cursor()
}
