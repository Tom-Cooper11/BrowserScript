package createAlert

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"encoding/json"
	"time"
)

type DataMap map[string]map[string]interface{}

type responseData struct {
	Data bool `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type responseData1 struct {
	Data string `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type uploadResps struct {
	Data struct {
		Id 	 string `json:"id"`
		Name string `json:"name"`
		ReferenceId string `json:"referenceId"`
	}
	Msg string `json:"msg"`
}

type JsonData struct {
	EventName  string `json:"eventName"`
	EventLevel string `json:"eventLevel"`
	Type 	 string `json:"type"`
	SubType    string `json:"subType"`
	StartTime  string `json:"startTime"`
	EventDescription string `json:"eventDescription"`
	SrcIps     []string `json:"srcIps"`
	Asset      []string `json:"asset"`
	UnkownAsset []interface{} `json:"unkownAsset"`
	FileList   []string `json:"fileList"`
	ResponseMeasures string `json:"responseMeasures"`
}

type PostData1 struct {
	EventId        string   `json:"eventId"`
	JudgeType      string   `json:"judgeType"`
	DealDeadline   string   `json:"dealDeadline"`
	Require        string   `json:"require"`
	DealPrincipal  string   `json:"dealPrincipal"`
	DealAssist     []string `json:"dealAssist"`
}

type PostData2 struct {
	FileList         []string `json:"fileList"`
	JudgeConclusion  string   `json:"judgeConclusion"`
	EventId		     string   `json:"eventId"`
	HandleSuggestion string   `json:"handleSuggestion"`
}

type PostData3 struct {
	EventId   string `json:"eventId"`
	TaskType  string `json:"taskType"`
	JudgeType string `json:"judgeType"`
	TaskName  string `json:"taskName"`
	DealUser  string `json:"dealUser"`
}

type SourceData struct {
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
	Total       int    `json:"total"`
	EventId     string `json:"eventId"`
	CreateType  string `json:"createType"`
}

type SaveAlert struct {
	FlowTaskId         string `json:"flowTaskId"`
	FormInfo           string `json:"formInfo"`
	TaskId             string `json:"taskId"`
	TemplateId         string `json:"templateId"`
	PlanTemplateId     string `json:"planTemplateId"`
	TaskType           int    `json:"taskType"`
	TaskName           string `json:"taskName"`
	TemplateInstanceId string `json:"templateInstanceId"`
	TaskSubType        string `json:"taskSubType"`
	Assignee           string `json:"assignee"`
	IsApprove          bool   `json:"isApprove"`
}

type Info struct {
	ContactUserName    string `json:"contactUserName"`
	ContactUserPhone   string `json:"contactUserPhone"`
	TemplateId         string `json:"templateId"`
	PlanTemplateId     string `json:"planTemplateId"`
	TaskType           int    `json:"taskType"`
	TemplateInstanceId string `json:"templateInstanceId"`
	TaskSubType        string `json:"taskSubType"`
	Assignee           string `json:"assignee"`
	IsApprove          bool   `json:"isApprove"`
	DocxPath		 []interface{} `json:"docxPath"`
}

func Rand_Sleep() {
	rand.Seed(time.Now().UnixNano())
	sleepDuration := rand.Intn(10) + 11
	time.Sleep(time.Duration(sleepDuration) * time.Second)
}

func Find(paths []interface{}, key string) string {
	for _, path := range paths {
		DocxPath := filepath.Join(path.(string), fmt.Sprintf("%s.docx", key))
		_, err := os.Stat(DocxPath)
		if err == nil {
			fmt.Println("File exists:", DocxPath)
			return DocxPath
		}
	}
	return ""
}

func readJsonFile(filePath string) (DataMap) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	
	jsonData := make(DataMap)
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		os.Exit(1)
	}

	return jsonData
}

func writeJsonFile(filePath string, jsonData DataMap) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(jsonData)
	if err != nil {
		fmt.Println("Error encoding file:", err)
		os.Exit(1)
	}

	err = os.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}
}

func uploadFile(filePath string) (*bytes.Buffer, string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		fmt.Println("Error creating form file:", err)
		os.Exit(1)
	}

	_, err = io.Copy(part, bytes.NewReader(file))
	if err != nil {
		fmt.Println("Error copying file data:", err)
		os.Exit(1)
	}
	writer.Close()
	contentType := writer.FormDataContentType()

	return body, contentType
}

func getPostfix(key string) string {
	runes := []rune(key)
	length := len(runes)
	if length < 4 {
		return ""
	}
	postfix := string(runes[length-4:])
	return postfix
}