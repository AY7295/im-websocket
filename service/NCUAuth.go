package service

import (
	"encoding/json"
	"fmt"
	io "io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"webSocket-be/model"
)

func GetUserByToken(token string) (*model.User, error) {
	url := "https://os.ncuos.com/api/user/profile/basic"
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "passport "+token)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err1 := Body.Close()
		if err1 != nil {

		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("json unmarshal error:", err)
	}

	user := &model.User{
		Id:   data["base_info"].(map[string]interface{})["xh"].(string),
		Name: data["base_info"].(map[string]interface{})["xm"].(string),
	}

	return user, nil
}
