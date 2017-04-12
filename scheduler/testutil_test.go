package scheduler

import (
	"math/rand"
	"testing"
	"time"

	"github.com/AkihiroSuda/go-dag"
	"github.com/docker/docker/pkg/testutil/assert"
)

func determineFakeWorkload(g *dag.Graph, maxSleep time.Duration, seed int64) map[dag.Node]time.Duration {
	rnd := rand.New(rand.NewSource(seed))
	res := make(map[dag.Node]time.Duration)
	for _, n := range g.Nodes {
		sleep := time.Duration(rnd.Int63n(int64(maxSleep)))
		res[n] = sleep
	}
	return res
}

func sum(workload map[dag.Node]time.Duration) time.Duration {
	x := time.Duration(0)
	for _, d := range workload {
		x += d
	}
	return x
}

func max(cands ...time.Duration) time.Duration {
	x := time.Duration(0)
	for _, cand := range cands {
		if cand > x {
			x = cand
		}
	}
	return x
}

func testExecute(t *testing.T, g *dag.Graph, concurrency uint, workload map[dag.Node]time.Duration) []dag.Node {
	if concurrency > 0 {
		t.Skip("concurrency > 0 unimplemented yet")
	}
	c := make(chan dag.Node, len(g.Nodes))
	err := Execute(g, concurrency, func(n dag.Node) error {
		time.Sleep(workload[n])
		c <- n
		return nil
	})
	assert.NilError(t, err)
	var got []dag.Node
	for i := 0; i < len(g.Nodes); i++ {
		n := <-c
		got = append(got, n)
	}
	t.Logf("got: %v", got)
	return got
}

func indexOf(nodes []dag.Node, node dag.Node) int {
	for i, n := range nodes {
		if n == node {
			return i
		}
	}
	panic("node not found")
}
