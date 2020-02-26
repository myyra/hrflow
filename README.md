# hrflow

[HR Flow](https://hrflow.accountor.fi/) CLI

## Usage

`hrflow help`

### Creating Reports

```
NAME:
   hrflow report - add a new hour report

USAGE:
   hrflow report [command options] [arguments...]

OPTIONS:
   --duration value, -d value     number of hours to report, formatted as 8h30m. Will be ignored if both start and end time are defined. (default: "8h")
   --start TIME, -s TIME          Set workday start to TIME. (default: now - HOURS)
   --end TIME, -e TIME            Set workday end to TIME. (default: now)
   --project PROJECT, -p PROJECT  which PROJECT to assign to the report. (default: none)
   --comment COMMENT, -c COMMENT  assign a COMMENT to the report. (default: empty)
   --date DATE                    DATE for the report, format 'd.M.' (years not supported) (default: today)
   --help, -h                     show help (default: false)
```

#### Quickly Reporting 8 Hours

Running `hrflow report` will report an 8 hour workday ending at current time.

## Installation

Create a file for the login at `~/.hrflow` with the contents:

```
username: USERNAME
password: PASSWORD
```

### Homebrew (macOS and Linux)

```
brew tap myyra/hrflow https://github.com/myyra/hrflow
brew install hrflow
```

### Go

Make sure Go is installed on your machine and `$GOPATH/bin` is in you `$PATH`. Then run

```
go get -u github.com/myyra/hrflow
```

### Binary

Binaries are available from the [Releases](https://github.com/myyra/hrflow/releases) page.
