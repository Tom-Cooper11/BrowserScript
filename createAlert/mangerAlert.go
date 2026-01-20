package createAlert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"
)

import "alert/encrypt"

func Alert_Save(auth string, execDir string) {
	dbPath := filepath.Join(execDir, "db.json")
	configPath := filepath.Join(execDir, "configuration.json")

	jsonData := readJsonFile(dbPath)
	config := readJsonFile(configPath)
	docxPath := config["path"]["docxPath"].([]interface{})
	
	if data, ok := jsonData[day]; ok {
		for key, value := range data {
			valueMap := value.(map[string]interface{})
			uneventName := valueMap["eventName"]
			taskId := valueMap["taskId"].(string)

			var DocxPath = Find(docxPath, key)
			if DocxPath == "" {
				fmt.Println("File does not exist:", key, DocxPath)
				return
			}

			body, contentType := uploadFile(DocxPath)

			url1 := "http://10.76.211.100/system/api/minio/v1.0/sync/upload"

			req, err := http.NewRequest("POST", url1, body)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}

			req.Header.Set("Content-Type",contentType)
			req.Header.Set("x-authorization", auth)
			req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux aarch64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
			req.Header.Set("Accept", "*/*")			
			req.Header.Set("Origin", "http://10.76.211.100")
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
			var fileId = uploadResp.Data.Id
			valueMap["fileId"] = fileId
			fmt.Printf("File upload success: %s\n", key)
			Rand_Sleep()

			var timestamps1 = time.Now().UnixMilli()
			var timestamp1 = strconv.FormatInt(timestamps1, 10)

			var info Info
			b1, err := json.Marshal(config["info"])
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return
			}

			err = json.Unmarshal(b1, &info)
			if err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				return
			}

			taskDesc := valueMap["eventDescription"].(string)
			infectLocation := valueMap["unkownAsset"].([]interface{})[0].(map[string]interface{})["ip"].(string)
			contactUserName := info.ContactUserName
			contactUserPhone := info.ContactUserPhone
			postfix := getPostfix(key)
			var influence, dealSuggestion string
			if _, exist := config["risk"][postfix]; exist {
				influence = config["risk"][postfix].(string)
				dealSuggestion = config["judge"]["IN01"].(string)
			} else {
				influence = config["risk"]["vuln"].(string)
				dealSuggestion = config["judge"]["IN02"].(string)
			}
			
			url2 := "http://10.76.211.100/instruct/api/task/v1.0/arrange/save"
			postData1 := SaveAlert {
				FlowTaskId: "",
				FormInfo: fmt.Sprintf("[{\"id\":\"fileCode\",\"value\":null},{\"id\":\"security_level_normal\",\"value\":null},{\"id\":\"multi_broRange\",\"value\":null},{\"id\":\"manulBroRange\",\"value\":null},{\"id\":\"taskDesc\",\"value\":\"%s\"},{\"id\":\"infectLocation\",\"value\":\"%s\"},{\"id\":\"influence\",\"value\":\"%s\"},{\"id\":\"dealSuggestion\",\"value\":\"%s\"},{\"id\":\"datetime_feedbackTime\",\"value\":\"%s\"},{\"id\":\"multi_ccOrg\",\"value\":null},{\"id\":\"multi_publishMethod\",\"value\":[\"平台发布\",\"即时通讯\"]},{\"id\":\"contactUserName\",\"value\":\"%s\"},{\"id\":\"contactUserLandlineNum\",\"value\":\"\"},{\"id\":\"contactUserPhone\",\"value\":\"%s\"},{\"id\":\"pubOrgName\",\"value\":\"四川省网络安全应急办公室\"},{\"id\":\"report_attachment\",\"value\":[{\"id\":\"%s\",\"name\":\"%s.docx\",\"referenceId\":null,\"uid\":%s,\"status\":\"success\"}]}]", taskDesc, infectLocation, influence, dealSuggestion, dayReply, contactUserName, contactUserPhone, fileId, key, timestamp1),
				TaskId: taskId,
				TemplateId: info.TemplateId,
				PlanTemplateId: info.PlanTemplateId,
				TaskType: info.TaskType,
				TaskName: fmt.Sprintf("关于%s的通报", uneventName),
				TemplateInstanceId: info.TemplateInstanceId,
				TaskSubType: info.TaskSubType,
				Assignee: info.Assignee,
				IsApprove: info.IsApprove,
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

			encryptData := &encrypt.Request{
				URL: url2,
				Params: url.Values{},
				Data: postData1,
			}
			sign := encrypt.GenerateSign(encryptData, timestamp1)

			req1.Header.Set("x-authorization", auth)
			req1.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux aarch64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
			req1.Header.Set("Accept", "application/json, text/plain, */*")
			req1.Header.Set("Content-Type", "application/json")
			req1.Header.Set("Origin", "http://10.76.211.100")
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
				fmt.Println("Status Error response from server:", resp1.Status)
				return
			}

			var respData1 responseData1
			err = json.NewDecoder(resp1.Body).Decode(&respData1)
			if err != nil {
				fmt.Println("Error decoding response JSON:", err)
				return
			}
			if respData1.Msg != "success" {
				fmt.Println("Error response from server:", respData1)
				return
			}
			fmt.Printf("Event save completed: %s\n", uneventName)
			Rand_Sleep()
		}
	}
	
	writeJsonFile(dbPath, jsonData)
}