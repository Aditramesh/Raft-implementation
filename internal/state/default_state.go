package state

type logType struct {
	term      int
	data      int
	dataindex int
}

type State struct {
	persistantState *persistantState
	nextIndex       []int
	matchIndex      []int
}

type persistantState struct {
	stateType   int
	currentTerm int
	votedFor    int
	log         []logType
	commitIndex int
	lastApplied int
}

type AppendArguments struct {
	term         int
	leaderID     int
	prevLogIndex int
	prevLogTerm  int
	entries      []logType
	leaderCommit int
}

type AppendReturn struct {
	term    int
	success bool
}

type RequestVoteArguments struct {
	term         int
	candidateId  int
	lastLogIndex int
	lastLogTerm  int
}

type RequestVoteReturn struct {
	term        int
	voteGranted bool
}

func NewDefaultService() *State {
	initLog := make([]logType, 0)
	persistantState := &persistantState{
		stateType:   follower,
		currentTerm: 0,
		votedFor:    -1,
		log:         initLog,
		commitIndex: 0,
		lastApplied: 0,
	}
	return &State{
		persistantState: persistantState,
		nextIndex:       []int{1, 1, 1, 1, 1},
		matchIndex:      []int{0, 0, 0, 0, 0},
	}
}

func (s *State) becomeFollower() {
	s.persistantState.stateType = follower
}

func (s *State) becomeCandidate() {
	s.persistantState.stateType = candidate
}

func (s *State) becomeLeader() {
	s.persistantState.stateType = leader
}

func (s *State) PersistToDisk() {

}

func (s *State) AppendEntiresRPC(appendArguemnts *AppendArguments, appendReturn *AppendReturn) error {
	if appendArguemnts.term < s.persistantState.currentTerm {
		appendReturn.success = false
		appendReturn.term = s.persistantState.currentTerm
		return nil
	}

	if appendArguemnts.term > s.persistantState.currentTerm {
		s.persistantState.currentTerm = appendArguemnts.term
	}

	if appendArguemnts.prevLogIndex >= len(s.persistantState.log) || s.persistantState.log[appendArguemnts.prevLogIndex].term != appendArguemnts.prevLogTerm {
		appendReturn.success = false
		appendReturn.term = s.persistantState.currentTerm
		return nil
	}

	if len(appendArguemnts.entries) > 0 {
		s.persistantState.log = append(s.persistantState.log[:appendArguemnts.prevLogIndex+1], appendArguemnts.entries...)
	}

	if appendArguemnts.leaderCommit > s.persistantState.commitIndex {
		s.persistantState.commitIndex = min(appendArguemnts.leaderCommit, len(s.persistantState.log)-1)
	}
	appendReturn.success = true
	appendReturn.term = s.persistantState.currentTerm
	return nil

}

func (s *State) RequestVotesRPC(args *RequestVoteArguments, ret *RequestVoteReturn) error {
	if args.term < s.persistantState.currentTerm {
		ret.voteGranted = false
		ret.term = s.persistantState.currentTerm
		return nil
	}

	if args.term > s.persistantState.currentTerm {
		s.persistantState.currentTerm = args.term
		s.persistantState.votedFor = -1
		s.becomeFollower()
	}

	if args.term == s.persistantState.currentTerm && args.lastLogIndex >= len(s.persistantState.log)-1 {
		if s.persistantState.votedFor == -1 || s.persistantState.votedFor == args.candidateId {
			s.persistantState.votedFor = args.candidateId
			ret.voteGranted = true
			return nil
		}
		ret.voteGranted = false
		return nil
	}
	ret.term = s.persistantState.currentTerm
	return nil
}
