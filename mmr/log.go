package mmr

type logger interface {
	Errorf(format string, params ...interface{})
	Error(params ...interface{})
}

var log logger

// UseLogger takes any type that satisfies the logger interface as a argument. It the sets the global log variable to the
// the value passed as an argument.
func UseLogger(l logger) {
	log = l
}
