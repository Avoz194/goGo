package entities

type Status int

const (
	active Status = iota
	done
)

func (s Status) String() string {
	switch s {
	case active:
		return "active"
	case done:
		return "done"
	}
	return "unknown"
}

func (s Status) isDone() bool {
	return s == done
}
