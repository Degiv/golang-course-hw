package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

type operationType string

const (
	zeroValue operationType = ""
	income    operationType = "income"
	outcome   operationType = "outcome"
	plus      operationType = "+"
	minus     operationType = "-"
)

type Operation struct {
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Id        interface{} `json:"id"`
	CreatedAt string      `json:"created_at"`
}

type CompanyOperation struct {
	Company        string    `json:"company"`
	InnerOperation Operation `json:"operation"`
	Operation
}

func (operation *CompanyOperation) fill() {
	if operation.InnerOperation.Id != nil {
		operation.Id = operation.InnerOperation.Id
	}

	if operation.InnerOperation.Value != nil {
		operation.Value = operation.InnerOperation.Value
	}

	if operation.InnerOperation.CreatedAt != "" {
		operation.CreatedAt = operation.InnerOperation.CreatedAt
	}

	if operation.InnerOperation.Type != "" {
		operation.Type = operation.InnerOperation.Type
	}
}

type CompanySummary struct {
	Company             string        `json:"company"`
	ValidOperationCount int           `json:"valid_operations_count"`
	Balance             int           `json:"balance"`
	InvalidOperations   []interface{} `json:"invalid_operations,omitempty"`
}

func (cs *CompanySummary) update(operation CompanyOperation) {
	operationType, typeOk := parseOperationType(operation.Type)
	value, valueOk := parseValue(operation.Value)

	isValid := typeOk && valueOk

	if isValid {
		cs.ValidOperationCount++
		cs.Balance = apply(cs.Balance, operationType, value)
	} else {
		if cs.InvalidOperations == nil {
			cs.InvalidOperations = []interface{}{}
		}
		cs.InvalidOperations = append(cs.InvalidOperations, operation.Id)
	}
}

func getPathFromStdin() string {
	failPath, _ := io.ReadAll(os.Stdin)

	if failPath == nil {
		return ""
	}

	return string(failPath)
}

func getPathFromENV() string {
	filePath, _ := os.LookupEnv("FILE")
	return filePath
}

func getFilePath() string {
	var filePathPtr *string = flag.String("file", "", "Path to file")
	flag.Parse()
	filePath := *filePathPtr
	if filePath != "" {
		return filePath
	}

	filePath = getPathFromENV()
	if filePath != "" {
		return filePath
	}

	filePath = getPathFromStdin()
	return filePath
}

func parseOperationType(operationType string) (operationType, bool) {
	switch operationType {
	case string(income):
		return income, true
	case string(outcome):
		return outcome, true
	case string(plus):
		return plus, true
	case string(minus):
		return minus, true
	default:
		return zeroValue, false
	}
}

func parseValue(value interface{}) (int, bool) {
	intValue, err := strconv.Atoi(fmt.Sprint(value))
	return intValue, err == nil
}

func idOk(id interface{}) bool {
	switch id.(type) {
	case string:
		return true
	case float64:
		_, err := strconv.Atoi(fmt.Sprint(id))
		return err == nil
	default:
		return false
	}
}

func dateOk(createdAt string) bool {
	_, err := time.Parse(time.RFC3339, createdAt)
	return err == nil
}

func companyOk(company string) bool {
	return company != ""
}

func apply(balance int, operation operationType, value int) int {
	switch operation {
	case income, plus:
		return balance + value
	case outcome, minus:
		return balance - value
	default:
		return 0
	}
}

func NewCompanySummary(operation CompanyOperation) *CompanySummary {
	operationType, typeOk := parseOperationType(operation.Type)
	value, valueOk := parseValue(operation.Value)

	isValid := typeOk && valueOk

	newSummary := CompanySummary{
		Company:             operation.Company,
		ValidOperationCount: 0,
		Balance:             0,
		InvalidOperations:   nil,
	}

	if isValid {
		newSummary.ValidOperationCount++
		newSummary.Balance = apply(newSummary.Balance, operationType, value)
	} else {
		newSummary.InvalidOperations = []interface{}{operation.Id}
	}

	return &newSummary

}

func getSummaries(companyOperations []CompanyOperation) []CompanySummary {
	for i := range companyOperations {
		operation := &companyOperations[i]
		operation.fill()
	}

	sort.SliceStable(companyOperations, func(i, j int) bool {
		return companyOperations[i].CreatedAt < companyOperations[j].CreatedAt
	})
	companySummariesMap := map[string]*CompanySummary{}
	for _, operation := range companyOperations {
		if !(companyOk(operation.Company) && idOk(operation.Id) && dateOk(operation.CreatedAt)) {
			continue
		}
		_, ok := companySummariesMap[operation.Company]
		if ok {
			companySummariesMap[operation.Company].update(operation)
		} else {
			if companyOk(operation.Company) && dateOk(operation.CreatedAt) && idOk(operation.Id) {
				companySummariesMap[operation.Company] = NewCompanySummary(operation)
			}
		}
	}
	companySummaries := make([]CompanySummary, 0)
	for _, value := range companySummariesMap {
		companySummaries = append(companySummaries, *value)
	}
	return companySummaries
}

func main() {
	filePath := getFilePath()
	if filePath == "" {
		fmt.Println("No file")
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(file)
	var operations []CompanyOperation
	err = json.Unmarshal(data, &operations)

	summaries := getSummaries(operations)
	sort.Slice(summaries, func(i int, j int) bool {
		return summaries[i].Company < summaries[j].Company
	})
	out, _ := os.Create("hw2/out.json")
	marshaled, err := json.MarshalIndent(summaries, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := out.Write(marshaled); err != nil {
		log.Fatal(err)
	}
}
