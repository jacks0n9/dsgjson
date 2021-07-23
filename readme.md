# DsgJson - Convert json files to discord bots
## How to use
```
dsgjson bot.json
```
## How to write your bot.json (or whatever you want to call it)
```json
{
    "token":"your discord bot token",
    // Config for commands
    "commands":{
        "prefix":"!",
        // Actual commands
        "commands":{
           // This will reply to your command with "pong"
           "ping":"pong",
           // Or if you want to send dm's
           "test":{
               // Having a reply or dm field in here are completely optional, but if you don't need to send a dm, just follow the syntax we used in the "ping" command.
               "reply":"replies work",
               "dm":"looks like your dms are open!"
           }
        }
    }
}
```
This is very basic for now, I might add more "events" or things you could listen for, but that's all I have for now.