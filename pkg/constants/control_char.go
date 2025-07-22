package constants

var ControlChars = map[string]string{
	string(Enter):     "Enter",
	string(Backspace): "Backspace",

	string(ArrowUp):    "ArrowUp",
	string(ArrowDown):  "ArrowDown",
	string(ArrowRight): "ArrowRight",
	string(ArrowLeft):  "ArrowLeft",

	string(CtrlR): "CtrlR",
	string(CtrlC): "CtrlC",
	string(CtrlZ): "CtrlZ",
	string(CtrlU): "CtrlU",

	string(SemiColon): ";",
	string(Ampersand): "&",
	string(Pipe):      "|",

	string(OutputRedirection): ">",
	string(InputRedirection):  "<",

	string(SingleQuote): "'",
	string(DoubleQuote): "\"",
	string(Backtick):    "`",

	string(Comma):           ",",
	string(Colon):           ":",
	string(ExclamationMark): "!",

	string(LeftParenthesis):  "(",
	string(RightParenthesis): ")",
	string(LeftBracket):      "[",
	string(RightBracket):     "]",
	string(LeftBrace):        "{",
	string(RightBrace):       "}",

	string(DollarSign):     "$",
	string(EqualSign):      "=",
	string(RightBackslash): "\\",
	string(LeftBackslash):  "/",
}
