package scheduler

import (
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/AkihiroSuda/go-dag"
)

// Callback is used for Execute
type Callback func(dag.Node) error

// Execute may return one of errors from callback routines.
// Concurrency is set to maximum if set to zero (unimplemented)
func Execute(g *dag.Graph, concurrency uint, cb Callback) error {
	s := &scheduler{
		waitGroups: make(map[dag.Node]*sync.WaitGroup, 0),
	}
	for _, n := range g.Nodes {
		s.init(n)
	}
	var eg errgroup.Group
	for _, n := range g.Nodes {
		n := n
		eg.Go(func() error {
			err := s.routine(n, g.DirectDependees(n), cb)
			return err
		})
	}
	return eg.Wait()
}

// scheduler is not designed for a graph with many nodes
type scheduler struct {
	waitGroups map[dag.Node]*sync.WaitGroup
}

func (s *scheduler) init(n dag.Node) {
	s.waitGroups[n] = &sync.WaitGroup{}
	s.waitGroups[n].Add(1)
}

func (s *scheduler) waitForCompletion(n dag.Node) {
	s.waitGroups[n].Wait()
}

func (s *scheduler) routine(n dag.Node, dependency []dag.Node, cb Callback) error {
	for _, dep := range dependency {
		s.waitForCompletion(dep)
	}
	err := cb(n)
	s.waitGroups[n].Done()
	return err
}
