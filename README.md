# Go-try

Keeps executing a command while or until it returns successfully.


## Usage:

    try while|until commands [command arguments...]

Verbs:

- while: tries until failure
- until: tries until success

Options:

- -h: Shows the help screen


## Examples:

Beep when Google is reachable again:

    try until ping -c1 google.com && beep

Wait until a file, then relaunch a service.

    try while [ -f ./daemon.lock ] && daemon


## Bugs, caveats and roadmap

- Sleep time is hardcoded to one second.
- Currently, commands cannot be piped together.
