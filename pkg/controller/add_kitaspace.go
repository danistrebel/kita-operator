package controller

import (
	"github.com/danistrebel/kita-operator/pkg/controller/kitaspace"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, kitaspace.Add)
}
