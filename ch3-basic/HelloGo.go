/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2020/12/11 16:04
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type requestBody struct {
	Key    string `json:"key"`
	Info   string `json:"info"`
	UserId string `json:"user_id"`
}

type responseBody struct {
	Code int      `json:"code"`
	Text string   `json:"text"`
	List []string `json:"list"`
	Url  string   `json:"url"`
}

func process(inputChan <-chan string, userId string) {
	for {
		input := <-inputChan
		if input == "EOF" {
			break
		}
		reqData := &requestBody{
			Key:    "792bcf45156d488c92e9d11da494b085",
			Info:   input,
			UserId: userId,
		}

		byteData,_:= json.Marshal(&reqData)

		req,err := http.NewRequest("POST",
			"http://www.tuling123.com/openapi/api",
				bytes.NewReader(byteData))

		req.Header.Set("Content-Type","application/json;charset=UTF-8")
		client := http.Client{
			Timeout:       2*time.Second,
		}
		resp,err := client.Do(req)
		if err !=nil{
			fmt.Println("Network error!")
			fmt.Println(err)
		}else{
			body,_ := ioutil.ReadAll(resp.Body)
			var resData responseBody
			json.Unmarshal(body,&resData)
			fmt.Println("AI:"+ resData.Text)
		}
		if resp !=nil{
			resp.Body.Close()
		}



	}
}
func main() {
	var input string
	fmt.Println("我是一个机器人，和你聊开心的和不开心的")
	channel := make(chan string)
	defer close(channel)
	go process(channel,string(rand.Int63()))

	for {
		fmt.Scanf("%s",&input)
		channel <-input
		if input == "EOF"{
			fmt.Println("Bye")
			break
		}
	}

}
