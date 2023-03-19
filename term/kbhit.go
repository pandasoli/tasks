package term
import (
  "syscall"
  "unsafe"
)


func KBHit() bool {
  var oldt syscall.Termios
  if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&oldt))); err != 0 {
    return false
  }

  newt := oldt
  newt.Lflag &^= (syscall.ICANON | syscall.ECHO)
  if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&newt))); err != 0 {
    return false
  }

  defer func() {
    syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&oldt)))
  }()

  oldf, _, err := syscall.Syscall(syscall.SYS_FCNTL, uintptr(syscall.Stdin), syscall.F_GETFL, 0)
  if err != 0 { return false }

  if _, _, err := syscall.Syscall(syscall.SYS_FCNTL, uintptr(syscall.Stdin), syscall.F_SETFL, oldf|syscall.O_NONBLOCK); err != 0 {
    return false
  }

  defer func() {
    syscall.Syscall(syscall.SYS_FCNTL, uintptr(syscall.Stdin), syscall.F_SETFL, oldf)
  }()

  var buf [1]byte

  n, _ := syscall.Read(syscall.Stdin, buf[:])
  if n == 1 { return true }

  return false
}
