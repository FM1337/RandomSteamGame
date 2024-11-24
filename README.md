# RandomSteamGame

A tool that uses the IClientCommService (the same thing the mobile apps use) to install a random steam game and then launch it.

## Why?
I was bored, wasn't sure what I wanted to play, so I decided to build a program to choose for me and to do so as a surprise.

## Usage
Create a `.env` file from the `.env.example` and fill in the required values

Then make sure you have a computer with steam installed and logged in.

Once you've done that, just run the application and it'll pick a random game from your library to install, and then it will install followed by launching it.

## Note
Technically this is an [undocumented API endpoint](https://steamapi.xpaw.me/#IClientCommService) being used and could easily stop working at any time, nor am I 100% certain it's allowed to even be used. Use this at your own risk, I'm not responsible for anything that goes wrong (bans, etc).
