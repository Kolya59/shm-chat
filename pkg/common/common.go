package common

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"

	"github.com/ghetzel/shmtool/shm"
)

const BufSize = 1024

func ReadMsg(segment *shm.Segment, done chan interface{}) {
	for {
		select {
		case <-done:
			return
		default:
			reader := bufio.NewReader(segment)
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				time.Sleep(time.Second)
				continue
			}
			if err != nil {
				log.Println("Failed to read from server: ", err)
			}

			log.Print("Got msg from server: ", line)
			segment.Reset()
			if line == "ESC\n" {
				log.Print("Server stop messaging")
				return
			}
		}
	}
}

func WriteMsg(segment *shm.Segment, done chan interface{}) {
	tmpReader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-done:
			return
		default:
			line, err := tmpReader.ReadString('\n')
			if err != nil {
				log.Print("Failed to read msg from stdin: ", err)
				<-done
				continue
			}

			if len(line) > BufSize {
				log.Println("Invalid msg size")
				continue
			}

			if line == "ESC\n" {
				// TODO Refactor
				log.Print("Close chat")
				<-done
				continue
			}

			if len(line) == 0 {
				log.Println("Empty msg")
				continue
			}

			if _, err := segment.Write([]byte(line)); err != nil {
				log.Println("Failed to write msg to chat: ", err)
				<-done
				continue
			}
		}
	}
}
