package app

import "fmt"

type State string

const (
	StateInit             State = "INIT"
	StateSelectDevice     State = "SELECT_DEVICE"
	StateCaptureScreen    State = "CAPTURE_SCREEN"
	StateCheckFudaiIcon   State = "CHECK_FUDAI_ICON"
	StateOpenFudaiDetail  State = "OPEN_FUDAI_DETAIL"
	StateParsePopupLayout State = "PARSE_POPUP_LAYOUT"
	StateOCRContent       State = "OCR_CONTENT"
	StateDecideJoin       State = "DECIDE_JOIN"
	StateClickJoin        State = "CLICK_JOIN"
	StateWaitResult       State = "WAIT_RESULT"
	StateCheckResult      State = "CHECK_RESULT"
	StateClaimReward      State = "CLAIM_REWARD"
	StateSwitchRoom       State = "SWITCH_ROOM"
	StateRecoverState     State = "RECOVER_STATE"
	StateHandleVerify     State = "HANDLE_VERIFICATION"
)

type StateMachine struct {
	current State
}

func NewStateMachine() *StateMachine {
	return &StateMachine{current: StateInit}
}

func (m *StateMachine) Current() State {
	return m.current
}

func (m *StateMachine) Transition(next State) {
	m.current = next
}

func (m *StateMachine) String() string {
	return fmt.Sprintf("StateMachine(current=%s)", m.current)
}
