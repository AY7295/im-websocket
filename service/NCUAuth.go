package service

import (
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"webSocket-be/model"
)

var (
	url    = "https://os.ncuos.com/api/user/profile/basic"
	method = "GET"
)

func GetUserByToken(token string) (*model.User, error) {

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "passport "+token)

	res, err := (&http.Client{}).Do(req)
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

	return &model.User{
		Id:   gjson.Get(string(body), "base_info.xh").Str,
		Name: gjson.Get(string(body), "base_info.xm").Str,
	}, nil

}
