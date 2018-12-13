# Goro Goso 🧟‍
Minimal Golang process auto reload by glob watcher


# Feature
- Support Http server reload ( I only need this one )


# Usage
```
$ gorogoso -watch="**/*.go" -entry="main.go"
```
### Also support multiple globs
```
$ gorogoso -watch"test/**/*.go,lib/**/*.go"
```


# Build from source
```
$ make build
```


# Todo
- ignore path args
