package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/luismingati/pokedexcli/internal/apiget"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args ...string) error
}

var cmd map[string]cliCommand

func init() {
	currentOffset := 0

	cmd = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "displays the names of 20 location areas in the Pokemon world",
			callback: func(args ...string) error {
				apiget.GetAreas(currentOffset)
				currentOffset += 20
				return nil
			},
		},
		"mapb": {
			name:        "mapb",
			description: "displays the previous 20 locations",
			callback: func(args ...string) error {
				if currentOffset >= 20 {
					if currentOffset > 20 {
						currentOffset -= 20
					}
					currentOffset -= 20
					return apiget.GetAreas(currentOffset)
				}
				return fmt.Errorf("you already are on the first page")
			},
		},
		"explore": {
			name:        "explore",
			description: "Shows the pokemons that can be found in a location area",
			callback: func(args ...string) error {
				if len(args) == 0 {
					return fmt.Errorf("you need to send the location area")
				}
				apiget.Explore(args[0])
				return nil
			},
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			callback: func(args ...string) error {
				if len(args) == 0 {
					return fmt.Errorf("you need to send the pokemon name")
				}
				apiget.Catch(args[0])
				return nil
			},
		},
		"inspect": {
			name:        "inspect",
			description: "show all the pokemon data and stats",
			callback: func(args ...string) error {
				if len(args) == 0 {
					return fmt.Errorf("you need to send the pokemon name")
				}
				apiget.Inspect(args[0])
				return nil
			},
		},
		"pokedex": {
			name:        "pokedex",
			description: "show all the pokemon in the pokedex",
			callback: func(args ...string) error {
				apiget.Pokedex()
				return nil
			},
		},
	}
}

func commandHelp(args ...string) error {
	fmt.Println("\nWelcome to the pokedexCLI!")
	fmt.Println("Usage:")

	for _, c := range cmd {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	fmt.Println()
	return nil
}

func commandExit(args ...string) error {
	os.Exit(0)
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		input = strings.TrimSpace(input)
		parts := strings.Split(input, " ")

		if len(parts) > 0 {
			command := parts[0]
			args := parts[1:]

			if cmd, ok := cmd[command]; ok {
				if err := cmd.callback(args...); err != nil {
					fmt.Println("Error executing command:", err)
				}
			} else {
				fmt.Println("Unknown command:", input)
			}
		}
	}
}
