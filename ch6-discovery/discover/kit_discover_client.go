/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2020/12/18 17:45
 */

package discover

import (
	"github.com/hashicorp/consul/api/watch"
	"log"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"strconv"
	"sync"
)

type KitDiscoverClient struct {
	Host string
	Port int
	client consul.Client
	config *api.Config
	mutex sync.Mutex
	instanceMap sync.Map
}

func NewkitDiscoverClient(consulHost string,consulPort int)(DiscoveryClient,error){
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost+":"+strconv.Itoa(consulPort)
	apiClient,err := api.NewClient(consulConfig)
	if err !=nil{
		return nil,err
	}
	client := consul.NewClient(apiClient)
	return &KitDiscoverClient{
		Host:consulHost,
		Port:consulPort,
		config:consulConfig,
		client:client,
	},nil
}

func(consulClient *KitDiscoverClient) Register(serviceName,instanceId,healthCheckUrl string,instanceHost string, instancePort int,meta map[string]string,logger *log.Logger) bool{
	serviceRegitration :=&api.AgentServiceRegistration{

		ID:                instanceId,
		Name:              serviceName,
		Port:              instancePort,
		Address:           instanceHost,
		Meta:              meta,

		Check:             &api.AgentServiceCheck{
			Interval:                       "15s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,

			DeregisterCriticalServiceAfter: "30s",
		},
	}

	err := consulClient.client.Register(serviceRegitration)
	if err != nil{
		log.Println("Register service Error!")
		return false
	}
	log.Println("Register service Success!")
	return true
}

func (consulClient *KitDiscoverClient)DeRegister(instanceId string,logger *log.Logger) bool{
	serviceRegistration := &api.AgentServiceRegistration{
		ID:instanceId,
	}

	err := consulClient.client.Deregister(serviceRegistration)
	if err !=nil{
		logger.Println("Deregister Service Error!")
		return false
	}
	log.Println("Deregister Service Success!")
	return true
}


func(consulClient *KitDiscoverClient) DiscoverServices(serviceName string,logger *log.Logger)[]interface{}{
	instanceList,ok := consulClient.instanceMap.Load(serviceName)
	if ok {
		return instanceList.([]interface{})
	}
	consulClient.mutex.Lock()
	defer consulClient.mutex.Unlock()
	instanceList,ok = consulClient.instanceMap.Load(serviceName)
	if ok {
		return instanceList.([]interface{})
	}else{
		go func() {
			params := make(map[string]interface{})
			params["type"] = "service"
			params["service"] = serviceName
			plan,_:= watch.Parse(params)
			plan.Handler = func(u uint64, i interface{}) {
				if i ==nil{
					return
				}
				v,ok :=i.([]*api.ServiceEntry)
				if !ok{
					return
				}

				if len(v) == 0{
					consulClient.instanceMap.Store(serviceName,[]interface{}{})

				}
				var healthServices []interface{}
				for _,service := range v{
					if service.Checks.AggregatedStatus() == api.HealthPassing{
						healthServices = append(healthServices,service.Service)
					}
				}
				consulClient.instanceMap.Store(serviceName,healthServices)
			}
			defer plan.Stop()
			plan.Run(consulClient.config.Address)
		}()
	}
	entries,_,err := consulClient.client.Service(serviceName,"",false,nil)
	if err !=nil{
		consulClient.instanceMap.Store(serviceName,[]interface{}{})
		logger.Println("Discover Service Error")
		return nil
	}
	instances := make([]interface{},len(entries))
	for i :=0;i<len(entries);i++{
		instances[i] = entries[i].Service
	}
	consulClient.instanceMap.Store(serviceName,instances)
	return instances
}
