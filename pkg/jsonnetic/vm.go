package jsonnetic

import (
	"github.com/google/go-jsonnet"
	"github.com/neticdk/go-jsonnetic/pkg/jsonnetic/native"
)

func MakeVM(jPath []string, maxStack int) *jsonnet.VM {
	vm := jsonnet.MakeVM()

	vm.Importer(&jsonnet.FileImporter{
		JPaths: jPath,
	})

	for _, nf := range native.Funcs() {
		vm.NativeFunction(nf)
	}

	// Allows for increasing the max stack size if needed from the default of 500
	if maxStack > 0 {
		vm.MaxStack = maxStack
	}

	return vm
}
