package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var command_name = os.Args[0]
var help_requested = flag.Bool("h", false, "Show this help screen.")
var debug_mode = flag.Bool("d", false, "Show debug information.")
var wait_time = flag.Int("w", 1, "Wait time between trials in seconds.")

var show_command_does_not_exist_warning = true

func show_usage() {
	fmt.Fprintf(os.Stderr, "Keeps executing a command while or until it returns successfully.\n\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "    %s while|until command [command arguments...]\n\n", command_name)
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "    -h: Show this help screen\n")
	fmt.Fprintf(os.Stderr, "    -d: Show debug information\n")
	fmt.Fprintf(os.Stderr, "    -w time: Wait time between trials in seconds\n")
	os.Exit(0)
}

func log(format_string string, arguments ...interface{}) {
	fmt.Fprintf(os.Stderr, format_string, arguments...)
}

func debug_log(format_string string, arguments ...interface{}) {
	if *debug_mode {
		log(format_string, arguments...)
	}
}

func die_on_error(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(1)
}

func wait() {
	debug_log("- Trying again in %d seconds...\n", *wait_time)
	time.Sleep(time.Duration(*wait_time) * time.Second)
}

func get_exit_code(exit_status error) int {
	if exit_status == nil {
		return 0
	} else {
		return exit_status.(*exec.ExitError).Sys().(syscall.WaitStatus).ExitStatus()
	}
}
func execs_okay(command_line string) bool {
	shell := os.Getenv("SHELL")
	shell_command := exec.Command(shell, "-c", command_line)
	debug_log("- Executing: %s", command_line)
	exit_code := get_exit_code(shell_command.Run())
	debug_log(" -> %d.\n", exit_code)
	if exit_code == 127 && show_command_does_not_exist_warning {
		log("Warning: one or more commands may not exist.\n")
		show_command_does_not_exist_warning = false
	}
	return exit_code == 0
}

func try_while(command_line string) {
	debug_log("Trying until command fails...\n")
	for execs_okay(command_line) {
		wait()
	}
	debug_log("Command has finally failed.\n")
}

func try_until(command_line string) {
	debug_log("Trying until command succeeds...\n")
	for !execs_okay(command_line) {
		wait()
	}
	debug_log("Command has finally succeeded.\n")
}

func main() {
	flag.Parse()
	predicate := flag.Args()
	// We must have at least a verb an a complement.
	if len(predicate) == 0 {
		show_usage()
	} else if len(predicate) < 2 {
		die_on_error(fmt.Errorf("You must specify at least a verb and a command."))
	}
	// Split the predicate into its components.
	verb := predicate[0]
	command_line := strings.Join(predicate[1:], " ")
	// Execute the right verb function, or die out.
	switch verb {
	case "while":
		try_while(command_line)
	case "until":
		try_until(command_line)
	default:
		die_on_error(fmt.Errorf("Uknown verb, `%s`.", verb))
	}
}
