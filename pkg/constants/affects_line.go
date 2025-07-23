package constants

var ControlCharAffectsLine = map[string]struct{}{
	string(Backspace):  {},
	string(ArrowLeft):  {},
	string(ArrowRight): {},
	string(CtrlC):      {},
}

func IsControlCharAffectsLine(b []byte) bool {
	_, affects := ControlCharAffectsLine[string(b)]
	return affects
}
