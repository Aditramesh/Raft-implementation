package state

const (
	follower  = iota
	candidate = 1
	leader    = 2
)

type StateMachine interface {
	becomeLeader()
	becomeCandidate()
	becomeFollower()
	AppendEntiresRPC(appendArguemnts *AppendArguments, appendReturn *AppendReturn) error
	RequestVotesRPC(args *RequestVoteArguments, ret *RequestVoteReturn) error
}
