package person

type Task struct{} // just for no errors

func (t *Task) isDone() bool {
	return true
}

// func init() {

// }

type Person struct {
	Id, Name, Email string
	Tasks           []Task
}

func (p *Person) addTask(task Task) {
	p.Tasks = append(p.Tasks, task)
}

func (p *Person) tasksList() []Task {
	return p.Tasks
}

func (p *Person) isAllDone() bool {
	for _, task := range p.Tasks {
		if task.isDone() == false {
			return false
		}
	}
	return true
}

func (p *Person) allActiveTasks() []Task {
	active := []Task{}

	for _, task := range p.Tasks {
		if task.isDone() == false {
			active = append(active, task)
		}
	}
	return active
}
