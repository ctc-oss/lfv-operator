package controller

import (
	"github.com/jw3/example-operator/pkg/controller/datavolume"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, datavolume.Add)
}
