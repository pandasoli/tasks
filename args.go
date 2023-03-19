package main

import (
	"os"
	"strings"
)


type Args struct {
  Escope string
}

func renderArgs() Args {
  var final_args Args
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

  final_args.Escope = args[0]

  return final_args
}
