package term
import (
  "fmt"
	"os"
)


func Getch() (ch int, err error) {
  buff := make([]byte, 1)
  _, err = os.Stdin.Read(buff)

  if err != nil {
    fmt.Println("Error reading from stdin:", err)
    return 0, err
  }

  return int(buff[0]), nil
}
