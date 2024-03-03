package frontend

var (
	// ClearLine clears the current line.
	ClearLine = "\r\x1b[2K"
	// ClearTerminal clears the terminal.
	ClearTerminal = "\033[2J"
	// CursorUp moves the cursor up.
	CursorUp = "\x1b[1A"
)
