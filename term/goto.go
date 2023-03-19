package term
import "fmt"


func GoToXY(x, y int) {
  // The terminal's buffer's positions start at 1
  fmt.Printf("\033[%d;%dH", y, x + 1)
}
