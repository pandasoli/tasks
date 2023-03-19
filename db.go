package main
import (
  "os"
	"gopkg.in/yaml.v3"
)


func Read() ([]Task, error) {
  file, err := os.Open("tasks.yml")
  if err != nil { return nil, err }
  defer file.Close()

  content := make([]byte, 1024)
  count, err := file.Read(content)
  if err != nil { return nil, err }

  data := []Task {}
  err = yaml.Unmarshal(content[:count], &data)

  return data, err
}

func Write(data []Task) error {
  res, err := yaml.Marshal(data)
  if err != nil { return err }

  file, err := os.Create("tasks.yml")
  if err != nil { return err }
  defer file.Close()

  _, err = file.Write(res)
  if err != nil { return err }

  return nil
}
