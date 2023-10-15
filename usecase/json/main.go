package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type DataSSC struct {
	KeyRequest string `json:"@key"`
	Message    string `json:"@message"`
	Time       string `json:"@timestamp"`
}

type DataRequestSSC struct {
	KeyRequest string `json:"@key"`
	Message    string `json:"@message"`
	TimeStart  string `json:"@timeStart"`
}

type DataResponseSSC struct {
	KeyRequest string `json:"@key"`
	Message    string `json:"@message"`
	TimeEnd    string `json:"@timeEnd"`
}

type DataReqAndRespSSC struct {
	KeyRequest string `json:"@key"`
	Message    string `json:"@message"`
	TimeStart  string `json:"@TimeStart"`
	TimeEnd    string `json:"@TimeEnd"`
}

func buildKey(data []DataSSC) {
	// Open a file for writing (you can also create the file if it doesn't exist)
	file, err := os.Create("output.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Encode the data into JSON and write it to the file
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("Data written to output.json")
}

func readJson(fileName string) []DataSSC {
	var dataSlice []DataSSC

	// Open the JSON file for reading
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return dataSlice
	}
	defer file.Close()

	// Create a decoder for reading the JSON data
	decoder := json.NewDecoder(file)

	// Create a variable to store the decoded data

	// Decode the JSON data into the data variable
	if err := decoder.Decode(&dataSlice); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return dataSlice
	}
	// return dataSlice
	// var tempData []DataSSC
	for i, item := range dataSlice {
		fmt.Println("Message:", item.Message)
		fmt.Println()

		strArr := strings.Split(item.Message, ",")

		value := strArr[0]

		valueArr := strings.Split(value, " ")
		// itemTmp := DataSSC{
		// 	KeyRequest: valueArr[len(valueArr)-1],
		// 	Message:    strings.Join(valueArr[:(len(valueArr)-1)], " "),
		// 	// TimeStart: ,
		// 	Time: item.Time,
		// }
		// fmt.Println(itemTmp)

		// tempData = append(tempData, itemTmp)
		key := valueArr[len(valueArr)-1]
		dataSlice[i].KeyRequest = key
	}
	// buildKey(tempData)
	return dataSlice
}

func buildObjectKey(data []DataSSC) {
}

func mergeJson(filePath, filePath2, out string) []DataReqAndRespSSC {
	file1, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file1.json:", err)
		panic(err)
	}
	defer file1.Close()
	var data1 []DataReqAndRespSSC
	if err := json.NewDecoder(file1).Decode(&data1); err != nil {
		fmt.Println("Error decoding file1.json:", err)
		panic(err)

	}

	// Open and parse the second JSON file
	file2, err := os.Open(filePath2)
	if err != nil {
		fmt.Println("Error opening file2.json:", err)
		panic(err)

	}
	defer file2.Close()
	var data2 []DataReqAndRespSSC
	if err := json.NewDecoder(file2).Decode(&data2); err != nil {
		fmt.Println("Error decoding file2.json:", err)
		panic(err)

	}

	// Merge the two data slices
	// mergedData := append(data1, data2...)

	mergedData := make(map[string]DataReqAndRespSSC)
	// for _, item := range data1 {
	// 	mergedData[item.KeyRequest] = DataReqAndRespSSC{
	// 		KeyRequest: item.KeyRequest,
	// 		Message:    item.Message,
	// 		TimeStart:  item.TimeStart,
	// 		TimeEnd:    item.TimeEnd,
	// 	}
	// }
	// mergedData2 := make(map[string]DataReqAndRespSSC)

	for _, item := range data2 {
		mergedData[item.KeyRequest] = DataReqAndRespSSC{
			KeyRequest: item.KeyRequest,
			Message:    item.Message,
			TimeStart:  item.TimeStart,
			TimeEnd:    item.TimeEnd,
		}
	}

	// Convert the merged data back to a slice
	// var mergedDataSlice []DataReqAndRespSSC
	for i, item := range data1 {
		value, ok := mergedData[item.KeyRequest]
		if ok {
			data1[i].TimeEnd = value.TimeEnd
		} else {
			fmt.Println("errrrrr:", item.KeyRequest)
			// panic("")
		}
	}

	// Encode the merged data and write it to a new file
	mergedFile, err := os.Create(out)
	if err != nil {
		fmt.Println("Error creating merged.json:", err)
		panic(err)

	}
	defer mergedFile.Close()
	if err := json.NewEncoder(mergedFile).Encode(data1); err != nil {
		fmt.Println("Error encoding merged data:", err)
		panic(err)
	}

	fmt.Println("Merged data saved to merged.json")

	return data1
}

func writeJson(data interface{}, fileName string) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		panic(err)
	}
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		panic(err)

	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing JSON data to the file:", err)
		panic(err)

	}

	fmt.Println("Data written to output.json")
}

// find max, min, average
func find(data []DataReqAndRespSSC) {
	var (
		max    int64  = 0
		keyMax string = ""
		keyMin string = ""
		min    int64  = 10000
		sum    int64  = 0
	)
	for _, item := range data {
		startTime := item.TimeStart

		endTime := item.TimeEnd

		if startTime == "" || endTime == "" {
			continue
		}

		startT, err := time.Parse("2006-01-02 15:04:05.000", startTime)
		if err != nil {
			panic(err)
		}

		startE, err := time.Parse("2006-01-02 15:04:05.000", endTime)
		if err != nil {
			panic(err)
		}

		timeExec := startE.Sub(startT).Milliseconds()

		if timeExec > max {
			keyMax = item.KeyRequest
			max = timeExec
		}
		if timeExec < min {
			keyMin = item.KeyRequest
			min = timeExec
		}

		sum += timeExec
	}
	average := float64(sum) / float64(len(data))

	// Print the results
	fmt.Println("Maximum:", keyMax, max)
	fmt.Println("Minimum:", keyMin, min)
	fmt.Println("Average:", average)

}

func main() {
	dataCallRequest := readJson("call_time.json")
	// fmt.Println(dataCallRequest)
	var sscRequest []DataRequestSSC
	for _, item := range dataCallRequest {
		sscRequest = append(sscRequest, DataRequestSSC{
			KeyRequest: item.KeyRequest,
			Message:    item.Message,
			TimeStart:  item.Time,
		})
	}
	// fmt.Println(sscRequest)

	writeJson(sscRequest, "sscRequest.json")

	dataCallResponse := readJson("call_success_time.json")
	// fmt.Println(dataCallResponse)

	var sscResponse []DataResponseSSC
	for _, item := range dataCallResponse {
		sscResponse = append(sscResponse, DataResponseSSC{
			KeyRequest: item.KeyRequest,
			Message:    item.Message,
			TimeEnd:    item.Time,
		})
	}
	// fmt.Println(sscResponse)
	writeJson(sscResponse, "sscResponse.json")

	dataFinal := mergeJson("sscRequest.json", "sscResponse.json", "ssc_final.json")
	// fmt.Println(dataFinal)
	find(dataFinal)
}
