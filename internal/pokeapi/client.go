package pokeapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Get makes a GET request to the given PokeAPI endpoint and returns the response body as bytes.
func Get(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/%s", endpoint)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PokeAPI request failed: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
