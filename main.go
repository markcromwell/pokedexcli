package main

import (
	"bufio"
	"fmt"
	"os"
)

func commandExit(commands map[string]cliCommand) error {
	fmt.Println("\nClosing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(commands map[string]cliCommand) error {
	fmt.Println(
		`Welcome to the Pokedex!
Usage:

Commands:`)
	for _, cmd := range commands {
		fmt.Printf(" - %s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(map[string]cliCommand) error
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for fmt.Print("Pokedex > "); scanner.Scan(); fmt.Print("Pokedex > ") {
		command := scanner.Text()

		input := cleanInput(command)

		if len(input) == 0 {
			continue
		}

		cmd, exists := commands[input[0]]
		if exists {
			err := cmd.callback(commands)
			if err != nil {
				fmt.Printf("Error executing command '%s': %v\n", input[0], err)
				continue
			}

		}
	}

}
