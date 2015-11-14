## Err
This is what I use for error handling in Go. It's similar to many other libraries and there are probably better ones out there. Never the less, I thought I would share (also, pushing to github makes it even easier for me to import to my own projects!).

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
  if e!=nil {
    panic(e)
  }
  useData(data)
```

### Warn and Log
Warn and Log are very similar. The difference is how the developer intends to handle the error in production. During development, Warn and Log will behave identically, they take a possible error and return OK (true if nil, false if error). If err.Out is not nil, they will write to err.Out.

The difference is what they do when err.Out is nil. In this case, Warn will panic and Log will do nothing. Use Warn for errors that will not be handled at production and Log for errors that will be handled at production.

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
Todo and Depricated are wrappers for Issue that add "Todo: " or "Depricated: " and serve to make using documentation errors a little cleaner.

```go
  func oldFoo(){
    err.Depricated("oldFoo; use newFoo")
  }

  func newFoo(){
    err.Todo("finish newFoo")
  }
```