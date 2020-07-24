
[Copy from here](https://talks.golang.org/2014/names.slide#1)

# Bad
```
func RuneCount(buffer []byte) int {
    runeCount := 0
    for index := 0; index < len(buffer); {
        if buffer[index] < RuneSelf {
            index++
        } else {
            _, size := DecodeRune(buffer[index:])
            index += size
        }
        runeCount++
    }
    return runeCount
}
```

# Good
```
func RuneCount(b []byte) int {
    count := 0
    for i := 0; i < len(b); {
        if b[i] < RuneSelf {
            i++
        } else {
            _, n := DecodeRune(b[i:])
            i += n
        }
        count++
    }
    return count
}
```

# My View

## Short body

Function should be short, so the name of variable can be easy to rember when you scan the body of function.

Otherwise, try to refactor your code. Use sub function. 

## When use long name

In the above example, count is a longer name compared to i, b, n.

because
```
The greater the distance between a name's declaration and its uses,
the longer the name should be.
```

*count* is returned at the end of the function with the declaration from the start.

But sometimes, if names are a little longer, it makes clearer.

For example, in io package
```
type Reader interface {
  Read(b []byte) (n int, err error)
}

type Writer interface {
  Write(b []byte) (n int, err error)
}
```

I feel it is clearer if the declarations are like this
```
type Reader interface {
  Read(to []byte) (n int, err error)
}

type Writer interface {
  Write(from []byte) (n int, err error)
}
```

