## Err
This is what I use for error handling in Go. It's similar to many other
libraries and there are probably better ones out there. Never the less, I
thought I would share (also, pushing to github makes it even easier for me to
import to my own projects!).

I have my terminal setup to color stderr, which helps a lot. This is in my .bashrc:
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
  if data,e := SomeFunc(); err.Warn(e){
    useData(data)
  }

  if data,e := SomeFunc(); err.Log(e){
    useData(data)
  } else {
    handleError(e)
  }

  // same as
  if data,e := SomeFunc(); e == nil{
    useData(data)
  } else {
    fmt.Println(e) //for debug only, remove before release
    // panic(e) <-- this is for production
  }

  if data,e := SomeFunc(); e == nil{
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

### Issue, Todo and Depricated
Issue will issue and error to Warn.
```go
  err.Issue("This call is to slow to release")
  callSlowFunction()
```
Todo and Depricated are wrappers for Issue that add "Todo: " or "Depricated: "
and serve to make using documentation errors a little cleaner.

```go
  func oldFoo(){
    err.Depricated("oldFoo; use newFoo")
  }

  func newFoo(){
    err.Todo("finish newFoo")
  }
```

### Debug
Debug is similar to fmt.Println, actually, it's a wrapper around it. Only if
err.DebugEnabled is true will Debug print the values passed in. It will also
print the filename and line number if err.DebugShowFile is true. By default,
DebugEnabled is false and DebugShowFile is true.

During debugging, it can be useful to toggle DebugEnabled on and off so that you
can see output from sections of code that are under investigation, but supress
the same output statements when invoked from sections of code that are not of
interest.

```go
  func foo(){
    bar(20)
    reset := err.DebugEnabled
    err.DebugEnabled = true
    bar(10)
    err.DebugEnabled = reset
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