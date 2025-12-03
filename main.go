package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/markcromwell/pokedexcli/internal/pokeapi"
)

type config struct {
	nextURL        *string
	prevURL        *string
	caughtPokemons map[string]*pokeapi.Pokemon
}

func commandExit(commands map[string]cliCommand, cfg *config, param []string) error {
	fmt.Println("\nClosing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(commands map[string]cliCommand, cfg *config, param []string) error {
	fmt.Println(
		`Welcome to the Pokedex!
Usage:

Commands:`)
	for _, cmd := range commands {
		fmt.Printf(" - %s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(commands map[string]cliCommand, cfg *config, param []string) error {
	if cfg.nextURL == nil || *cfg.nextURL == "" {
		cfg.nextURL = new(string)
		*cfg.nextURL = pokeapi.PokeAPIBaseURLMap
	}

	locations, err := pokeapi.GetLocationAreas(*cfg.nextURL)
	if err != nil {
		return err
	}

	for _, loc := range locations.Results {
		fmt.Printf("%s\n", loc.Name)
	}

	prevURL := cfg.nextURL
	cfg.nextURL = locations.Next
	cfg.prevURL = prevURL

	return nil
}

func commandMapb(commands map[string]cliCommand, cfg *config, param []string) error {
	if cfg.prevURL == nil || *cfg.prevURL == "" {
		fmt.Println(`you're on the first page`)
		return nil
	}

	locations, err := pokeapi.GetLocationAreas(*cfg.prevURL)
	if err != nil {
		return err
	}

	for _, loc := range locations.Results {
		fmt.Printf("%s\n", loc.Name)
	}

	cfg.nextURL = locations.Next
	cfg.prevURL = locations.Previous
	return nil
}

func commandExplore(commands map[string]cliCommand, cfg *config, param []string) error {
	if len(param) == 0 {
		return fmt.Errorf("please specify a location to explore")
	}

	locationArea, err := pokeapi.GetLocationArea(param[0])
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", locationArea.Name)
	fmt.Printf("Found Pokemon:\n")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

// commandCatch expects only 1 param, the pokemon name to catch. it uses the client.go GetPokemon to see if the pokemon exists, returns an error if it doesn't. then pritns "Throwing a Pokeball at %s... where %s is the name of the Pokemon. A percentage chance will be used at catching the pokemon will be calculated based on base experience (higher the harder). If the roll suceeds a map of Pokemon index by name will be used to store the Pokemon. Also should check if the Pokemon is already caught before trying to catch again."
func commandCatch(commands map[string]cliCommand, cfg *config, param []string) error {
	if len(param) == 0 {
		return fmt.Errorf("please specify a Pokemon to catch")
	}

	// check if already caught
	if _, caught := cfg.caughtPokemons[param[0]]; caught {
		fmt.Printf("You have already caught %s.\n", param[0])
		return nil
	}

	pokemonName := param[0]
	pokemon, err := pokeapi.GetPokemon(pokemonName)
	if err != nil {
		return fmt.Errorf("could not find Pokemon '%s': %v", pokemonName, err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	// Simple catch chance calculation based on base experience
	catchChance := 50.0 - (float64(pokemon.BaseExperience) / 10.0)
	if catchChance < 5.0 {
		catchChance = 5.0 // Minimum 5% chance
	}

	roll := float64(rand.Intn(100))
	if roll < catchChance {
		if cfg.caughtPokemons == nil {
			cfg.caughtPokemons = make(map[string]*pokeapi.Pokemon)
		}

		cfg.caughtPokemons[pokemon.Name] = pokemon
		fmt.Printf("Congratulations! You caught %s!\n", pokemon.Name)
		fmt.Printf("You may now inspect it with the inspect command.\n")
	} else {
		fmt.Printf("Oh no! %s escaped the Pokeball!\n", pokemon.Name)
	}

	return nil
}

/*
	commandInspect takes the name of a Pokemon and prints the name, height, weight, stats and type(s) of the Pokemon. Example usage:

For example:
Pokedex > inspect pidgey
you have not caught that pokemon
Pokedex > catch pidgey
Throwing a Pokeball at pidgey...
pidgey was caught!
Pokedex > inspect pidgey
Name: pidgey
Height: 3
Weight: 18
Stats:

	-hp: 40
	-attack: 45
	-defense: 40
	-special-attack: 35
	-special-defense: 35
	-speed: 56

Types:
  - normal
  - flying
*/
func commandInspect(commands map[string]cliCommand, cfg *config, param []string) error {
	if len(param) == 0 {
		return fmt.Errorf("please specify a Pokemon to inspect")
	}

	pokemonName := param[0]
	pokemon, caught := cfg.caughtPokemons[pokemonName]
	if !caught {
		fmt.Printf("You have not caught that pokemon\n")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

// commandPokedex takes no parameters and lists all caught Pokemon by name.
// Pokedex > pokedex
//Your Pokedex:
// - pidgey
// - caterpie

func commandPokedex(commands map[string]cliCommand, cfg *config, param []string) error {
	if len(cfg.caughtPokemons) == 0 {
		fmt.Println("You have not caught any Pokemon yet.")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range cfg.caughtPokemons {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(commands map[string]cliCommand, cfg *config, param []string) error
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
	"map": {
		name:        "map",
		description: "Displays the map (scroll)",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Displays the map (backwards scroll)",
		callback:    commandMapb,
	},
	"explore": {
		name:        "explore",
		description: "Explores the specified location",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Catch a Pokemon (not implemented)",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "Inspect a caught Pokemon",
		callback:    commandInspect,
	},
	"pokedex": {
		name:        "pokedex",
		description: "List all caught Pokemon",
		callback:    commandPokedex,
	},
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var cfg config

	for fmt.Print("Pokedex > "); scanner.Scan(); fmt.Print("Pokedex > ") {
		command := scanner.Text()

		input := cleanInput(command)

		if len(input) == 0 {
			continue
		}

		cmd, exists := commands[input[0]]
		if exists {
			err := cmd.callback(commands, &cfg, input[1:])
			if err != nil {
				fmt.Printf("Error executing command '%s': %v\n", input[0], err)
				continue
			}

		}
	}

}
