package rsp

// replyError is the target's error packet (E NN) as an error: the code
// is the hex pair the stub sent.
type replyError struct {
	code string
}

func newReplyError(code string) replyError {
	return replyError{code: code}
}

func (e replyError) Error() string {
	return "assembly/rsp: target error " + e.code
}
