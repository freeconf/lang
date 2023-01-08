package main

import "C"

//export hello
func hello(path C.int, yangfile C.int) C.int {
	return C.int(path + yangfile)
	// ypath := source.Path(path)
	// _, err := parser.LoadModule(ypath, yangfile)
	// if err != nil {
	// 	return C.int(1)
	// }

	// return C.int(0)
}
