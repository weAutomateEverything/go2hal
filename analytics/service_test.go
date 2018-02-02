package analytics

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/zamedic/go2hal/alert"
	"github.com/zamedic/go2hal/chef"
)

func TestService_SendAnalyticsAlert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlert := alert.NewMockService(ctrl)
	mockStore := chef.NewMockStore(ctrl)

	service := NewService(mockAlert, mockStore)

	msg := "{\"text\": \"Chef Converge success - legion-minion-s3\", \"attachments\": [{\"color\": \"#58B957\", \"fields\": [{\"short\": true, \"value\": \"pchfsvr1v.standardbank.co.za\", \"title\": \"Chef Server\"}, {\"short\": true, \"value\": \"chopchop\", \"title\": \"Organisation\"}, {\"short\": true, \"value\": \"Updated: 1, Total: 330\", \"title\": \"Resources\"}, {\"value\": \"[recipe[puff-base-server], recipe[legion::worker], recipe[trevor]]\", \"title\": \"Run List\"}], \"title\": \"legion-minion-s3\"}]}"

	recipes := []chef.Recipe{chef.Recipe{Recipe: "puff-base-server"}}
	mockStore.EXPECT().GetRecipes().Return(recipes, nil)

	expected := "*analytics Event*\nChef Converge success - legion-minion-s3\n*Chef Server* pchfsvr1v.standardbank.co.za\n*Organisation* chopchop\n*Resources* Updated: 1, Total: 330\n*Run List* [recipe[puff-base-server], recipe[legion::worker], recipe[trevor]]\n"
	mockAlert.EXPECT().SendAlert(expected)
	service.SendAnalyticsAlert(msg)

}
