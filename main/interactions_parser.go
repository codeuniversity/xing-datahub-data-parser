package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/codeuniversity/xing-datahub-protocol"
)

type interaction protocol.RawInteraction
type interactionMap map[string]*interaction

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

	counter := 0
	fmt.Println("reading csv:...")
	for {
		counter++
		if counter%10000 == 0 {
			fmt.Print("line ", counter, "\r")
		}
		data, err := csvReader.Read()
		if data == nil {
			break
		} else {
			if err == nil {
				addToMap(parseArrayToInteraction(data))
			}
		}
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Writing:...")
	counter = 0
	csvWriter := csv.NewWriter(outputFile)
	csvWriter.Comma = '\t'
	csvWriter.Write([]string{"user_id", "item_id", "interaction_type", "created_at"})
	for _, i := range outputInteractions {
		counter++
		if counter%10000 == 0 {
			fmt.Print("line ", counter, "\r")
		}
		csvWriter.Write(i.toArray())
	}

}

func addToMap(i *interaction) {
	previousInteraction := outputInteractions[i.hashKey()]
	if previousInteraction == nil {
		outputInteractions[i.hashKey()] = i
	} else if previousInteraction.InteractionType != "4" {
		prev, _ := strconv.ParseInt(previousInteraction.InteractionType, 10, 64)
		current, _ := strconv.ParseInt(i.InteractionType, 10, 64)
		if prev < current {
			outputInteractions[i.hashKey()] = i
		}
	}
}

func parseArrayToInteraction(arr []string) *interaction {
	return &interaction{
		UserId:          arr[0],
		ItemId:          arr[1],
		InteractionType: arr[2],
		CreatedAt:       arr[3],
	}
}

func (i *interaction) toArray() []string {
	return []string{i.UserId, i.ItemId, i.InteractionType, i.CreatedAt}
}

func (i *interaction) hashKey() string {
	return i.UserId + "/" + i.ItemId
}
