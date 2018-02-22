# GO2HAL Alert Module

The module allows for alert messages and images to be send to 3 types of groups

## Non Technical

Use this to send messages and images to a group that consists of non developers but individuals who have an interest in
your system.

So, lets say you want to send an alert that a server is down, for the non technical users, information such as server name,
ip, ect will not be valuable, instead they might just want something that says a node is down, its being investigated.

We also use it to inform users that a critical batch has succeeded or failed, or that a new verion of code has been deployed.

## Technical

Messages that are targetted towards a technical audience, such as "Response time from system X has breached a response time of Y seconds".

## Heartbeat

This group is for the bot administrators, used mainly when the bot encouteres an internal error or panic. Also used when
testing new changes that you dont want to sent to the other groups yet.

#Server

The server connects to telegram to send alerts.

example server can be found in server.go in the examples folder.

There are 3 commands

## SetGroup

Sets the Alert Group to the one where this command was run

## SetNonTechGroup

Sets the No Technical Group to the one where this command was run

## SetHeartbeatGroup

Sets the Heartbeat Group to the one where this command was run

#Client

If you have HAL service running and wish to send alerts from an external application, we have provided a proxy service which
will send Rest requests to send the alerts.

see examples in client.go in the examples folder