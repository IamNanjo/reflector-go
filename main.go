package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/IamNanjo/logger-go"
)

var (
	statusUrl = "https://archlinux.org/mirrors/status/json/"
)

type Url struct {
	Url           string  `json:"url"`
	Protocol      string  `json:"protocol"`
	CompletionPct float32 `json:"completion_pct"`
	Score         float64 `json:"score"`
	Active        bool    `json:"active"`
	Country       string  `json:"country"`
	Details       string  `json:"details"`
}

type StatusResponse struct {
	Urls []Url `json:"urls"`
}

func getMirrorLists() *StatusResponse {
	res, err := http.Get(statusUrl)
	if err != nil {
		logger.Err.Fatalln("Could not get mirror lists")
	}

	defer res.Body.Close()

	parsedResponse := StatusResponse{}

	err = json.NewDecoder(res.Body).Decode(&parsedResponse)
	if err != nil {
		logger.Err.Fatalln("Could not parse body of mirror lists")
	}

	return &parsedResponse
}

// Filter, sort and limit results.
// Sorting happens based on score.
func filterResults(results []Url, protocol string, country string, number int) []Url {
	filtered := make([]Url, 0, len(results))

	countries := strings.Split(country, ",")

	for _, url := range results {
		matchesFilters := slices.ContainsFunc(countries, func(country string) bool {
			return url.Active && url.CompletionPct == 1 && // Only active mirrors with 100% completion
				slices.Contains(strings.Split(protocol, ","), url.Protocol) && // Protocol in allowed protocols list
				(country == "" || country == url.Country) // Country matches
		})

		if matchesFilters {
			filtered = append(filtered, url)
		}
	}

	filteredLen := len(filtered)
	limit := min(filteredLen, number)

	filtered = filtered[:limit]

	sort.SliceStable(filtered, func(i, j int) bool { return filtered[i].Score > filtered[j].Score })

	return filtered
}

func formatOutput(urls []Url) string {
	urlCount := len(urls)
	result := make([]string, len(urls))

	for i := range urlCount {
		url := urls[i]
		result[i] = fmt.Sprintf(
			"# Country=%s Score=%.2f\nServer = %s$repo/os/$arch",
			url.Country,
			url.Score,
			url.Url,
		)
	}

	return strings.Join(result, "\n\n")
}

func saveOutput(filepath string, urls string) {
	err := os.WriteFile(filepath, []byte(urls), 0644)
	if err != nil {
		logger.Err.Fatalln("Output writing failed")
	}
}

func main() {
	savePath := flag.String("save", "/etc/pacman.d/mirrorlist", "Save to file")
	protocol := flag.String("protocol", "https", "Comma separated list of protocols")
	country := flag.String("country", "Finland,Sweden,Estonia,Norway,Denmark,Latvia,Lithuania,Poland,Germany,France", "Comma separated list of countries")
	limit := flag.Int("limit", 20, "Maximum number of mirrorlists to use")
	flag.Parse()

	if savePath == nil || protocol == nil || country == nil || limit == nil {
		logger.Err.Fatalln("Command-line flag parsing failed")
	}

	response := getMirrorLists()
	outputUrls := filterResults(response.Urls, *protocol, *country, *limit)
	formattedOutput := formatOutput(outputUrls)

	logger.Debug.Println("\n" + formattedOutput)
	saveOutput(*savePath, formattedOutput)
	logger.Out.Printf("Saved %d mirrors into mirrorlist at %s", len(outputUrls), *savePath)
}
