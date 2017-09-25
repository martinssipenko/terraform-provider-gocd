package gocd

import (
	"fmt"
)

// Equal is true if the two materials are logically equivalent. Not neccesarily literally equal.
func (m Material) Equal(a *Material) bool {
	if m.Type != a.Type {
		return false
	}
	switch m.Type {
	case "git":
		return m.Attributes.equalGit(&a.Attributes)
	default:
		panic(fmt.Errorf("Material comparison not implemented for '%s'", m.Type))
	}
}

func (a1 MaterialAttributes) equalGit(a2 *MaterialAttributes) bool {
	if a1.URL == a2.URL {
		if a1.Branch == a2.Branch {
			// Check if branches are equal
			return true
		} else if a1.Branch == "" && a2.Branch == "master" || a1.Branch == "master" && a2.Branch == "" {
			// Check if branches are master equal
			return true
		}
	}
	return false
}
