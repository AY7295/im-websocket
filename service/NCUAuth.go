package service

import (
	"errors"
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
	bodyStr := string(body)

	user := &model.User{
		Id:   gjson.Get(bodyStr, "base_info.xh").String(),
		Name: gjson.Get(bodyStr, "base_info.xm").String(),
	}

	if user.Name == "" || user.Id == "" {
		return nil, errors.New(bodyStr)
	}
	return user, nil
}
