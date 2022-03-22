# watch
Implementation of the Unix watch command in go. This is for Windows users who aren't lucky enough to have it built in 
and who don't want to use WSL. Will _NOT_ run on non Windows platforms because you should be using the natively 
provided watch command. 

### Building

```go build``` 

### Command Line Options

-e    Exit on non-zero return of command

-g    Exit when output changes

-n #  Interval in seconds (default 5)

-p    Try to run at precise interval

-t    Hide title bar

-x    Run with command processor (cmd.exe/pwsh.exe/powershell.exe)


### Exit codes
1 - exit on error selected and watched command error'ed

2 - exit on change selected and watched command's output changed

3 - didn't include a command to run

4 - Couldn't find command processor to run under

5 - Not running on Windows