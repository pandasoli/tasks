package term
import "fmt"


// HideCursor hides the console cursor.
func HideCursor() {
  fmt.Print("\033[?25l")
}

// ShowCursor shows the console cursor.
func ShowCursor() {
  fmt.Print("\033[?25h")
}

func IBeam_cursor() { fmt.Print("\033[6 q") }
func Block_cursor() { fmt.Print("\033[2 q") }
func Underline_cursor() { fmt.Print("\033[4 q") }

func IBeam_blinking_cursor() { fmt.Print("\033[5 q") }
func Block_blinking_cursor() { fmt.Print("\033[1 q") }
func Underline_blinking_cursor() { fmt.Print("\033[3 q") }
