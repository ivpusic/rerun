rerun
============
[![Build Status](https://travis-ci.org/ivpusic/rerun.svg?branch=master)](https://travis-ci.org/ivpusic/rerun)

Recompiling and rerunning go apps when source changes

## Features
- specify list of files/directories to ignore
- specify list of file suffixes to watch (.go, .html, etc.)
- provide application arguments
- configuration using cli-flags and/or json file
- Cross-platform support (Linux, OSX, Windows)

### How to install?
```shell
go get github.com/ivpusic/rerun
```

## Usage
```
usage: rerun [<flags>]

Flags:
  --help               Show context-sensitive help (also try --help-long and --help-man).
  -v, --verbose        Verbose mode. It will show rerun internal messages. Default: false
  -r, --root="."       Project Root.
  -i, --ignore=IGNORE  List of ignored files and directories.
  -a, --args=ARGS      Application arguments.
  -s, --suffixes=SUFFIXES  
                       File suffixes to watch.
  -c, --config=CONFIG  JSON configuration location
  --version            Show application version.
```

To run with default settings just type
```
rerun
```

### Examples

#### CLI flags
```
rerun -r root -a arg1,arg2 -i bower_components,node_modules,test
```

You have troubles? Use verbose mode (``-v`` flag)! You will see a lot of usefull information about rerun internals.
```
rerun -v
```

#### JSON config
Create json file with content, with name for example conf.json
```
{
	"root": ".",
	"ignore": ["some/path/to/ignore1", "some/path/to/ignore2"],
	"args": ["dev", "test"],
	"suffixes": [".go", ".html", ".tpl"]
}
```
and then
```
rerun --conf conf.json
```

#### CLI + JSON
If the same option is provided by cli flag and json config, one from cli will survive.

Example of json config:
```
{
	"ignore": ["some/path/to/ignore"]
}
```
and then
```
rerun --conf conf.json --root server
```

#### ENV variables
You can use environment variables inside your configurations.

##### Linux/OSX
```
{
    "ignore": ["$GOPATH/hello/how/are/you"]
}
```

##### Windows
```
{
    "ignore": ["%GOPATH%/hello/how/are/you"]
}
```

#### Wildcard paths
```
{
	"ignore": ["/some/path", "/some/other/**/*.go"]
}
```

# License
MIT
