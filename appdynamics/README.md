# Appdynamics

The module allows HAL to communicate with App Dynamics


## REST
Please refer to the swagger endpoint for details
### Alerts
#### Text Alert
HAL has an endpoint available for appdynamics to send a HTTP rest request to HAL.
HAL will then send the alert to telegram

Endpoint is available as a POST to /api/appdynamics/group id

*Sample request*
```json
{
  "environment": "test",
  "events": [
    {
      "severity": "ERROR",
      "application": {
        "name": "TESTAPP"
      },
      "tier": {
        "name": "TESTTIER"
      },
      "node": {
        "name": "test"
      },
      "displayName": "Business Transaction Error",
      "eventMessage": "[Error]  - ServletException: java.lang.NumberFormatException: For input string: \"\""
    }
  ]
}
```

To setup a appdynamics action
1. go to "Alert & Respond"
1. Create a new HTTP Request template
1. Method is POST
1. RAW URL is {hal base url}/api/appdynamics/{group id}
so if your base url is http://hal.interwebz.com and your group id is 45678 then the final url would be
http://hal.interwebz.com/api/appdynamics/45678
1. Under payload, paste the JSON Template below. Change environment as you please.
1. Once done, under actions, create a new HTTP Requets action and link it to the HTTP request template.
1. Use your policies to define the rules you want alerted on.


*JSON Template*
```json
#macro(EntityInfo $item)
	"name": "${item.name}"
#end	
	
	{
	"environment": "DEV / SIT",
	"policy" : {
		"digestDurationInMins": "${policy.digestDurationInMins}",
		"name": "${policy.name}"

	},
	"action": {
		"triggerTime":  "${action.triggerTime}",
		"name": "${action.name}"
		
	},
	"events": [
		#foreach(${event} in ${fullEventList})
		{
			"severity": "${event.severity}",
			"application": {
				#EntityInfo($event.application)
			},
			"tier": {
				#EntityInfo($event.tier)
			},
			"node": {
				#EntityInfo($event.node)
			},

			"displayName": "$event.displayName",
			"eventMessage": "$event.eventMessage"
		}
		#if($foreach.count != $fullEventList.size()) , #end
		#end
	]
}
```

## Config
To get HAL to query Appdynamics, you will need to provide HAL with the information it needs to query Appdynamics

### Endpoint
HAL needs to know the URL, user and password to use to query appdynamics

so, if your user is A-user, your group is customer1 and your password is secret, with app dynamics available on http://appd.yourcomany.com:8090
then the Endpoint address would be

http://A-user%40customer1:secret@appd.yourcomany.com:8090

endpoint is a POST request to /api/appdynamics/"group id"/endpoint

```json

```



# Custom Metrics
## IBM MQ Client

