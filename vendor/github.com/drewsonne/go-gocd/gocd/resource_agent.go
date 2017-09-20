package gocd

// GetLinks returns HAL links for agent
func (a *Agent) GetLinks() *HALLinks {
	return a.Links
}
