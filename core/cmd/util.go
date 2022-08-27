package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// TODO: move me to better location
func getOutfilename(dir, ext string) string {
	hostname, _ := os.Hostname()
	now := time.Now()

	dir = strings.TrimSuffix(dir, "/")

	fn := fmt.Sprintf(
		"%s/%s.%02d%02d%02d.%02d%02d%02d.%s",
		dir,
		hostname,
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		ext,
	)

	return fn
}
