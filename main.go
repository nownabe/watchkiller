package main

import (
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/nownabe/cenv"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()
	defer sugar.Sync()

	watchFile := cenv.MustString("KILLER_WATCH_FILE")
	interval := cenv.MustInt32("KILLER_INTERVAL")

	args := os.Args[1:]
	if len(args) == 0 {
		sugar.Fatal("A command is required. Usage: watchkiller <command> [args...]")
	}

	sugar.Infof("Watch file = %s", watchFile)
	sugar.Infof("Interval = %d", interval)
	sugar.Infof("Command = %s", args)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		sugar.Fatalf("Failed to execute %s. %s", args[0], err)
	}

	go func() {
		for {
			time.Sleep(time.Duration(interval) * time.Second)
			if _, err := os.Stat(watchFile); err == nil {
				sugar.Infof("Found %s! killing %s.", watchFile, args[0])
				cmd.Process.Kill()
				break
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				sugar.Errorf("Command %s exited with status code %d.", args[0], status.ExitStatus())
			}
		} else {
			sugar.Fatalf("Failed to wait command. %s", err)
		}
	} else {
		sugar.Infof("Command %s exited successfully.", args[0])
	}
}
