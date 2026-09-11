package dap

// newCapabilities is the initialize reply: what this adapter supports
// (the whole truth - undeclared requests never arrive).
func newCapabilities() capabilities {
	return capabilities{
		SupportsConfigurationDoneRequest: true,
		SupportsFunctionBreakpoints:      true,
		SupportsInstructionBreakpoints:   true,
		SupportsEvaluateForHovers:        true,
		SupportsReadMemoryRequest:        true,
		SupportsDisassembleRequest:       true,
		SupportsTerminateRequest:         true,
	}
}

// capabilities is the DAP feature negotiation of the adapter.
type capabilities struct {
	SupportsConfigurationDoneRequest bool `json:"supportsConfigurationDoneRequest"`
	SupportsFunctionBreakpoints      bool `json:"supportsFunctionBreakpoints"`
	SupportsInstructionBreakpoints   bool `json:"supportsInstructionBreakpoints"`
	SupportsEvaluateForHovers        bool `json:"supportsEvaluateForHovers"`
	SupportsReadMemoryRequest        bool `json:"supportsReadMemoryRequest"`
	SupportsDisassembleRequest       bool `json:"supportsDisassembleRequest"`
	SupportsTerminateRequest         bool `json:"supportsTerminateRequest"`
}

// launchArgs are the launch configuration attributes - also the schema
// of the VSCode launch.json contribution: a source (assembled
// in-process) or a binary with its symbol sidecar, the base address,
// the deterministic clock. Command delegates the whole session to a
// program that serves DAP itself (a prog example in its -debug mode):
// the adapter spawns it and relays the conversation.
type launchArgs struct {
	Arch    string `json:"arch"`
	Source  string `json:"source"`
	Bin     string `json:"bin"`
	Sym     string `json:"sym"`
	Base    string `json:"base"`
	Icount  bool   `json:"icount"`
	Command string `json:"command"`
}

// sourceRef names a source file (the path is what the editor opened).
type sourceRef struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

// sourceBreakpoint is one requested line breakpoint.
type sourceBreakpoint struct {
	Line   int `json:"line"`
	Column int `json:"column,omitempty"`
}

// setBreakpointsArgs replaces all breakpoints of one source file.
type setBreakpointsArgs struct {
	Source      sourceRef          `json:"source"`
	Breakpoints []sourceBreakpoint `json:"breakpoints"`
}

// breakpoint is one installed (or rejected) breakpoint: the editor
// renders verified as a solid dot, otherwise grey.
type breakpoint struct {
	Verified bool   `json:"verified"`
	Line     int    `json:"line,omitempty"`
	Message  string `json:"message,omitempty"`
}

// setBreakpointsResult carries the verdict per requested breakpoint.
type setBreakpointsResult struct {
	Breakpoints []breakpoint `json:"breakpoints"`
}

// functionBreakpoint is one requested named breakpoint (a label).
type functionBreakpoint struct {
	Name string `json:"name"`
}

// setFunctionBreakpointsArgs replaces all function breakpoints.
type setFunctionBreakpointsArgs struct {
	Breakpoints []functionBreakpoint `json:"breakpoints"`
}

// instructionBreakpoint is one requested address breakpoint.
type instructionBreakpoint struct {
	InstructionReference string `json:"instructionReference"`
	Offset               int    `json:"offset,omitempty"`
}

// setInstructionBreakpointsArgs replaces all instruction breakpoints.
type setInstructionBreakpointsArgs struct {
	Breakpoints []instructionBreakpoint `json:"breakpoints"`
}

// thread is one debuggee execution context (a bare-metal machine has
// exactly one).
type thread struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// threadsResult is the thread list.
type threadsResult struct {
	Threads []thread `json:"threads"`
}

// stackFrame is the single frame at the pc: the source line when the
// line map covers it, the raw address for the disassembly view.
type stackFrame struct {
	Id                          int        `json:"id"`
	Name                        string     `json:"name"`
	Line                        int        `json:"line"`
	Column                      int        `json:"column"`
	Source                      *sourceRef `json:"source,omitempty"`
	InstructionPointerReference string     `json:"instructionPointerReference,omitempty"`
}

// stackTraceResult is the frames of one thread.
type stackTraceResult struct {
	StackFrames []stackFrame `json:"stackFrames"`
	TotalFrames int          `json:"totalFrames"`
}

// scope is one named variable container.
type scope struct {
	Name               string `json:"name"`
	PresentationHint   string `json:"presentationHint,omitempty"`
	VariablesReference int    `json:"variablesReference"`
	Expensive          bool   `json:"expensive"`
}

// scopesResult is the containers of one frame.
type scopesResult struct {
	Scopes []scope `json:"scopes"`
}

// variable is one named value; memoryReference hands the address to
// the editor's memory viewer.
type variable struct {
	Name               string `json:"name"`
	Value              string `json:"value"`
	VariablesReference int    `json:"variablesReference"`
	MemoryReference    string `json:"memoryReference,omitempty"`
}

// variablesResult is the children of one scope.
type variablesResult struct {
	Variables []variable `json:"variables"`
}

// variablesArgs asks for the children of one scope.
type variablesArgs struct {
	VariablesReference int `json:"variablesReference"`
}

// evaluateArgs is one watch/hover/repl expression (a label or an
// address).
type evaluateArgs struct {
	Expression string `json:"expression"`
	FrameId    int    `json:"frameId,omitempty"`
	Context    string `json:"context,omitempty"`
}

// evaluateResult is the rendered expression value.
type evaluateResult struct {
	Result             string `json:"result"`
	Type               string `json:"type,omitempty"`
	VariablesReference int    `json:"variablesReference"`
}

// continueArgs resumes one thread (all of them here).
type continueArgs struct {
	ThreadId int `json:"threadId"`
}

// continueResult reports what resumed.
type continueResult struct {
	AllThreadsContinued bool `json:"allThreadsContinued"`
}

// disassembleArgs asks for a listing window at an address: Offset is
// in bytes, InstructionOffset in instructions from the reference.
type disassembleArgs struct {
	MemoryReference   string `json:"memoryReference"`
	Offset            int    `json:"offset,omitempty"`
	InstructionOffset int    `json:"instructionOffset,omitempty"`
	InstructionCount  int    `json:"instructionCount"`
}

// disassembledInstruction is one listing line.
type disassembledInstruction struct {
	Address     string `json:"address"`
	Instruction string `json:"instruction"`
}

// disassembleResult is the listing.
type disassembleResult struct {
	Instructions []disassembledInstruction `json:"instructions"`
}

// readMemoryArgs is one memory window request.
type readMemoryArgs struct {
	MemoryReference string `json:"memoryReference"`
	Offset          int    `json:"offset,omitempty"`
	Count           int    `json:"count"`
}

// readMemoryResult is the window: Data rides base64 through the json
// encoding; UnreadableBytes reports the tail a bad address swallowed.
type readMemoryResult struct {
	Address         string `json:"address"`
	Data            []byte `json:"data,omitempty"`
	UnreadableBytes int    `json:"unreadableBytes,omitempty"`
}

// outputBody is the console output event body (the program's own
// printing and the adapter's warnings).
type outputBody struct {
	Category string `json:"category,omitempty"`
	Output   string `json:"output"`
}

// stoppedBody is the stop event body: why, which thread.
type stoppedBody struct {
	Reason            string `json:"reason"`
	ThreadId          int    `json:"threadId"`
	AllThreadsStopped bool   `json:"allThreadsStopped"`
}

// exitedBody is the process exit event body.
type exitedBody struct {
	ExitCode int `json:"exitCode"`
}
