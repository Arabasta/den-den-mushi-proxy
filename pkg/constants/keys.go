package constants

var (
	Enter     = []byte{13}
	Backspace = []byte{127}

	ArrowUp    = []byte{27, 91, 65}
	ArrowDown  = []byte{27, 91, 66}
	ArrowRight = []byte{27, 91, 67}
	ArrowLeft  = []byte{27, 91, 68}

	CtrlR = []byte{18}
	CtrlC = []byte{3}
	CtrlZ = []byte{26}
	CtrlU = []byte{21}

	PasteStart = []byte{27, 91, 50, 48, 48, 126}
	PasteEnd   = []byte{27, 91, 50, 48, 49, 126}

	SemiColon = []byte{59}
	Ampersand = []byte{38}
	Pipe      = []byte{124}

	// todo: handle >>, 2>, &>, >&2, <<<, <<, 2>&1, etc
	OutputRedirection = []byte{62} // >
	InputRedirection  = []byte{60} // <

	SingleQuote = []byte{39} // '
	DoubleQuote = []byte{34} // "
	Backtick    = []byte{96} // `

	Comma           = []byte{44} // ,
	Colon           = []byte{58} // :
	ExclamationMark = []byte{33} // !

	LeftParenthesis  = []byte{40}  // (
	RightParenthesis = []byte{41}  // )
	LeftBracket      = []byte{91}  // [
	RightBracket     = []byte{93}  // ]
	LeftBrace        = []byte{123} // {
	RightBrace       = []byte{125} // }

	DollarSign     = []byte{36} // $
	EqualSign      = []byte{61} // =
	RightBackslash = []byte{92} // \
	LeftBackslash  = []byte{47} // /
)
