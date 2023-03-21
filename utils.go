package main


func containsString(s []string, search string) bool {
  for _, str := range s {
    if str == search {
      return true
    }
  }

  return false
}
