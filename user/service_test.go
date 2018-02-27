package user

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestService_parseInput(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockUser := NewMockStore(ctrl)
	service := NewService(mockUser)

	mockUser.EXPECT().AddUpdateUser("emp1", "nam1", "jira1")
	mockUser.EXPECT().AddUpdateUser("emp2", "nam2", "jira2")

	service.parseInputRequest("emp1;nam1;jira1\nemp2;nam2;jira2")

}
