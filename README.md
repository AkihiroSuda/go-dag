# go-dag

[![Build Status](https://travis-ci.org/AkihiroSuda/go-dag.svg)](https://travis-ci.org/AkihiroSuda/go-dag)
[![GoDoc](https://godoc.org/github.com/AkihiroSuda/go-dag?status.svg)](https://godoc.org/github.com/AkihiroSuda/go-dag)

Minimalistic DAG utility with concurrent scheduler

## Example

See [scheduler/scheduler_test.go](scheduler/scheduler_test.go).

```
  0      1
  |      |
  2      3
  |
 / \
4   5
```

```go
import (
    "github.com/AkihiroSuda/go-dag"
    "github.com/AkihiroSuda/go-dag/scheduler"
)

g := &dag.Graph{
	Nodes: []dag.Node{0, 1, 2, 3, 4, 5},
	Edges: []dag.Edge{
		{Depender: 2, Dependee: 0},
		{Depender: 3, Dependee: 1},
		{Depender: 4, Dependee: 2},
		{Depender: 5, Dependee: 2},
	},
}
concurrency := 0
scheduler.Execute(g, concurrency, func(n dag.Node) { appendTo(got, n) })
assert(t, indexOf(got, 0) < indexOf(got, 2))
assert(t, indexOf(got, 2) < indexOf(got, 4))
assert(t, indexOf(got, 2) < indexOf(got, 5))
assert(t, indexOf(got, 1) < indexOf(got, 3))
```
