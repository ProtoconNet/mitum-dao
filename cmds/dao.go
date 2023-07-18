package cmds

type DAOCommand struct {
	CreateDAO     CreateDAOCommand     `cmd:"" name:"create-dao" help:"create dao to contract account"`
	Propose       ProposeCommand       `cmd:"" name:"propose" help:"propose new proposal"`
	CancelPropose CancelProposeCommand `cmd:"" name:"cancel-propose" help:"cancel propose"`
	Register      RegisterCommand      `cmd:"" name:"register" help:"register to vote"`
	PreSnap       PreSnapCommand       `cmd:"" name:"pre-snap" help:"snap voting powers"`
	Vote          VoteCommand          `cmd:"" name:"vote" help:"vote to proposal"`
	PostSnap      PostSnapCommand      `cmd:"" name:"post-snap" help:"snap voting powers"`
	Execute       ExecuteCommand       `cmd:"" name:"execute" help:"execute proposal"`
}
