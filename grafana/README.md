A webhook for grafana alerts

# API

## REST

POST /api/graphana/{chatid:\[0-9\]+}

```json
{
  "evalMatches": [
    {
      "value": 100,
      "metric": "High value",
      "tags": null
    },
    {
      "value": 200,
      "metric": "Higher Value",
      "tags": null
    }
  ],
  "imageUrl": "http://grafana.org/assets/img/blog/mixed_styles.png",
  "message": "Someone is testing the alert notification within grafana.",
  "ruleId": 0,
  "ruleName": "Test notification",
  "ruleUrl": "http://localhost:3000/",
  "state": "alerting",
  "title": "[Alerting] Test notification"
}


```