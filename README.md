# ghorg
Analyse a github organization to find top contributors by number of commits and number of repos

## Usage

### Main
```
NAME:
   ghorg - Analyse a github organization to find top contributors by number of commits and number of repos

USAGE:
   ghorg [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     analyse, a  Analyse a github org

GLOBAL OPTIONS:
   --ghtoken value, -g value  (Mandatory) Your github token [$GH_TOKEN]
   --verbose, --vvv           To use it in verbose mode
   --help, -h                 show help
   --version, -v              print the version
```

### Analyse command

```
NAME:
   ghorg analyse - Analyse a github org

USAGE:
   ghorg analyse [command options] name-of-the-org-to-scan

OPTIONS:
   --repos-details, --rd  To show in table all repos names owned by user
   --csv value            file name to ouput in csv format
   --no-markdown          To hide markdown result
```