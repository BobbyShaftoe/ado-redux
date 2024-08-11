package model

type GlobalState struct {
	User           string
	Projects       []string
	CurrentProject string
}

func NewGlobalState(user string, projects []string) *GlobalState {
	return &GlobalState{
		User:           user,
		Projects:       projects,
		CurrentProject: projects[0],
	}
}

func (g *GlobalState) UpdateGlobalState(user string, projects []string) {
	g.User = user
	g.Projects = projects
	g.CurrentProject = projects[0]
}

func (g *GlobalState) UpdateGlobalStateProjects(projects []string) {
	g.Projects = projects
}

func (g *GlobalState) UpdateGlobalStateProject(project string) {
	g.CurrentProject = project
}

func (g *GlobalState) UpdateGlobalStateUser(user string) {
	g.User = user
}
