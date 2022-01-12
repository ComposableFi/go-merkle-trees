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

type logg struct{}

// Errorf has no log body, it simply satisfies the logger interface
func (l logg) Errorf(format string, params ...interface{}) {}

// Error has no log body, it simply satisfies the logger interface
func (l logg) Error(params ...interface{}) {}

func init() {
	// set log to default logger which prints nothing if no logger interface is passed to the UseLogger method
	log = logg{}
}
