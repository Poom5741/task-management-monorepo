package dependency

import "errors"

var (
	ErrCircularDependency = errors.New("circular dependency detected")
	ErrSelfDependency     = errors.New("task cannot depend on itself")
	ErrDependencyExists   = errors.New("dependency already exists")
)

type DependencyGraph struct {
	adjacencyList map[string][]string
}

func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		adjacencyList: make(map[string][]string),
	}
}

func (g *DependencyGraph) AddDependency(taskID, dependsOnID string) error {
	if taskID == dependsOnID {
		return ErrSelfDependency
	}

	for _, dep := range g.adjacencyList[taskID] {
		if dep == dependsOnID {
			return ErrDependencyExists
		}
	}

	if g.wouldCreateCycle(taskID, dependsOnID) {
		return ErrCircularDependency
	}

	g.adjacencyList[taskID] = append(g.adjacencyList[taskID], dependsOnID)
	return nil
}

func (g *DependencyGraph) wouldCreateCycle(taskID, dependsOnID string) bool {
	visited := make(map[string]bool)
	return g.hasPath(dependsOnID, taskID, visited)
}

func (g *DependencyGraph) hasPath(from, to string, visited map[string]bool) bool {
	if from == to {
		return true
	}

	if visited[from] {
		return false
	}

	visited[from] = true

	for _, neighbor := range g.adjacencyList[from] {
		if g.hasPath(neighbor, to, visited) {
			return true
		}
	}

	return false
}

func (g *DependencyGraph) GetDependencies(taskID string) []string {
	return g.adjacencyList[taskID]
}
