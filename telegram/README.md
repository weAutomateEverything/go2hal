
The service provides a mechaism to communicate directly with telegram

# Environment Variables
HAL_API_SERVICES - Audit endpoint to send all messages too. 

# Admin
HAL exposes a number of API's that allow the user to configure how the bot can assist your group.

The admin services require a JWT Token. There are 2 ways to obtain the token

## Telegram Commands

execute command `/token` and a JWT token will be returned that can be used to admin the group functions. The token is 
only valid for the group in which is was created.  

## API



# Auditing
HAL has been designed to allow for all messages to be send to an audit endpoint.

The following are available:
* halMessageClassification - https://github.com/weAutomateEverything/halMessageClassification