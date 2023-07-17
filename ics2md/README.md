## ics2md

golang sample for converting ics (iCal format) to markdown.

```bash
# you will get converted output as stdout
$ go run ics2md.go test.ics
## 2011/10/29

test1
test1
test1

## 2017/10/02

### 2017/10/02 22:30:00

test2
test2


### 2017/10/02 23:30:00

test3test3test3test3test3test3test3test3test3test3
test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3test3


# if you want as markdown file, use redirect
$ go run ics2md.go test.ics > ~/Desktop/test.md
```

<img width="773" alt="image" src="https://user-images.githubusercontent.com/1454332/222964218-586aa10e-8e14-4986-b6ee-da7f578d1e56.png">
