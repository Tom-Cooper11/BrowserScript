package createAlert

import (
	"alert/encrypt"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var days = time.Now()
var day = days.Format("20060102")
var dayReplys = days.AddDate(0, 0, 7)
var dayReply = dayReplys.Format("2006-01-02 15:04:05")

func Alert_Follow(auth string, execDir string) {
	dbPath := filepath.Join(execDir, "db.json")
	configPath := filepath.Join(execDir, "configuration.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("configuration.json file does not exist in the executable directory.")
		return
	}
	config := readJsonFile(configPath)

	jsonPath := config["path"]["jsonPath"].(string)
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		fmt.Println("postData.json file does not exist in the executable directory.")
		return
	}
	jsonData := readJsonFile(jsonPath)

	url1 := "http://10.76.211.100/judge/api/judge/v1.0/create"
	
	if data, ok := jsonData[day]; ok {
		for _, value := range data {
			valueMap := value.(map[string]interface{})
			var valueJson JsonData
			b1, err := json.Marshal(value)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return
			}

			err = json.Unmarshal(b1, &valueJson)
			if err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				return
			}
			jsonData1, err := json.Marshal(valueJson)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return
			}

			request, err := http.NewRequest("POST", url1, bytes.NewBuffer(jsonData1))
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}

			var timestamps = time.Now().UnixMilli()
			var timestamp = strconv.FormatInt(timestamps, 10)
			encryptData := &encrypt.Request{
				URL: url1,
				Params: url.Values{},
				Data: valueJson,
			}
			sign := encrypt.GenerateSign(encryptData, timestamp)

			request.Header.Set("Sign", sign)
			request.Header.Set("x-authorization", auth)
			request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux aarch64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
			request.Header.Set("Accept", "application/json, text/plain, */*")
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Origin", "http://10.76.211.100")
			request.Header.Set("Referer", "http://10.76.211.100/web-judgem/system/createEvent?dataIndex=0&title=%E6%96%B0%E5%A2%9E%E4%BA%8B%E4%BB%B6&source=")			
			request.Header.Set("timestamp", timestamp)
			request.Header.Set("Cookie", "token_flag=1")
			request.Header.Set("accept-encoding", "gzip, deflate")
			request.Header.Set("accept-language", "zh")

			client := &http.Client{
				Timeout: 5 * time.Second,
			}
			response, err := client.Do(request)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				fmt.Println("Error response from create server:", response.Status)
				return
			}

			var respData responseData1
			err = json.NewDecoder(response.Body).Decode(&respData)
			if err != nil {
				fmt.Println("Error decoding response JSON:", err)
				return
			}
			if respData.Msg != "success" {
				fmt.Println("Error response from server:", respData)
				return
			}
			valueMap["eventId"] = respData.Data
			fmt.Printf("Event create success: %s\n", valueMap["eventName"])
			Rand_Sleep()
		}
	}

	jsonData2 := make(DataMap)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		jsonData2[day] = jsonData[day]
	} else {
		jsonData2 = readJsonFile(dbPath)
		jsonData2[day] = jsonData[day]
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(jsonData2)
	if err != nil {
		fmt.Println("Error encoding file:", err)
		return
	}

	err = os.WriteFile(dbPath, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	url2 := "http://10.76.211.100/judge/api/judge/v1.0/start"
	for _, value := range jsonData[day] {
		valueMap := value.(map[string]interface{})
		eventId := valueMap["eventId"].(string)
		uneventName := valueMap["eventName"]
		ueventName := url.QueryEscape(uneventName.(string))
		eventName := url.QueryEscape(ueventName)
		referUrl := fmt.Sprintf("http://10.76.211.100/web-judgem/system/event_manage_detail?title=%s&eventId=%s&taskStatus=0&source=&status=0", eventName, eventId)

		postData1 := PostData1{
			EventId:       eventId,
			JudgeType:     "1",
			DealDeadline:  dayReply,
			Require:      config["judge"]["require"].(string),
			DealPrincipal: "1600311717",
			DealAssist:    []string{},
		}

		jsonData2, err := json.Marshal(postData1)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req1, err := http.NewRequest("POST", url2, bytes.NewReader(jsonData2))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		var timestamps1 = time.Now().UnixMilli()
		var timestamp1 = strconv.FormatInt(timestamps1, 10)
		encryptData := &encrypt.Request{
				URL: url2,
				Params: url.Values{},
				Data: postData1,
			}
		sign := encrypt.GenerateSign(encryptData, timestamp1)

		req1.Header.Set("Sign", sign)
		req1.Header.Set("x-authorization", auth)
		req1.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux aarch64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
		req1.Header.Set("Accept", "application/json, text/plain, */*")
		req1.Header.Set("Content-Type", "application/json")
		req1.Header.Set("Origin", "http://10.76.211.100")
		req1.Header.Set("Referer", referUrl)
		req1.Header.Set("timestamp", timestamp1)
		req1.Header.Set("Cookie", "token_flag=1")
		req1.Header.Set("accept-encoding", "gzip, deflate")
		req1.Header.Set("accept-language", "zh")

		client1 := &http.Client{
			Timeout: 5 * time.Second,
		}
		resp1, err := client1.Do(req1)
		if err != nil {
			fmt.Println("Error sending to judge request:", err)
			return
		}
		defer resp1.Body.Close()
		if resp1.StatusCode != http.StatusOK {
			fmt.Println("Error response from submit server:", resp1.Status)
			return
		}

		var respData1 responseData
		err = json.NewDecoder(resp1.Body).Decode(&respData1)
		if err != nil {
			fmt.Println("Error decoding response JSON:", err)
			return
		}
		if respData1.Msg != "success" {
			fmt.Println("Error json from submit server:", respData1)
			return
		}
		fmt.Printf("Event submit success: %s\n", valueMap["eventName"])
		Rand_Sleep()
	}
}

func Alert_Judge(auth string, execDir string) {
	dbPath := filepath.Join(execDir, "db.json")
	configPath := filepath.Join(execDir, "configuration.json")

	jsonData := readJsonFile(dbPath)
	config := readJsonFile(configPath)
	docxPath := config["path"]["docxPath"].([]interface{})
	
	if data, ok := jsonData[day]; ok {
		for key, value := range data {
			valueMap := value.(map[string]interface{})
			eventId := valueMap["eventId"].(string)
			uneventName := valueMap["eventName"]
			ueventName := url.QueryEscape(uneventName.(string))
			eventName := url.QueryEscape(ueventName)
			referUrl := fmt.Sprintf("http://10.76.211.100/web-judgem/system/event_manage_detail_support?title=%s&eventId=%s&taskStatus=0&source=plat&status=1", eventName, eventId)

			var DocxPath = Find(docxPath, key)
			if DocxPath == "" {
				fmt.Println("File does not exist:", key, DocxPath)
				return
			}

			body, contentType := uploadFile(DocxPath)

			url1 := "http://10.76.211.100/system/api/minio/v1.0/uploadCipher"

			req, err := http.NewRequest("POST", url1, body)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}
			var timestamps = time.Now().UnixMilli()
			var timestamp = strconv.FormatInt(timestamps, 10)

			req.Header.Set("Content-Type", contentType)
			req.Header.Set("x-authorization", auth)
			req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux aarch64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.5359.125 Safari/537.36 UOS Professional")
			req.Header.Set("Accept", "*/*")			
			req.Header.Set("Origin", "http://10.76.211.100")
			req.Header.Set("Referer", referUrl)
			req.Header.Set("timestamp", timestamp)
			req.Header.Set("Cookie", "token_flag=1")

			client := &http.Client{
				Timeout: 10 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error response from server:", resp.Status)
				return
			}
			var uploadResp uploadResps
			err = json.NewDecoder(resp.Body).Decode(&uploadResp)
			if err != nil {
				fmt.Println("Error decoding response JSON:", err)
				return
			}
			var fileList = uploadResp.Data.Id
			valueMap["fileList"] = fileList
			fmt.Printf("File upload success: %s\n", key)
			Rand_Sleep()

			url2 := "http://10.76.211.100/judge/api/judge/v1.0/judge"
			postData1 := PostData2 {
				FileList: []string{
					fileList,
				},
				JudgeConclusion: config["judge"]["conclusion"].(string),
				EventId: eventId,
				HandleSuggestion: config["judge"][valueMap["type"].(string)].(string),
			}

			jsonData1, err := json.Marshal(postData1)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return
			}
			req1, err := http.NewRequest("POST", url2, bytes.NewReader(jsonData1))
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}

			var timestamps1 = time.Now().UnixMilli()
			var timestamp1 = strconv.FormatInt(timestamps1, 10)
			encryptData := &encrypt.Request{
				URL: url2,
				Params: url.Values{},
				Data: postData1,
			}
			sign := encrypt.GenerateSign(encryptData, timestamp1)

			req1.Header.Set("x-authorization", auth)
			req1.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux aarch64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.5359.125 Safari/537.36 UOS Professional")
			req1.Header.Set("Accept", "*/*")
			req1.Header.Set("Content-Type", "application/json")
			req1.Header.Set("Origin", "http://10.76.211.100")
			req1.Header.Set("Referer", referUrl)
			req1.Header.Set("timestamp", timestamp1)
			req1.Header.Set("Cookie", "token_flag=1")
			req1.Header.Set("Sign", sign)
			req1.Header.Set("accept-encoding", "gzip, deflate")
			req1.Header.Set("accept-language", "zh")

			client1 := &http.Client{
				Timeout: 5 * time.Second,
			}
			resp1, err := client1.Do(req1)
			if err != nil {
				fmt.Println("Error sending to judge request:", err)
				return
			}
			defer resp1.Body.Close()
			if resp1.StatusCode != http.StatusOK {
				fmt.Println("Error response from server:", resp1.Status)
				return
			}

			var respData1 responseData
			err = json.NewDecoder(resp1.Body).Decode(&respData1)
			if err != nil {
				fmt.Println("Error decoding response JSON:", err)
				return
			}
			if respData1.Msg != "success" {
				fmt.Println("Error response from server:", respData1.Msg)
				return
			}
			fmt.Printf("Event judge completed: %s\n", uneventName)
			Rand_Sleep()
		}
	}
	
	writeJsonFile(dbPath, jsonData)
}

