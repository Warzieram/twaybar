# TWAYBAR
(or twitchbar if that doesn't get me sued)

## Description 

Twaybar is a basic custom waybar module for twitch integration

## Installation 

- clone this repo then :
```bash
cd twaybar
make build
```
This will create the `twaybar` executable

- Put this in your waybar config:
```
"custom/twaybar": {
    "format": "{}",
    "exec": "path/to/waybar",
    "return-type": "json"
}
```

Note: for now the module doesn't work when executed through hyprland runtime like on startup or:
```bash
hyprctl exec waybar
```

So to make this work you will have to run:
```bash
killall waybar
waybar &
```

Until I figure out a way to make it work (help is welcome)

## Usage

The module uses twitch API so will need a configuration step on twitch

### Step 1: Create a twitch application

Go to https://dev.twitch.tv/console and create an application there   
You can name it whatever you want, you just need to set the OAuth redirection URL to `localhost:[port_of_your_choice]/callback`   
Then you will need to create the `.env` file and paste the content of `.env.example` into it.   

### Step 2: Set the environment 

Here are the fields to fill in `.env`:

`CLIENT_ID`: The client id of the app you created in the twitch dev console   
`CLIENT_SECRET`: The client secret of the app you created in the twitch dev console (available on the app's page)   
`BROADCASTER_LOGIN`: The login of the broadcaster you want info from (usually the username in lowercase)   
`USER_LOGIN`: Your own login, note that you need the rights to access the informations you'll subscribe to, for example
if you're banned from a channel you cannot see its chat messages, and you can only see subs of channel that you are a moderator for.   
`PORT`: The port you want to use for the OAuth server (if you don't know just put 8080)    

### Step 3: Configuration

The `config.json` file is used to chose what events you want to subscribe to :   
`chat` set to `true` will add the last message and its sender username to the text output   
`subs` will add the last sub's username to the tooltip output   
`resubs` will add the last resub's username and its message to  the tooltip output   
`debug` will enable logging off all requests for authorization token and events subscriptions   

### Step 4: Execution

Once all the configuration is done you can launch by restarting waybar :
```bash
killall waybar
waybar &
```

A link should be opened in your default browser to authorize the app to access your twitch account just click "Accept" (if it doesn't open just click it in your terminal)

Then the messages should start being displayed in waybar so Enjoy !

## Troubleshooting

You can set the debug option to `true` in `conf.json` to debug the requests   
The authorization token should refresh automatically when it's expired resulting in the OAuth page opening on startup, but if it doesn't happen you can delete the
`token.json` file and try again.
Do not hesitate to open an issue if something goes wrong

## Contributing

This code is *VERY FAR* from being perfect I did this to learn some go and just have fun so feel free to contribute via pull requests if you see something you
can improve
    

