package term
import (
  "os"
  "strings"
)


func Getch() (string, error) {
  buff := make([]byte, 10)

  _, err := os.Stdin.Read(buff)
  if err != nil { return "", err }

  for range buff {
    for i, ch := range buff {
      if ch == 0 {
        buff = append(buff[:i], buff[i + 1:]...)
        break
      }
    }
  }

  // How the `buff` has always 10 positions
  // and I don't need them, I'm removing them
  res := string(buff)
  res = strings.TrimRight(res, "\x00")

  return res, nil
}
