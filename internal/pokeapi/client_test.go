package pokeapi

import "testing"

func TestGet(t *testing.T) {
	endpoints := []string{
		"pokemon/1",  // Bulbasaur
		"ability/65", // Overgrow
		"item/1",     // Master Ball
		"move/15",    // Cut
		"location/1", // Pallet Town
	}
	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			// Get expects a full URL. Prepend the base URL constant.
			url := PokeAPIBaseURL + endpoint
			data, err := Get(url)
			if err != nil {
				t.Fatalf("Expected no error for endpoint %s, got: %v", endpoint, err)
			}
			if len(data) == 0 {
				t.Errorf("Expected non-empty response for endpoint %s", endpoint)
			}
		})
	}
}
