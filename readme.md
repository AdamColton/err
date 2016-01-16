## Err
A library to help with Go error handling and debugging. In particular, helping
with the transition from development to production code.

This library has 3 main sections error handling, documentation errors and
debugging output.

Error has three core functions Log, Panic and Warn. Warn being a special
case that can be switched to behave like either Log or Panic. The idea being
that during development, we optimistically attempt to continue, but when we're
doing release testing or we've released, we can panic on unrecoverable errors.
This also lets us direct non-panic inducing error output to stdErr during
development, but direct that information to either a file or memory cache in
production.

Documentation errors are issued to err.Warn - so they will not cause a panic
during development but they will panic during release testing.

Debug output can be used as needed or left in place. A useful technique is to
leave debugging statements in place, but turn off debugging output (err.DebugOut
= nil). Then, if a program enters a state that requires investigation, the
output can be directed to stdOut during development, or to a file in production.

### Panic
```go
  data, e := SomeFunc()
  err.Panic(e)
  useData(data)
  // same as
  data, e := SomeFunc()
  if e != nil {
    panic(e)
  }
  useData(data)
```

### Warn and Log
Warn and Log are very similar. The difference is how the developer intends to
handle the error in production. During development, Warn and Log will behave
identically, they take a possible error and return OK (true if nil, false if
error). If err.Out is not nil, they will write to err.Out.

The difference is what they do when PanicOnWarn is true. In this case, Warn will
panic. Use Warn for errors that will not be handled at production and Log for
errors that will be handled (or intentionally ignored) at production.

For production, it is recommended to set PanicOnWarn to true and err.Out to a
file. By default, err.Out is set to stderr.

```go
  if data, e := SomeFunc(); err.Warn(e){
    useData(data)
  }

  if data, e := SomeFunc(); err.Log(e){
    useData(data)
  } else {
    handleError(e)
  }

  // same as
  if data, e := SomeFunc(); e == nil{
    useData(data)
  } else {
    fmt.Println(e) //for debug only, remove before release
    // panic(e) <-- this is for production
  }

  if data, e := SomeFunc(); e == nil{
    useData(data)
  } else {
    fmt.Println(e) //for debug only, remove before release
    handleError(e)
  }
```

It is worth noting that err.Out is a simple interface:
```go
type stringWriter interface {
  WriteString(string) (int, error)
}

var Out = stringWriter(os.Stderr)
```
Both *os.File and bufio.Writer implement this interface, and it's easy to match.

### Test
This function is specifically for unit testing.
```go
  e := mundaneCallInTest()
  err.Test(e,t)
  // same as
  e := mundaneCallInTest()
  if e != nil {
    t.Error(e)
  }
```

### Issue, Todo and Deprecated
Issue will issue and error to Warn.
```go
  err.Issue("This call is to slow to release")
  callSlowFunction()
```
Todo and Deprecated are wrappers for Issue that add "Todo: " or "Depricated: "
and serve to make using documentation errors a little cleaner.

```go
  func oldFoo(){
    err.Deprecated("oldFoo; use newFoo")
  }

  func newFoo(){
    err.Todo("finish newFoo")
  }
```

### Debug
The err.Debug function take any input and outputs a string (using fmt.Sprintln).
If err.DebugShowFile is true, it will precede the output with the file and line
number of the caller.

The file and line number can be disabled by setting err.DebugShowFile to false.
By default, the output is sent to StdOut, but it can be redirected to anything
that implements WriteString(string) (int, error), which includes *os.File and
bufio.Writer.

During debugging, it can be useful to toggle DebugEnabled on and off so that you
can see output from sections of code that are under investigation, but suppress
the same output statements when invoked from sections of code that are not of
interest.

```go
  func foo(){
    bar(20)
    err.DebugOut = err.Stdout
    bar(10)
    err.DebugOut = nil
  }

  func bar(x float64) float64{
    s := float64(0)
    for i := float64(0); i<15; i++{
      y := math.Sqrt(x-i)
      err.Debug(x-i,y)
      s += y
    }
    return s
  }
```

### Bash Terminal Coloring
I have my terminal setup to color stderr, which I find extremely helpful. This
is in my .bashrc:
```bash
color()(set -o pipefail;"$@" 2>&1>&3|sed $'s,.*,\e[38;5;218m&\e[m,'>&2)3>&1
function g () {
  if [ -z $1 ]
  then
    clear && color go run main.go
  else
    clear && color go run $1
  fi
}

alias gt="clear && color go test"
```

Anything sent to stdErr will now be pink-ish, while everything sent to stdOut
will be white (or whatever your default terminal coloring is).

### To-Do
Caching: Build an interface to wrap the necessary methods on bufio.ReadWriter so
that we can write information to a cache (from err.Log, err.Debug or the
application) and if we panic, we can try to write the cache to either a file,
stdOut or errOut. We'd also need a cache limit so that only the last X bytes
written to the cache would be stored. That way a long running process (like a
server) won't consume G's of memory for the error cache. Also, we should try to
write the error that caused the panic to this cache.