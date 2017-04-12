package dag

import (
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
)

func TestConnectedComponentRoots(t *testing.T) {
	assert.DeepEqual(t,
		&Graph{
			Nodes: []Node{0, 1, 2, 3},
			Edges: []Edge{
				{Depender: 2, Dependee: 0},
				{Depender: 3, Dependee: 1},
			},
		}.ConnectedComponentRoots(), []Node{0, 1})
	assert.DeepEqual(t,
		&Graph{
			Nodes: []Node{0, 1, 2, 3, 4, 5, 6},
			Edges: []Edge{
				{Depender: 2, Dependee: 0},
				{Depender: 3, Dependee: 1},
				{Depender: 4, Dependee: 2},
				{Depender: 5, Dependee: 2},
				{Depender: 6, Dependee: 5},
			},
		}.ConnectedComponentRoots(), []Node{0, 1})
}

func TestDirectDependers(t *testing.T) {
	assert.DeepEqual(t,
		&Graph{
			Nodes: []Node{0, 1, 2, 3, 4, 5, 6},
			Edges: []Edge{
				{Depender: 2, Dependee: 0},
				{Depender: 3, Dependee: 1},
				{Depender: 4, Dependee: 2},
				{Depender: 5, Dependee: 2},
				{Depender: 6, Dependee: 5},
			},
		}.DirectDependers(2), []Node{4, 5})
}

func TestDirectDependees(t *testing.T) {
	assert.DeepEqual(t,
		&Graph{
			Nodes: []Node{0, 1, 2, 3, 4, 5, 6},
			Edges: []Edge{
				{Depender: 2, Dependee: 0},
				{Depender: 3, Dependee: 1},
				{Depender: 4, Dependee: 2},
				{Depender: 5, Dependee: 2},
				{Depender: 6, Dependee: 5},
			},
		}.DirectDependees(2), []Node{0})
}
