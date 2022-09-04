# draco

a multiversion gophertunnel proxy to join the latest MC version without renderdragon

# How to use?
1. Download the draco executable off https://github.com/cqdetdev/draco/releases/tag/v0.0.3
2. Run `CheckNetIsolation LoopbackExempt -a -n="Microsoft.MinecraftUWP_8wekyb3d8bbwe"` as an **administrator** in Powershell
3. Run the draco.exe file
4. Follow the link and code the program tells you, this is how the proxy knows that you are you and can authenicate to MS services
5. Once you authenticate, wait for a "Listening" message to pop up
6. Create a server with address `127.0.0.1` and default port
7. Join this server and enjoy!

# Purpose

mojang can't seem to actually make a good update to their game; they decided to add renderdragon to all versions of MC (including x86 builds) in 1.18.30

this update caused the game to run at 1.294 fps, cause 2 hour input delays, and the infamous pink glitch (which mojang has pretended to fix for the past 2 years)

this proxy basically allows the client to join on a proxy that can translate the packets from 1.18.10 <=> 1.19.20 using the gophertunnel multi-protocol api

this allows players to then be able to use the 1.18.10 x86 build (which does NOT have renderdragon) and play on MC without any stutters or input lag

also uses dragonfly chunk code for chunk translation

# Notes

this should work fairly flawlessly but the code is absolutely dogshit as of now and can be significantly improved, with
potential for versions alongside 1.19.20 too. this also doesn't forward packs automatically, at least for now.

expect 20-30ms delay increase as it is a network based proxy

# Current Issues
- command arguments are not handled properly for the latest version
- only works on dragonfly servers

# Further Notes

fuck mojang
