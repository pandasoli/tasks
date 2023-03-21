package program


type Task struct {
  Title string
  Done bool
}

type Escope struct {
  Title string
  Tasks []Task
}

type Selection struct {
  Escope,
  Task int
}
