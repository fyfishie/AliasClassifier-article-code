/*
  - @Author: fyfishie
  - @Date: 2023-03-27:17

* @LastEditors: fyfishie

* @LastEditTime: 2023-05-13:16
  - @Description: check http request format
  - @email: fyfishie@outlook.com
*/
package main

import (
	"aliasParseMaster/lib"
	"aliasParseMaster/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// func streamReqParse(w http.ResponseWriter, r *http.Request) (lib.UserInput, error) {
// 	userInput := lib.UserInput{}
// 	bs, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		return userInput, err
// 	}
// 	err = json.Unmarshal(bs, &userInput)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write(utils.MakeMessageBytes(false, err.Error()))
// 		return userInput, err
// 	}
// 	if len(userInput.IPToDoList) == 0 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write(utils.MakeMessageBytes(false, "ip set should not be in zero length"))
// 		return userInput, errors.New("zero length of list input")
// 	}
// 	return userInput, nil
// }

// parse request params, if any error occurs, automatically responses to client and returns err
// if all runs well, returns err(==nil)
func parseVPRegist(w http.ResponseWriter, r *http.Request) (vpURL, vpName string, err error) {
	err = r.ParseForm()
	if err == nil {
		urls, ok := r.Form["slave_url"]
		if ok {
			if len(urls) > 0 {
				names, ok := r.Form["name"]
				if ok {
					if len(names) > 0 {
						return urls[0], names[0], nil
					}
					err = errors.New("more than zero name args expected")
				}
				err = errors.New("can not parse slave name")
			}
			err = errors.New("more than zero url args expected")
		}
		err = errors.New("can not parse slave_url")
	}
	return "", "", err
}

type UserInput struct {
	TaskName string   `json:"task_name"`
	IPList   []string `json:"ip_list"`
}

// TODO:valid check
func StreamWorkReqCheck(w http.ResponseWriter, r *http.Request) (userInput *lib.UserInput, err error) {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write(utils.MakeMessageBytes(false, "can not read request body, connect to admin please"))
		w.WriteHeader(http.StatusInternalServerError)
		return nil, nil
	}
	tmp := lib.UserInput{}
	err = json.Unmarshal(bs, &tmp)
	if err != nil {
		w.Write(utils.MakeMessageBytes(false, "can not read request body, connect to admin please"))
		w.WriteHeader(http.StatusInternalServerError)
		return nil, nil
	}
	return &tmp, nil
}
