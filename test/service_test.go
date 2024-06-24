package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/mbatimel/Data_Showcase_HW/internal/service"
)

func TestConnect(t *testing.T) {
	_,err := service.NewRedisConnection(3)
	if err != nil {
		t.Errorf("%q",err)
	}
}

func TestCups(t *testing.T) {
	serv1,_ := service.NewRedisConnection(3)
	serv2,_ := service.NewRedisConnection(4)
	serv3,_ := service.NewRedisConnection(5)
	serv4,_ := service.NewRedisConnection(6)

	res1 := serv1.Cap()
	res2 := serv2.Cap()
	res3 := serv3.Cap()
	res4 := serv4.Cap()
	if res1 != 3{
		t.Errorf("got %q, wanted %q", res1, 3)
	}
	if res2 != 4{
		t.Errorf("got %q, wanted %q", res2, 4)
	}
	if res3 != 5{
		t.Errorf("got %q, wanted %q", res3, 5)
	}
	if res4 != 6{
		t.Errorf("got %q, wanted %q", res4, 6)
	}

}

func TestSet(t *testing.T) {
	serv,_ := service.NewRedisConnection(3)
	serv.Add("key1", "value1")
	serv.Add("key2", "value2")
	serv.Add("key3", "value3")
	for i:=1; i<=3; i++ {
		key_i:="key" + strconv.Itoa(i)
		value, ok := serv.Get(key_i)
		if !ok{
			t.Errorf("The value was not found")
		}else{
			fmt.Println(value)
		}	
	}
}

func TestGet(t *testing.T) {
	serv,_ := service.NewRedisConnection(3)
	serv.Add("key1", "value1")
	serv.Add("key2", "value2")
	serv.Add("key3", "value3")
	for i:=1; i<=3; i++ {
		key_i:="key" + strconv.Itoa(i)
		value, ok := serv.Get(key_i)
		if !ok{
			t.Errorf("The value was not found")
		}else{
			fmt.Println(value)
		}	
	}
		value, ok := serv.Get("key4")
		if !ok{
			fmt.Println("The value was not found")
		}else{
			t.Error(value)
		}	

}

func TestRemove(t *testing.T) {
	serv,_ := service.NewRedisConnection(3)
	serv.Add("key1", "value1")
	value, ok := serv.Get("key1")
	if !ok{
		t.Errorf("The value was not found")
	}else{
		fmt.Println(value)
	}
	serv.Remove("key1")
	value, ok = serv.Get("key1")
	if !ok{
		fmt.Println("The value was not found")
	}else{
		t.Errorf("value is %q",value)
	}	
	
}

func TestAddTller(t *testing.T) {
	serv,_ := service.NewRedisConnection(3)
	serv.AddWithTTL("key1", "value1",3)
	value, ok := serv.Get("key1")
	if !ok{
		t.Errorf("The value was not found")
	}else{
		fmt.Println(value)
	}
	time.Sleep(6* time.Second)
	value, ok = serv.Get("key1")
	if !ok{
		fmt.Println("The value was not found")
	}else{
		t.Errorf("value is %q",value)
	}	
}


func TestLRU(t *testing.T) {
		serv, err := service.NewRedisConnection(3)
		if err != nil {
			t.Fatalf("Failed to create cache service: %v", err)
		}
	
		serv.Add("key1", "value1")
		serv.Add("key2", "value2")
		serv.Add("key3", "value3")
		serv.Add("key4", "value4")
	
		if _, ok := serv.Get("key1"); ok {
			t.Error("Expected key1 to be evicted")
		}
		if _, ok := serv.Get("key2"); !ok {
			t.Error("Expected key2 to be present")
		}
		if _, ok := serv.Get("key3"); !ok {
			t.Error("Expected key3 to be present")
		}
		if _, ok := serv.Get("key4"); !ok {
			t.Error("Expected key4 to be present")
		}
}

func TestClear(t *testing.T) {
	serv,_ := service.NewRedisConnection(3)
	err := serv.Clear()
	if err != nil {
		t.Error(err)
	}
}

