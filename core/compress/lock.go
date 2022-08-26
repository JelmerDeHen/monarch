package compress

import (
	"fmt"
	"os"
)

// Create empty files to indicate if file is used
type Lock struct {
	// Filename of the lockfile
	Name string
}

func (l *Lock) String() string {
	return l.Name
}

func (l *Lock) Lock() error {
	if l.IsLocked() {
		return fmt.Errorf("%s already locked!", l.Name)
	}

	f, err := os.Create(l.Name)
	if err != nil {
		return err
	}

	defer f.Close()

	return nil
}

func (l *Lock) Release() error {
	return os.Remove(l.Name)
}

func (l *Lock) IsLocked() bool {
	if _, err := os.Stat(l.Name); err == nil {
		return true
	}

	return false
}

func NewLock(name string) *Lock {
	return &Lock{
		Name: name,
	}
}

/*
  err := testLock()
  if err != nil {
    fmt.Println(err)
  }
func testLock () error {
  name := "/tmp/locktest.lock"

  l := NewLock(name)
  err := l.Lock()
  if err != nil {
    return fmt.Errorf("[-] Could not create lockfile %s: %s", l.Name, err)
  }

  fmt.Printf("[+] Created lockfile %s\n", l.Name)

  err = l.Release()
  if err != nil {
    return fmt.Errorf("[-] Could not release lockfile %s: %s", l.Name, err)
  }
  fmt.Printf("[+] Released lockfile %s\n", l.Name)

  return nil
}
*/
