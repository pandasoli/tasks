package term
import (
  "fmt"
  "os"
)


func WhereXY() (row, col int, err error) {
  // Send the ANSI command to get the cursor position.
  if _, err := os.Stdout.Write([]byte("\033[6n")); err != nil {
    return 0, 0, err
  }

  // Read the response from stdin.
  var buf [16]byte
  n, err := os.Stdin.Read(buf[:])
  if err != nil {
    return 0, 0, err
  }

  // Parse the cursor position from the response.
  if _, err := fmt.Sscanf(string(buf[:n]), "\033[%d;%dR", &row, &col); err != nil {
    return 0, 0, err
  }

  return row, col, nil
}