func Alert_Task(auth string, execDir string) {
	dbPath := filepath.Join(execDir, "db.json")
	jsonData := readJsonFile(dbPath)

	url1 := "http://10.76.211.100/instruct/api/task/v1.0/create"

	if data, ok := jsonData[day]; ok {
		for _, value := range data {
			valueMap := value.(map[string]interface{})
			eventId := valueMap["eventId"].(string)
			uneventName := valueMap["eventName"]
			ueventName := url.QueryEscape(uneventName.(string))
			eventName := url.QueryEscape(ueventName)
			referUrl := fmt.Sprintf("http://10.76.211.100/web-judgem/system/event_manage_detail?title=%s&eventId=%s&taskStatus=0&source=&status=1", eventName, eventId)
			postData := PostData3 {
				DealUser:  "1060513426117365760",
				EventId:   eventId,
				JudgeType: "",
				TaskName:  fmt.Sprintf("关于%s的通报", uneventName),
				TaskType:  "2",
			}

			jsonData1, err := json.Marshal(postData)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return
			}
			req, err := http.NewRequest("POST", url1, bytes.NewReader(jsonData1))
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}

			var timestamps = time.Now().UnixMilli()
			var timestamp = strconv.FormatInt(timestamps, 10)
			encryptData := &encrypt.Request{
				URL: url1,
				Params: url.Values{},
				Data: postData,
			}
			var sign = encrypt.GenerateSign(encryptData, timestamp)

			req.Header.Set("x-authorization", auth)
			req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux aarch64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Origin", "http://10.76.211.100")
			req.Header.Set("Referer", referUrl)
			req.Header.Set("timestamp", timestamp)
			req.Header.Set("Cookie", "token_flag=1")
			req.Header.Set("Sign", sign)
			req.Header.Set("accept-encoding", "gzip, deflate")
			req.Header.Set("accept-language", "zh")
			
			client := &http.Client{
				Timeout: 5 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error response from server:", resp.Status)
				return
			}
			
			var respData1 responseData1
			err = json.NewDecoder(resp.Body).Decode(&respData1)
			if err != nil {
				fmt.Println("Error decoding response JSON:", err)
				return
			}
			if respData1.Msg != "success" {
				fmt.Println("Error response from server:", respData1.Msg)
				return
			}
			valueMap["taskId"] = respData1.Data
			fmt.Printf("Task create success: %s\n", uneventName)
			Rand_Sleep()
		}
	}
	
	writeJsonFile(dbPath, jsonData)
}

