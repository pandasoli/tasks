package main
import (
  "tasks/term"
)


func lastSpaceIndex(str string) int {
  for i := len(str) - 1; i >= 0; i-- {
    if str[i] == ' ' {
      return i
    }
  }

  return 0
}

func EditText(text *string, x int, end_callback func(x int)) {
  slice := []rune(*text)

  for {
    str := ""

    for str == "" {
      str_, err := term.Getch()
      if err != nil { panic(err) }

      str = str_
    }

    if str == "\n" { // Enter
      break
    }

    switch str[0] {
    case 127: // Backspace
      if x > 0 {
        x--

        copy(slice[x:], slice[x + 1:])
        slice = slice[:len(slice)-1]
      }

    case 23:
      i := lastSpaceIndex(string(slice[:x]))

      left := slice[:i]
      slice = append(left, slice[x:]...)

      x = i

    case '\033': // Handle escape char
      if str[1] == '[' {
        switch str[2] {
          case 'C': if x < len([]rune(*text)) { x++ } // Right arrow key
          case 'D': if x > 0 { x-- } // Left arrow key
          case 'F': x = len([]rune(*text)) // End key
          case 'H': x = 0 // Home key
          case '3': // Delete key
            if str[3] == '~' {
              if x < len(slice) {
                copy(slice[x:], slice[x+1:])
                slice = slice[:len(slice)-1]
              }
            }
        }
      }

    default:
      slice = append(
        slice[:x],
        append([]rune(str), slice[x:]...)...
      )

      x++
    }

    *text = string(slice)
    end_callback(x)
  }
}
