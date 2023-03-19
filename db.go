package main

import (
  "io/ioutil"
  "os/user"

	"gopkg.in/yaml.v3"
)


func ReadAll() (map[string][]Task, error) {
  currentUser, err := user.Current()
  if err != nil { panic(err) }

  content, err := ioutil.ReadFile(currentUser.HomeDir + "/tasks.yml")

  data := make(map[string][]Task)

  err = yaml.Unmarshal(content, &data)

  return data, err
}

func ReadScoped(escope string) ([]Task, error) {
  data, err := ReadAll()
  if err != nil { return []Task {}, err }

  return data[escope], nil
}

func Write(data []Task, escope string) error {
  alldata, err := ReadAll()
  if err != nil { return err }

  if escope == "" {
    escope = "global"
  }

  alldata[escope] = data

  res, err := yaml.Marshal(alldata)
  if err != nil { return err }

  currentUser, err := user.Current()
  if err != nil { panic(err) }

  err = ioutil.WriteFile(currentUser.HomeDir + "/tasks.yml", res, 0644)

  return err
}
