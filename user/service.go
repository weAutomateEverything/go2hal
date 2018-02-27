package user

import (
	"bufio"
	"strings"
)

type Service interface {
	parseInputRequest(in string) error
}

type service struct {
	store Store
}

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
		jiraId := tokens[2]

		s.store.AddUpdateUser(employeeNumber, name, jiraId)
	}
	return nil
}
