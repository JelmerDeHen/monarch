package v4l2

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Finds first camera named name
// Look for name in /sys/class/video4linux/*/name
func findCameraByName(name string) (string, error) {
	///sys/class/video4linux/video0/name
	entries, err := ioutil.ReadDir("/sys/class/video4linux/")
	if err != nil {
		return "", err
	}

	if len(entries) == 0 {
		return "", fmt.Errorf("/sys/class/video4linux/ is empty")
	}

	for _, entry := range entries {
		namefn := fmt.Sprintf("/sys/class/video4linux/%s/name", entry.Name())
		_, err := os.Stat(namefn)
		if err != nil {
			continue
		}
		b, err := os.ReadFile(namefn)
		devname := strings.TrimSpace(string(b))
		if devname == name {
			return entry.Name(), nil
		}
	}
	return "", fmt.Errorf("Could not find camera named %q", name)
}

func testFindCameraByName() {
	name, err := findCameraByName("Logitech BRIO")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("/dev/" + name)
}
