package controller

import (
	"github.com/philipsahli/hugo-operator/pkg/controller/site"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, site.Add)
}
