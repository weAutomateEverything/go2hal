package skynet



type rebuildNode struct {
}


/* Rebuild Node */
func (s *rebuildNode) commandIdentifier() string {
	return "RebuildNode"
}

func (s *rebuildNode) commandDescription() string {
	return "Rebuilds a node"
}

func (s *rebuildNode) execute(update tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
			SendError(errors.New(fmt.Sprint(err)))
			SendError(errors.New(string(debug.Stack())))

		}
	}()
	RecreateNode(update.Message.CommandArguments(), update.Message.From.UserName)
}


