package compiled

type circularDepFinder struct {
	findFn findFn
}

func newCircularDepFinder(findFn findFn) *circularDepFinder {
	return &circularDepFinder{findFn: findFn}
}

type findFn func(id string) []string

type chainDeps []string

func (c chainDeps) has(id string) bool {
	for _, curr := range c {
		if curr == id {
			return true
		}
	}
	return false
}

func (c circularDepFinder) find(id string) (bool, []string) {
	deps := c.doFind(id, c.findFn(id), make(chainDeps, 0))
	return len(deps) > 0, deps
}

func (c circularDepFinder) doFind(id string, deps []string, chain chainDeps) []string {
	cpChain := append(chain, id)
	if chain.has(id) {
		return cpChain
	}

	for _, d := range deps {
		if currentChain := c.doFind(d, c.findFn(d), cpChain); currentChain != nil {
			return currentChain
		}
	}

	return nil
}
