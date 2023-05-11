package main

import (
	"envoyer/command"
	"envoyer/config"
	"envoyer/dic"
)

func main() {
	config.LoadConfig()
	dic.InitContainer()
	command.Execute()
}
