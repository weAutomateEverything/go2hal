package commands

type SetGroup struct {

}

func (s *SetGroup) CommandIdentifier() string {
	return "SetGroup"
}

func (s *SetGroup) CommandDescription() string {
	return "Set Alert Group"
}

func (s *SetGroup) execute(arguments string){

}

func inti(){
	Register(func() Command {
		return &SetGroup{}
	})
}
