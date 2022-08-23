package x11grab

import (
	"github.com/JelmerDeHen/scrnsaver"
	"fmt"
	"bytes"
	"os"
	"os/exec"
	"time"
)

func ffmpegArgv() []string {
	// Create out file
	hostname, err := os.Hostname()
	if err != nil {
    panic(err)
	}
	now := time.Now()
	outfile := fmt.Sprintf(
		"/data/mon/srec_new/%s.%02d%02d%02d.%02d%02d%02d.mkv",
		hostname,
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
	)

  display := os.Getenv("DISPLAY")
  if display == "" {
    panic("$DISPLAY was empty")
  }

	arguments := []string{
		"-nostdin", "-hide_banner",
		"-loglevel", "warning",
		"-f", "x11grab",
    "-an",
		"-framerate", "25",
    "-video_size", scrnsaver.GetResolution(),
		"-i", os.Getenv("DISPLAY"),
		"-vcodec", "libx265",
    "-preset", "ultrafast",
		outfile,
	}
	return arguments
}

func Run() {
	start := time.Now()
	var running bool
	var cmd *exec.Cmd

	for {
		// Get idle time
		info, err := scrnsaver.GetXScreenSaverInfo()
		if err != nil {
			panic(err)
		}

		// Determine if ffmpeg proc died
		if cmd == nil || (cmd != nil && cmd.ProcessState != nil && cmd.ProcessState.Exited()) {
			running = false
		}

		// Start ffmpeg when not idle
		if info.Idle < time.Minute {
			//if info.Idle < time.Second * 5 {
			if !running {
				ffmpegargv := ffmpegArgv()
				cmd = exec.Command("ffmpeg", ffmpegargv...)

				var stdout, stderr bytes.Buffer
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr

				go cmd.Run()

				//fmt.Printf("cmd: %+v\n", cmd)

				//fmt.Printf("cmd.ProcessState.Success(): %v\n", cmd.ProcessState.Success())
				//fmt.Printf("cmd.ProcessState.Exited(): %v\n", cmd.ProcessState.Exited())

				// Give ffmpeg some time to start or crash
				// Prevent ffmpeg from retrying to exec when crashing immediatly
				time.Sleep(time.Second * 1)
				if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
					err := fmt.Errorf("Could not start ffmpeg: errno=%d\n%s%s\n", cmd.ProcessState.ExitCode(), stdout.String(), stderr.String())
					panic(err)
				}

				start = time.Now()
				running = true
			}
		}

		// Elapsed time >1h rotate proc
		etime := time.Since(start)
		if etime > time.Hour {
			//if etime > time.Second * 20 {
			if running {
				//fmt.Println("Respawning")
				cmd.Process.Kill()
				running = false
				time.Sleep(time.Second * 1)
			}
		}

		// Kill ffmpeg when idle
		if info.Idle > time.Minute*10 {
			//if info.Idle > time.Second * 5 {
			if running && cmd != nil && cmd.ProcessState == nil {
				//fmt.Println("User idle")
				cmd.Process.Kill()
				running = false
				time.Sleep(time.Second * 1)
			}
		}

		if int(etime.Seconds())%60 == 1 {
			fmt.Printf("etime=%04d; info.Idle=%04d; running=%v\n", int(etime.Seconds()), int(info.Idle.Seconds()), running)
		}
		//    time.Sleep(time.Second)
		time.Sleep(time.Minute)
	}
}

