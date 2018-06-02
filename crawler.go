package main

import (
	"fmt"
	"webcrawler/workerpool"
)

type crawlerJob struct {
}

func lol() {
	workerpool := workerpool.New(5)
	fmt.Println(workerpool)
}
