package gocd

// RemoveLinks gets the EnvironmentsResponse ready to be submitted to the GoCD API.
func (er *EnvironmentsResponse) RemoveLinks() {
	er.Links = nil
	for _, env := range er.Embedded.Environments {
		env.RemoveLinks()
	}
}

// RemoveLinks gets the Environment ready to be submitted to the GoCD API.
func (env *Environment) RemoveLinks() {
	env.Links = nil
	for _, p := range env.Pipelines {
		p.RemoveLinks()
	}
	for _, a := range env.Agents {
		a.RemoveLinks()
	}
}
