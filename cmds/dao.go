package cmds

type DAOCommand struct {
	CreateDAO CreateDAOCommand `cmd:"" name:"create-dao" help:"create dao to contract account"`
	Propose   ProposeCommand   `cmd:"" name:"propose" help:"propose new proposal"`
	Register  RegisterCommand  `cmd:"" name:"register" help:"register to vote"`
	PreSnap   PreSnapCommand   `cmd:"" name:"snap" help:"snap voting powers"`
}
