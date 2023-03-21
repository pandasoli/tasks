package program

import (
	"os"
	"strings"
)


func RenderArgs() []string {
  var args []string

  if strings.Contains(os.Args[0], "/go-build") {
    for i, arg := range os.Args {
      if arg == "--" {
        args = os.Args[i + 1:]
        break
      }
    }
	} else {
		args = os.Args[1:]
	}

  return args
}
