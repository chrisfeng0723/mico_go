/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2020/12/14 13:46
 */

package main

import (
	"fmt"
	"time"
)

func setVTO1(v *int){
	*v = 1
}

func setVTO2(v *int){
	*v = 2
}

func main(){
	v := new(int)
	go setVTO1(v)
	go setVTO2(v)
	time.Sleep(time.Second)
	fmt.Println(*v)
}
