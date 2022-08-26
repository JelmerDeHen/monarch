package compress

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func FfmpegCompress(infile string, outfile string, extraArgs []string) (err error) {
	args := []string{
		"-loglevel", "warning",
		"-i", infile,
	}
	args = append(args, extraArgs...)
	args = append(args, outfile)

	cmd := exec.Command("ffmpeg", args...)

	var stdcombined bytes.Buffer
	cmd.Stdout = &stdcombined
	cmd.Stderr = &stdcombined

	fmt.Printf("FfmpegCompress(): Exec: ffmpeg %s\n", strings.Join(args[:], " "))

	err = cmd.Run()
	if err != nil {
		return err
	}

	// Print stdout & stderr
	output := stdcombined.String()
	if output != "" {
		fmt.Printf("%s", stdcombined.String())
	}

	errno := cmd.ProcessState.ExitCode()
	if errno != 0 {
		return fmt.Errorf("FfmpegCompress(): exit code %d", errno)
	}

	return nil
}
