# go-cancelable-file-writer

A Go file writer that can be canceled with a `context.Context`. Particularly
useful when writing large files and there might be a need for a timeout or
manual cancel to be able to properly clean up incomplete files.

## Adding this package to your project

```bash
go get github.com/TSonono/go-cancelable-file-writer@latest
```

## Example usage

```go
// data := ...

// mainContext might already also have a prior manual cancel in addition to timeout
ctx, cancel := context.WithTimeout(mainContext, 10*time.Second)
defer cancel()

file, err := os.Create("test.txt")
if err != nil {
    log.Fatal(err.Error())
}
defer file.Close()

_, err = cancelablefw.CancelableFileWriter(ctx, data, file)
if err != nil {
    {
        err := os.Remove("test.txt")
        if err != nil {
            log.Fatal(err.Error())
        }
    }
    log.Fatal(err.Error())
}
```
