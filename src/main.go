package main

import (
	"fftester/ffmpeg"
	"fmt"
	"log"
	"runtime"
	"time"
)

func main() {
	log.Println("START")

	log.Println("runtime.GOMAXPROCS:", runtime.GOMAXPROCS(0))

	ffService, err := ffmpeg.NewService()
	if err != nil {
		log.Fatal(err)
	}

	pass1 := `-i test.mp4 -movflags +faststart -y -vcodec libx264 -filter:v scale=-2:720 -b:v 1M -acodec aac -pix_fmt yuv420p -preset slow -crf 22 -maxrate 1M -bufsize 2M -pass 1 -f mp4 /dev/null`
	pass2 := `-i test.mp4 -movflags +faststart -y -vcodec libx264 -filter:v scale=-2:720 -b:v 1M -acodec aac -pix_fmt yuv420p -preset slow -crf 22 -maxrate 1M -bufsize 2M -pass 2 -f mp4 out.mp4`

	pass1log := ""
	pass2log := ""
	start := time.Now()
	if _, _, pass1log, pass2log, err = ffService.Convert(1, "test.mp4", "out.mp4", pass1, pass2); err != nil {
		log.Println("pass1 log:", pass1log)
		log.Println("pass2 log:", pass2log)
		log.Println("error converting :(")
		return
	}
	fmt.Println()
	log.Printf("Test file (32sec) 2-pass converted in %s\n", time.Since(start))
	score := 32 / time.Since(start).Seconds()

	log.Printf("Speed: %.2fx on %d cores\n", score, runtime.GOMAXPROCS(0))

}
