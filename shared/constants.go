package shared

import "syscall"

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

const (
	enableProcessedInput = 0x1
	enableLineInput      = 0x2
	enableEchoInput      = 0x4
	enableWindowInput    = 0x8
	enableMouseInput     = 0x10
	enableInsertMode     = 0x20
	enableQuickEditMode  = 0x40
	enableExtendedFlag   = 0x80

	enableProcessedOutput = 1
	enableWrapAtEolOutput = 2

	keyEvent              = 0x1
	mouseEvent            = 0x2
	windowBufferSizeEvent = 0x4
)

const (
	rightAltPressed  = 1
	leftAltPressed   = 2
	rightCtrlPressed = 4
	leftCtrlPressed  = 8
	shiftPressed     = 0x0010
	ctrlPressed      = rightCtrlPressed | leftCtrlPressed
	altPressed       = rightAltPressed | leftAltPressed
)

var (
	procAllocConsole                = kernel32.NewProc("AllocConsole")
	procSetStdHandle                = kernel32.NewProc("SetStdHandle")
	procGetStdHandle                = kernel32.NewProc("GetStdHandle")
	procSetConsoleScreenBufferSize  = kernel32.NewProc("SetConsoleScreenBufferSize")
	procCreateConsoleScreenBuffer   = kernel32.NewProc("CreateConsoleScreenBuffer")
	procGetConsoleScreenBufferInfo  = kernel32.NewProc("GetConsoleScreenBufferInfo")
	procWriteConsoleOutputCharacter = kernel32.NewProc("WriteConsoleOutputCharacterW")
	procWriteConsoleOutputAttribute = kernel32.NewProc("WriteConsoleOutputAttribute")
	procGetConsoleCursorInfo        = kernel32.NewProc("GetConsoleCursorInfo")
	procSetConsoleCursorInfo        = kernel32.NewProc("SetConsoleCursorInfo")
	procSetConsoleCursorPosition    = kernel32.NewProc("SetConsoleCursorPosition")
	procReadConsoleInput            = kernel32.NewProc("ReadConsoleInputW")
	procGetConsoleMode              = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode              = kernel32.NewProc("SetConsoleMode")
	procFillConsoleOutputCharacter  = kernel32.NewProc("FillConsoleOutputCharacterW")
	procFillConsoleOutputAttribute  = kernel32.NewProc("FillConsoleOutputAttribute")
	procScrollConsoleScreenBuffer   = kernel32.NewProc("ScrollConsoleScreenBufferW")
	var cutNewlineReplacer = strings.NewReplacer("\r", "", "\n", "")
)

type wchar uint16
type short int16
type dword uint32
type word uint16

type coord struct {
	x short
	y short
}

type smallRect struct {
	left   short
	top    short
	right  short
	bottom short
}

type consoleScreenBufferInfo struct {
	size              coord
	cursorPosition    coord
	attributes        word
	window            smallRect
	maximumWindowSize coord
}
