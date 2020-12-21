/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2020/12/21 10:50
 */

package discover

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type InstanceInfo struct {
	ID string `json:"id"`
	Service string `json:"service,omitempty"`
	Name string `json:"name"`
	Tags []string `json:"tags,omitempty"`
	Address string `json:"address"`
	Port int `json:"port"`
	Meta map[string]string `json:"meta,omitempty"`
	EnableTagOverride bool `json:"enable_tag_override"`
	Check `json:"check,omitempty"`
	Weights `json:"weigths,omitempty"`
}

type Check struct {
	DeregisterCriticalServiceAfter string `json:"deregister_critical_service_after"`
	Args []string `json:"args,omitempty"`
	HTTP string `json:"http"`
	Interval string `json:"interval,omitempty"`
	TTL string `json:"ttl,omitempty"`
}

type Weights struct {
	Passing int `json:"passing"`
	Warning int `json:"warning"`
}


type MyDiscoverClient struct {
	Host string
	Port int
}

func NewMyDiscoverClient(consulHost string,consulPort int)(DiscoveryClient,error){
	return &MyDiscoverClient{
		Host: consulHost,
		Port: consulPort,
	},nil
}
func (consulClient *MyDiscoverClient) Register(serviceName, instanceId, healthCheckUrl string,instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool {
	instanceInfo := &InstanceInfo{
		ID:                instanceId,
		Name:              serviceName,
		Address:           instanceHost,
		Port:              instancePort,
		Meta:              meta,
		EnableTagOverride: false,
		Check:             Check{
			DeregisterCriticalServiceAfter:"30s",
			HTTP: "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval: "15s",
		},
		Weights:           Weights{
			Passing:10,
			Warning:1,
		},
	}
	byteData,_ := json.Marshal(instanceInfo)
	req,err := http.NewRequest("PUT",
			"http://"+consulClient.Host+":"+strconv.Itoa(consulClient.Port)+"/v1/agent/service/register",
			bytes.NewReader(byteData))
	if err ==nil{
		req.Header.Set("Content-Type","application/json;charset=UTF-8")
		client :=http.Client{}
		resp,err :=client.Do(req)
		if err !=nil{
			log.Println("Register service Error!")
		}else{
			resp.Body.Close()
			if resp.StatusCode == 200{
				log.Println("Register service Success!")
				return true
			}else{
				log.Println("Register service Error!")
			}

		}
	}

	return false
}

func (consulClient *MyDiscoverClient) DeRegister(instanceId string, logger *log.Logger) bool {
	req,err := http.NewRequest("PUT",
		"http://"+consulClient.Host+":"+strconv.Itoa(consulClient.Port)+"/v1/agent/service/deregister"+instanceId,
		nil)
	client := http.Client{}
	resp,err := client.Do(req)
	if err !=nil{
		log.Println("Deregister service Error!")
	}else{
		resp.Body.Close()
		if resp.StatusCode == 200{
			log.Println("Deregister service Success!")
			return true
		}else{
			log.Println("Deregister service Error!")
		}
	}
	return false
}

func (consulClient *MyDiscoverClient) DiscoverServices(serviceName string, logger *log.Logger) []interface{} {
	req,err := http.NewRequest("GET",

		"http://"+consulClient.Host+":"+strconv.Itoa(consulClient.Port)+"/v1/health/service/"+serviceName,nil)
	client := http.Client{}
	resp,err := client.Do(req)
	if err !=nil{
		log.Println("Discover Service Error!")
	}else if resp.StatusCode == 200{
		var servicelist []struct{
			Service InstanceInfo `json:"service"`
		}
		err = json.NewDecoder(resp.Body).Decode(&servicelist)
		resp.Body.Close()
		if err ==nil{
			instances := make([]interface{},len(servicelist))
			for i :=0;i<len(instances);i++{
				instances[i]= servicelist[i].Service
			}
			return instances
		}
	}
	return nil
}
