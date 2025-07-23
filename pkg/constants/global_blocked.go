package constants

var GloballyBlocked = map[string]string{
	string(ArrowUp):   "ArrowUp",
	string(ArrowDown): "ArrowDown",
	string(CtrlR):     "CtrlR",
	string(CtrlZ):     "CtrlZ",
	string(CtrlU):     "CtrlU",
}

func IsGloballyBlockedControlChar(data []byte) bool {
	_, isBlocked := GloballyBlocked[string(data)]
	return isBlocked
}
