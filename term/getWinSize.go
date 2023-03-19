package term
import "golang.org/x/sys/unix"


func GetWinSize() (width, height int, err error) {
  ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
  return int(ws.Col), int(ws.Row), err
}
