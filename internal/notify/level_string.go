package notify

// LevelFrom parses a string into a Level, returning LevelInfo and false
// if the string is not a recognised level.
func LevelFrom(s string) (Level, bool) {
	switch Level(s) {
	case LevelInfo, LevelWarn, LevelError:
		return Level(s), true
	default:
		return LevelInfo, false
	}
}

// String returns the string representation of a Level.
func (l Level) String() string {
	return string(l)
}

// IsValid reports whether the Level is one of the defined constants.
func (l Level) IsValid() bool {
	_, ok := LevelFrom(string(l))
	return ok
}
