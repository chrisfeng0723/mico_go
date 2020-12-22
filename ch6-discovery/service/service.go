/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2020/12/22 18:31
 */

package service

import (
	"context"
	"discovery/config"
	"discovery/discover"
	"errors"
)

type Service interface {
	HealthCheck() bool
	SayHello() string
	DiscoveryService(ctx context.Context,serviceName string)([]interface{},error)
}

var ErrNotServiceInstances = errors.New("instances are not existed")

type DiscoveryServiceImpl struct {
	discoveryClient discover.DiscoveryClient
}

func NewDiscoveryServiceImpl(discoveryClient discover.DiscoveryClient) Service{
	return &DiscoveryServiceImpl{discoveryClient:discoveryClient}
}

func(*DiscoveryServiceImpl) SayHello() string{
	return "hello world"
}

func(service *DiscoveryServiceImpl) DiscoveryService(ctx context.Context,serviceName string)([]interface{},error){
	instances := service.discoveryClient.DiscoverServices(serviceName,config.Logger)
	if instances ==nil || len(instances) == 0{
		return nil,ErrNotServiceInstances
	}
	return instances,nil
}

func (*DiscoveryServiceImpl) HealthCheck() bool{
	return true
}