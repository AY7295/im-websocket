package service

import (
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
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
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "passport "+token)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &model.User{
		Id:   gjson.Get(string(body), "base_info.xh").Str,
		Name: gjson.Get(string(body), "base_info.xm").Str,
	}, nil

}
