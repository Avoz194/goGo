package entities

type Status int


const (
	done Status = iota
	active
	unknown = -1
)

var UnkownStatus = Status(unknown)
var AllStatuses = []string{"done", "active"}
var AllStatusIDs = []Status{done, active}

func (s Status) String() string {
	return AllStatuses[s]
}

//Change default for Status
func CreateStatus(sString string) Status{
	switch  sString {
	case "active": return active
	case "done": return done
	case "": return active
	}
	return unknown;
}
func (s Status) isDone() bool {
	return s == done
}

