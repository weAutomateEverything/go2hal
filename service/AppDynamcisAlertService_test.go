package service

import "testing"

func TestAppDynamicsAlerts(t *testing.T) {
		SendAppdynamicsAlert("		{	\"environment\": \"DEV / SIT\",	\"policy\" : {		\"digestDurationInMins\": \"1\",		\"name\": \"My Policy\"	},	\"action\": {		\"triggerTime\":  \"Wed Sep 27 14:29:19 SAST 2017\",		\"name\": \"Telegram\"			},	\"events\": [				{			\"severity\": \"WARN\",			\"application\": {					\"name\": \"SBIS\"			},			\"tier\": {					\"name\": \"WEB\"			},			\"node\": {					\"name\": \"dsbis02v\"			},			\"displayName\": \"Ongoing Warning Health Rule Violation\",			\"eventMessage\": \"AppDynamics has detected a problem with Node <b>dsbis02v</b>.<br><b>Memory utilization is too high</b> continues to violate with <b>warning</b>.<br>All of the following conditions were found to be violating<br>For Node <b>dsbis02v</b>:<br>1) Hardware Resources|Memory|Used % Condition<br><b>Used %'s</b> value <b>88.0</b> was <b>greater than</b> the threshold <b>75.0</b> for the last <b>30</b> minutes<br>\"		}					]}")
}


