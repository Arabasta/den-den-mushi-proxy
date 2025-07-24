package constants

var ControlChars = map[string]string{
	string(Enter):     "Enter",
	string(Backspace): "Backspace",

	string(ArrowUp):    "ArrowUp",
	string(ArrowDown):  "ArrowDown",
	string(ArrowRight): "ArrowRight",
	string(ArrowLeft):  "ArrowLeft",

	string(CtrlA): "CtrlA",
	string(CtrlB): "CtrlB",
	string(CtrlC): "CtrlC",
	string(CtrlD): "CtrlD",
	string(CtrlE): "CtrlE",
	string(CtrlF): "CtrlF",
	string(CtrlG): "CtrlG",
	string(CtrlH): "CtrlH",
	string(CtrlK): "CtrlK",
	string(CtrlL): "CtrlL",
	string(CtrlM): "Enter",
	string(CtrlN): "CtrlN",
	string(CtrlO): "CtrlO",
	string(CtrlP): "CtrlP",
	string(CtrlQ): "CtrlQ",
	string(CtrlR): "CtrlR",
	string(CtrlS): "CtrlS",
	string(CtrlT): "CtrlT",
	string(CtrlU): "CtrlU",
	string(CtrlV): "CtrlV",
	string(CtrlW): "CtrlW",
	string(CtrlX): "CtrlX",
	string(CtrlY): "CtrlY",
	string(CtrlZ): "CtrlZ",

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
