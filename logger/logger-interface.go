package logger

type ILogger interface {
	D(format string, args ...interface{})
	I(format string, args ...interface{})
	N(format string, args ...interface{})
	W(format string, args ...interface{})
	E(format string, args ...interface{})
	C(format string, args ...interface{})
	Dc(c int, format string, args ...interface{})
	Ic(c int, format string, args ...interface{})
	Nc(c int, format string, args ...interface{})
	Wc(c int, format string, args ...interface{})
	Ec(c int, format string, args ...interface{})
	Cc(c int, format string, args ...interface{})
}
