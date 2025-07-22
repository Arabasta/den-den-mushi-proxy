package constants

type Key struct {
	Bytes                []byte
	Name                 string
	IsGloballyBlocked    bool
	IsHealthcheckBlocked bool
	IsChangeBlocked      bool
}

var defaultBlocked = [][]byte{
	ArrowUp,
	ArrowDown,
	CtrlR,
	CtrlC,
	CtrlZ,
	CtrlU,
}
