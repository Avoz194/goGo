package entities

type Status int

const (
	active Status = iota
	done
	unknown = -1
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

//Change default for Status
func CreateStatus(sString string) Status{
	switch  sString {
	case "active": return active
	case "done": return done
	}
	return -1;
}
func (s Status) isDone() bool {
	return s == done
}