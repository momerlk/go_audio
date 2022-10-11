package main

import (
	"github.com/dragonmaster101/go_audio/audio"
	"time"
	"fmt"
)

func printTime(){
	fmt.Println("Hours :" , time.Now().Hour() , ", Minutes :" , time.Now().Minute() , 
	", Seconds :" , time.Now().Second());
}

func main(){
	r := &audio.Recorder{}
	r.Init("./test.wav");
	err := r.Start();
	if err != nil {
		fmt.Println(err);
		return;
	}
	fmt.Println("Press ENTER to stop recording");
	fmt.Scanln();
	fmt.Println("done recording!");
	err = r.Stop();
	if err != nil {
		fmt.Println(err);
		return;
	}
}
