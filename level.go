package logger

type Level int

const (
	DebugLevel = Level(iota)
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

// LevelMap maps a level to a string.
var LevelMap = map[Level]string{
	DebugLevel: "debug",
	InfoLevel:  "info",
	WarnLevel:  "warn",
	ErrorLevel: "error",
	PanicLevel: "panic",
	FatalLevel: "fatal",
}

// String returns the string representation of the level.
func (l Level) String() string {
	return LevelMap[l]
}
