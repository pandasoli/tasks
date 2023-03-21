package main

import (
  "io/ioutil"
  "os/user"

	"gopkg.in/yaml.v3"
)


func ReadAll() ([]Escope, error) {
  currentUser, err := user.Current()
  if err != nil { panic(err) }

  content, err := ioutil.ReadFile(currentUser.HomeDir + "/tasks.yml")

  data := make([]Escope, 0, 254)
  err = yaml.Unmarshal(content, &data)

  return data, err
}

func ReadScopes(escopes []string) ([]Escope, error) {
  var res []Escope

  data, err := ReadAll()
  if err != nil { return nil, err }

  for _, escope := range data {
    if containsString(escopes, escope.Title) {
      res = append(res, escope)
    }
  }

  return res, nil
}

func Write(data []Escope) error {
  alldata, err := ReadAll()
  if err != nil { return err }

  for _, new_escope := range data {
    for i, old_scope := range alldata {
      if old_scope.Title == new_escope.Title {
        alldata[i].Tasks = new_escope.Tasks
        break
      }
    }
  }

  res, err := yaml.Marshal(alldata)
  if err != nil { return err }

  currentUser, err := user.Current()
  if err != nil { panic(err) }

  err = ioutil.WriteFile(currentUser.HomeDir + "/tasks.yml", res, 0644)

  return err
}
