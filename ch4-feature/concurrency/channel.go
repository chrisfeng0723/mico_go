/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2020/12/14 13:55
 */

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func printInput(ch chan string){
	for val:= range ch{
		if val == "EOF"{
			break
		}
		fmt.Printf("Input is %s\n",val)
	}
}

func consume(ch chan int){
	time.Sleep(time.Second*100)
	<-ch
}
func main(){
	ch := make(chan string)
	go printInput(ch)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan(){
		val := scanner.Text()
		ch <-val
		if val == "EOF"{
			fmt.Println("End the game!")
			break
		}
	}

	defer close(ch)

}

func send(ch chan int,begin int){
	for i :=begin;i <begin+10;i++{
		ch <- i
	}
}
