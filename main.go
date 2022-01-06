package main

import (
	"goAudio/audio"
	"time"
	"fmt"
)

func printTime(){
	fmt.Println("Hours :" , time.Now().Hour() , ", Minutes :" , time.Now().Minute() , 
	", Seconds :" , time.Now().Second());
}

func main(){
	fileName := "test"
	printTime();	
	audio.Record(fileName , 5);
	printTime();
	audio.Play(fileName + ".aiff");
}