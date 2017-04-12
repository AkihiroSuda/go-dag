package scheduler

import (
	"testing"
	"time"

	"github.com/AkihiroSuda/go-dag"
	"github.com/docker/docker/pkg/testutil/assert"
)

const (
	defaultSeed     = 42
	defaultMaxSleep = 50 * time.Millisecond
)

func TestExecute0_0(t *testing.T) {
	testExecute0(t, 0)
}

func TestExecute0_1(t *testing.T) {
	testExecute0(t, 1)
}

func TestExecute0_2(t *testing.T) {
	testExecute0(t, 2)
}

func TestExecute0_5(t *testing.T) {
	testExecute0(t, 5)
}

func TestExecute0_32(t *testing.T) {
	testExecute0(t, 32)
}

func testExecute0(t *testing.T, concurrency uint) {
	/*
			 0      1
			 |      |
			 2      3
		         |
			/ \
		       4   5
	*/
	g := &dag.Graph{
		Nodes: []dag.Node{0, 1, 2, 3, 4, 5},
		Edges: []dag.Edge{
			{Depender: 2, Dependee: 0},
			{Depender: 3, Dependee: 1},
			{Depender: 4, Dependee: 2},
			{Depender: 5, Dependee: 2},
		},
	}

	w := determineFakeWorkload(g, defaultMaxSleep, defaultSeed)
	ideal := max(w[0]+w[2]+w[4], w[0]+w[2]+w[5], w[1]+w[3])
	sequential := sum(w)
	t.Logf("Ideal: %v (%.2f times faster than sequential %v)",
		ideal, float64(sequential)/float64(ideal), sequential)

	begun := time.Now()
	got := testExecute(t, g, concurrency, w)
	took := time.Now().Sub(begun)
	t.Logf("Took: %v (%2.2f%% of ideal)", took, 100*float64(took)/float64(ideal))

	assert.Equal(t, len(got), 6)
	assert.Equal(t, indexOf(got, 0) < indexOf(got, 2), true)
	assert.Equal(t, indexOf(got, 2) < indexOf(got, 4), true)
	assert.Equal(t, indexOf(got, 2) < indexOf(got, 5), true)
	assert.Equal(t, indexOf(got, 1) < indexOf(got, 3), true)
}

func TestExecute1_0(t *testing.T) {
	testExecute1(t, 0)
}

func TestExecute1_1(t *testing.T) {
	testExecute1(t, 1)
}

func TestExecute1_2(t *testing.T) {
	testExecute1(t, 2)
}

func TestExecute1_5(t *testing.T) {
	testExecute1(t, 5)
}

func TestExecute1_32(t *testing.T) {
	testExecute1(t, 32)
}

func testExecute1(t *testing.T, concurrency uint) {
	/*
					 0     1
					 |     |
					 2     3
				         |    / \
					/ \  /   |
				       4   5     |
		                       |         |
		                       |  ______/
		                       | /
		                       6
	*/
	g := &dag.Graph{
		Nodes: []dag.Node{0, 1, 2, 3, 4, 5, 6},
		Edges: []dag.Edge{
			{Depender: 2, Dependee: 0},
			{Depender: 3, Dependee: 1},
			{Depender: 4, Dependee: 2},
			{Depender: 5, Dependee: 2},
			{Depender: 5, Dependee: 3},
			{Depender: 6, Dependee: 4},
			{Depender: 6, Dependee: 3},
		},
	}

	w := determineFakeWorkload(g, defaultMaxSleep, defaultSeed)
	ideal := max(w[0]+w[2]+w[4]+w[6], w[0]+w[2]+w[3]+w[5], w[0]+w[2]+w[5], w[1]+w[3]+w[5], w[1]+w[3]+w[6])
	sequential := sum(w)
	t.Logf("Ideal: %v (%.2f times faster than sequential %v)",
		ideal, float64(sequential)/float64(ideal), sequential)

	begun := time.Now()
	got := testExecute(t, g, concurrency, w)
	took := time.Now().Sub(begun)
	t.Logf("Took: %v (%2.2f%% of ideal)", took, 100*float64(took)/float64(ideal))

	assert.Equal(t, len(got), 7)
	assert.Equal(t, indexOf(got, 0) < indexOf(got, 2), true)
	assert.Equal(t, indexOf(got, 2) < indexOf(got, 4), true)
	assert.Equal(t, indexOf(got, 2) < indexOf(got, 5) && indexOf(got, 3) < indexOf(got, 5), true)
	assert.Equal(t, indexOf(got, 1) < indexOf(got, 3), true)
	assert.Equal(t, indexOf(got, 4) < indexOf(got, 6) && indexOf(got, 3) < indexOf(got, 6), true)
}

func TestExecute2_0(t *testing.T) {
	testExecute2(t, 0)
}

func TestExecute2_1(t *testing.T) {
	testExecute2(t, 1)
}

func TestExecute2_2(t *testing.T) {
	testExecute2(t, 2)
}

func TestExecute2_5(t *testing.T) {
	testExecute2(t, 5)
}

func TestExecute2_32(t *testing.T) {
	testExecute2(t, 32)
}

func testExecute2(t *testing.T, concurrency uint) {
	/*
	       42
	      /  \
	   100    200
	     \   / |
	      101  |
	        \  |
	         102
	*/
	g := &dag.Graph{
		Nodes: []dag.Node{42, 100, 101, 102, 200},
		Edges: []dag.Edge{
			{Depender: 100, Dependee: 42},
			{Depender: 200, Dependee: 42},
			{Depender: 101, Dependee: 100},
			{Depender: 101, Dependee: 200},
			{Depender: 102, Dependee: 101},
			{Depender: 102, Dependee: 101},
		},
	}

	w := determineFakeWorkload(g, defaultMaxSleep, defaultSeed)
	ideal := max(w[42]+w[100]+w[101]+w[102], w[42]+w[200]+w[101]+w[102], w[42]+w[200]+w[102])
	sequential := sum(w)
	t.Logf("Ideal: %v (%.2f times faster than sequential %v)",
		ideal, float64(sequential)/float64(ideal), sequential)

	begun := time.Now()
	got := testExecute(t, g, concurrency, w)
	took := time.Now().Sub(begun)
	t.Logf("Took: %v (%2.2f%% of ideal)", took, 100*float64(took)/float64(ideal))

	assert.Equal(t, len(got), 5)
	assert.Equal(t, indexOf(got, 42) < indexOf(got, 100) && indexOf(got, 42) < indexOf(got, 200), true)
	assert.Equal(t, indexOf(got, 100) < indexOf(got, 101) && indexOf(got, 200) < indexOf(got, 101), true)
	assert.Equal(t, indexOf(got, 101) < indexOf(got, 102) && indexOf(got, 200) < indexOf(got, 102), true)
}
