package process

import (
	"math"
	"sort"
)

type iProcess interface {
	pid() int
	parentPid() int
	processName() string
}

type ProcessTree struct {
	Pid         int
	ParentPid   int
	ProcessName string

	Parent    *ProcessTree
	Childrens []*ProcessTree

	X        float64
	Depth    float64
	MaxDepth float64
	mod      float64

	link     *ProcessTree
	ancestor *ProcessTree

	number int
	change float64
	shift  float64

	MinX float64
	MaxX float64
}

func (p *ProcessTree) pid() int            { return p.Pid }
func (p *ProcessTree) parentPid() int      { return p.ParentPid }
func (p *ProcessTree) processName() string { return p.ProcessName }

func (p *ProcessTree) nextRight() (link *ProcessTree) {
	if len(p.Childrens) != 0 {
		return p.Childrens[len(p.Childrens)-1]
	} else if p.link != nil {
		return p.link
	} else {
		return
	}

}

func (p *ProcessTree) nextLeft() (link *ProcessTree) {
	if len(p.Childrens) != 0 {
		return p.Childrens[0]
	} else if p.link != nil {
		return p.link
	} else {
		return
	}
}

func (p *ProcessTree) leftBrother() (b *ProcessTree) {
	if p.Parent != nil {
		for _, n := range p.Parent.Childrens {
			if n == p {
				return b
			} else {
				b = n
			}
		}
	}

	return b
}

func (p *ProcessTree) leftMostSibling() (s *ProcessTree) {
	if p.Parent != nil && p.Parent.Childrens != nil && p.Parent.Childrens[0] != p {
		s = p.Parent.Childrens[0]
		return s
	}

	return s
}

func ListProcesses() (*ProcessTree, error) {
	iProcesses, err := listProcesses()
	if err != nil {
		return nil, err
	}

	processes := []*ProcessTree{}
	// processes := make([]Process, len(iProcesses))
	for _, iProcess := range iProcesses {
		process := &ProcessTree{
			Pid:         iProcess.pid(),
			ParentPid:   iProcess.parentPid(),
			ProcessName: iProcess.processName(),

			Parent:    &ProcessTree{},
			Childrens: []*ProcessTree{},

			number: 1,
			X:      -1,
		}
		process.ancestor = process
		processes = append(processes, process)
	}

	sort.SliceStable(processes, func(i, j int) bool {
		return processes[i].ParentPid < processes[j].ParentPid
	})

	unflat, err := unflattenProcesses(processes)
	if err != nil {
		return nil, err
	}

	buchheim(unflat)

	var min, max *float64
	min = new(float64)
	*min = math.MaxFloat64
	max = new(float64)
	*max = math.MinInt

	getMinMaxX(unflat, min, max)

	unflat.MinX = *min
	unflat.MaxX = *max

	return unflat, err
}

func unflattenProcesses(processes []*ProcessTree) (result *ProcessTree, err error) {
	psMap := make(map[int]*ProcessTree, len(processes))

	for _, p := range processes {
		p1 := p
		psMap[p.Pid] = p1
	}

	for _, p := range psMap {
		if psMap[p.ParentPid] != nil && p.Pid != p.ParentPid {
			psMap[p.ParentPid].Childrens = append(psMap[p.ParentPid].Childrens, p)
			p.Parent = psMap[p.ParentPid]

		}
		// if psMap[p.ParentPid] != nil && len(psMap[p.ParentPid].Childrens) != 0 {
		// 	p.number = len(psMap[p.ParentPid].Childrens)
		// }
	}

	results := []*ProcessTree{}

	rootsnumber := 1
	for _, p := range psMap {
		if psMap[p.ParentPid] != nil && p.Pid != p.ParentPid {
			continue
		}

		// p.number = rootsnumber
		rootsnumber++
		results = append(results, p)
	}

	result = &ProcessTree{
		ProcessName: "fakeRoot",

		Parent:    nil,
		Childrens: results,

		number: 1,
		X:      -1,
	}

	result.ancestor = result

	sortTree(result)

	assignParents(result)

	return
}

func sortTree(t *ProcessTree) {
	sort.SliceStable(t.Childrens, func(i, j int) bool {
		return t.Childrens[i].Pid < t.Childrens[j].Pid
	})

	number := 1
	for _, c := range t.Childrens {
		sortTree(c)
		c.number = number
		number++
	}
}

func getMinMaxX(t *ProcessTree, min, max *float64) {
	for _, c := range t.Childrens {
		if c.X > *max {
			*max = c.X
		}
		if c.X < *min {
			*min = c.X
		}
		getMinMaxX(c, min, max)
	}
}

func assignParents(t *ProcessTree) {
	for _, c := range t.Childrens {
		c.Parent = t
		assignParents(c)
	}
}
