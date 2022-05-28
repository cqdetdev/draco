# draco

a multiversion gophertunnel proxy to join the latest MC version without renderdragon

# Purpose

mojang can't seem to actually make a good update to their game; they decided to add renderdragon to all versions of MC (including x86 builds) in 1.18.30

this update caused the game to run at 1.294 fps, cause 2 hour input delays, and the infamous pink glitch (which mojang has pretended to fix for the past 2 years)

this proxy basically allows the client to join on a proxy that can translate the packets from 1.18.10 <=> 1.18.30 using the gophertunnel multi-protocol api

this allows players to then be able to use the 1.18.10 x86 build (which does NOT have renderdragon) and play on MC without any stutters or input lag

# Notes

this project is still in development but does NOW SUPPORT a successful login into a server from 1.18.12 (as of 5/28/2022). However, note that it is still not practical to use this proxy since a) there are alot of unhandled packets that still need to be worked on and b) chunks are not handled properly meaning you will spawn into an empty void.

*Could this be used for client-side multiversion on other versions like 1.16.100 or 1.12.0?* 

# Further Notes

fuck mojang