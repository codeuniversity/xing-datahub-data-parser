package main

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"

	"github.com/codeuniversity/xing-datahub-protocol"
)

type interaction protocol.RawInteraction
type interactionCollection []*interaction
type interactionMap map[string]*interaction

var interactions interactionCollection
var outputInteractions = interactionMap{}

func main() {
	filepath := os.Getenv("file")
	if filepath == "" {
		panic("file not set")
	}
	outputFilePath := os.Getenv("output")
	if outputFilePath == "" {
		outputFilePath = "interactionsCleaned.csv"
	}

	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	csvReader.LazyQuotes = false
	csvReader.ReuseRecord = true

	for {
		data, err := csvReader.Read()
		if data == nil {
			break
		} else {
			if err == nil {
				addInteraction(data)
			}
		}
	}

	sort.Sort(interactions)
	addToMap(interactions)

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}
	csvWriter := csv.NewWriter(outputFile)
	csvWriter.Comma = '\t'
	for _, i := range outputInteractions {
		csvWriter.Write(i.toArray())
	}

}

func addToMap(c interactionCollection) {
	for _, i := range c {
		previousInteraction := outputInteractions[i.hashKey()]
		if previousInteraction == nil || previousInteraction.InteractionType != "4" {
			outputInteractions[i.hashKey()] = i
		}
	}
}

func addInteraction(arr []string) {
	interaction := parseArrayToInteraction(arr)
	interactions = append(interactions, interaction)
}

func parseArrayToInteraction(arr []string) *interaction {
	return &interaction{
		UserId:          arr[0],
		ItemId:          arr[1],
		InteractionType: arr[2],
		CreatedAt:       arr[3],
	}
}

func (c interactionCollection) Len() int {
	return len(c)
}

func (c interactionCollection) Less(i, j int) bool {
	createdAt1, err := strconv.ParseInt(c[i].CreatedAt, 10, 64)
	createdAt2, err := strconv.ParseInt(c[j].CreatedAt, 10, 64)

	if err != nil {
		return false
	}
	return createdAt1 < createdAt2
}

func (c interactionCollection) Swap(i, j int) {
	var tmpPointer *interaction
	tmpPointer = c[i]
	c[i] = c[j]
	c[j] = tmpPointer
}

func (i *interaction) toArray() []string {
	return []string{i.UserId, i.ItemId, i.InteractionType, i.CreatedAt}
}

func (i *interaction) hashKey() string {
	return i.UserId + "/" + i.ItemId
}
