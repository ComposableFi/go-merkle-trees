package mmr

type logger interface {
	Errorf(format string, params ...interface{})
	Error(params ...interface{})
}

var log logger

func UseLogger(l logger) {
	log = l
}
