package apiget

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/luismingati/pokedexcli/internal/pokecache"
)

var (
	cache    *pokecache.Cache
	once     sync.Once
	pokeball *pokecache.Cache
)

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type PokeEncounters struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type AreasResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func initCache() {
	once.Do(func() {
		cache = pokecache.NewCache(time.Second * 100)
		pokeball = pokecache.NewCache(time.Second * 10000)
	})
}

func GetAreas(offset int) error {
	initCache()
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=%d&limit=20", offset)

	if data, ok := cache.Get(url); ok {
		var goData AreasResponse
		if err := json.Unmarshal(data, &goData); err != nil {
			return fmt.Errorf("failed to unmarshal response from cache: %v", err)
		}
		for _, result := range goData.Results {
			fmt.Println(result.Name)
		}
		return nil
	}

	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	cache.Add(url, body)

	var goData AreasResponse

	if err := json.Unmarshal(body, &goData); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	for _, result := range goData.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func Explore(region string) error {
	initCache()

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", region)

	if data, ok := cache.Get(url); ok {
		var goData PokeEncounters
		if err := json.Unmarshal(data, &goData); err != nil {
			return fmt.Errorf("failed to unmarshal response from cache: %v", err)
		}
		fmt.Printf("Exploring %s...\n", region)
		for _, result := range goData.PokemonEncounters {
			fmt.Println(result.Pokemon.Name)
		}
		return nil
	}

	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	cache.Add(url, body)

	var goData PokeEncounters

	if err := json.Unmarshal(body, &goData); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	fmt.Printf("Exploring %s...\n", region)
	for _, result := range goData.PokemonEncounters {
		fmt.Println(result.Pokemon.Name)
	}

	return nil
}

// https://pokeapi.co/api/v2/pokemon/{id or name}/

func Catch(pokename string) error {
	initCache()

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokename)

	if _, ok := pokeball.Get(pokename); ok {
		fmt.Printf("You already have %s!\n", pokename)
		return nil
	}

	if data, ok := cache.Get(url); ok {
		var goData Pokemon

		if err := json.Unmarshal(data, &goData); err != nil {
			return fmt.Errorf("failed to unmarshal response from cache: %v", err)
		}
		fmt.Printf("Throwing a pokeball at %s...\n", pokename)
		rand := rand.Intn(goData.BaseExperience)

		if rand < goData.BaseExperience/2 {
			fmt.Printf("%s escaped\n", goData.Name)
			return nil
		}

		fmt.Printf("%s was caught!\n", goData.Name)
		pokeball.Add(goData.Name, data)

		return nil
	}

	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	cache.Add(url, body)

	var goData Pokemon

	if err := json.Unmarshal(body, &goData); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	fmt.Printf("Throwing a pokeball at %s...\n", pokename)
	rand := rand.Intn(goData.BaseExperience)

	if rand < goData.BaseExperience/2 {
		fmt.Printf("%s escaped\n", goData.Name)
		return nil
	}

	fmt.Printf("%s was caught!\n", goData.Name)
	pokeball.Add(goData.Name, body)

	return nil
}

func Inspect(pokename string) error {
	initCache()

	if data, ok := pokeball.Get(pokename); ok {
		var goData Pokemon

		if err := json.Unmarshal(data, &goData); err != nil {
			return fmt.Errorf("failed to unmarshal response from cache: %v", err)
		}
		fmt.Printf("Name: %s\n", goData.Name)
		fmt.Printf("Height: %d\n", goData.Height)
		fmt.Printf("Weight: %d\n", goData.Weight)
		fmt.Println("Stats:")
		for _, stat := range goData.Stats {
			fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range goData.Types {
			fmt.Printf("  -%s\n", t.Type.Name)
		}
		return nil
	}

	fmt.Println("You have not caught this pokemon yet!")
	return nil
}

func Pokedex() error {
	initCache()
	fmt.Println("Your pokedex:")
	pokeball.List()
	return nil
}
