package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	// "LearnServer/conf"
	"LearnServer/conf"
)

// PostAndReturn 向另外一个服务器Post并且直接将那台服务器返回的结果返回到前端
func PostAndReturn(path string, data interface{}, CSTREAM func(int, string, io.Reader) error) error {
	client := &http.Client{}

	url := conf.AppConfig.ContentServer + path

	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	return CSTREAM(statusCode, "text/json", resp.Body)
}

// PostAndGetData 向另一台服务器发送Post，得到返回的数据存在respData中。
// 需要给respData传指针，确保respData被Post函数修改
func PostAndGetData(path string, data interface{}, respData interface{}) (int, error) {
	client := &http.Client{}

	url := conf.AppConfig.ContentServer + path

	jsonData, err := json.Marshal(data)

	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	json.Unmarshal([]byte(body), &respData)
	statusCode := resp.StatusCode
	return statusCode, nil
}

// GetAndReturn 向另外一个服务器发送Get并且直接将那台服务器返回的结果返回到前端
func GetAndReturn(path string, paramStr string, CSTREAM func(int, string, io.Reader) error) error {
	url := conf.AppConfig.ContentServer + path + "?" + paramStr

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode

	return CSTREAM(statusCode, "text/json", resp.Body)
}

// GetAndGetData 向另一台服务器发送Get，得到返回的数据存在respData中。
func GetAndGetData(path string, paramStr string, respData interface{}) (int, error) {
	url := conf.AppConfig.ContentServer + path + "?" + paramStr

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	json.Unmarshal([]byte(body), &respData)
	statusCode := resp.StatusCode
	return statusCode, nil
}
