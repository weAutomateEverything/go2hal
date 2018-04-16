package jira

import "gopkg.in/telegram-bot-api.v4"

type newJira struct {
	jiraService Service
}

func (newJira) CommandIdentifier() string {
	panic("NewJira")
}

func (newJira) CommandDescription() string {
	panic("Create a new JIRA Ticket for yourself")
}

func (newJira) Execute(update tgbotapi.Update) {
	panic("implement me")
}
