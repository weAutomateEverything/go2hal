package user

import (
	"bufio"
	"strings"
)

/*
Service interface to mamage users on HAL
*/
type Service interface {
	parseInputRequest(in string) error
}

type service struct {
	store Store
}

/*
NewService returns a new User Service
*/
func NewService(store Store) Service {
	return &service{store: store}
}

func (s *service) parseInputRequest(in string) error {
	scanner := bufio.NewScanner(strings.NewReader(in))
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ";")
		employeeNumber := tokens[0]
		name := tokens[1]
		jiraID := tokens[2]

		s.store.AddUpdateUser(employeeNumber, name, jiraID)
	}
	return nil
}
