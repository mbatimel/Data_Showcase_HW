package main

import (
	"fmt"
	"log"

	"github.com/mbatimel/Data_Showcase_HW/internal/server"
)



func main() {
	serv,err := server.NewRedisConnection(3)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Redis is working.")

	fmt.Println(serv.Cap())

	err = serv.Clear()
	if err != nil {
		log.Fatalln(err)
	}
	serv.Add("key1", "value1")
    value, ok := serv.Get("key1")
    if ok {
        fmt.Println("Получено значение:", value)
    } else {
        fmt.Println("Ключ не найден")
    }
}