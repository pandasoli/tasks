package program


func ContainsString(s []string, search string) bool {
  for _, str := range s {
    if str == search {
      return true
    }
  }

  return false
}
