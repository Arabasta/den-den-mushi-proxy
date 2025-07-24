package constants

var GloballyBlocked = map[string]string{
	string(ArrowUp):   "ArrowUp",
	string(ArrowDown): "ArrowDown",
	string(CtrlR):     "CtrlR",
	string(CtrlZ):     "CtrlZ",
	string(CtrlU):     "CtrlU", // clear line
	string(CtrlA):     "CtrlA", // start of line
	string(CtrlE):     "CtrlE", // end of line
	string(CtrlW):     "CtrlW", // delete word
	string(CtrlD):     "CtrlD", // delete char / EOF
	string(CtrlK):     "CtrlK", // kill to end
	string(CtrlW):     "CtrlW", // delete word
	string(CtrlY):     "CtrlY", // paste from kill buffer

}

func IsGloballyBlockedControlChar(data []byte) bool {
	_, isBlocked := GloballyBlocked[string(data)]
	return isBlocked
}
