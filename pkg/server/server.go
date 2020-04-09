package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ghetzel/shmtool/shm"

	"github.com/kolya59/shm-chat/pkg/common"
)

func StartServer() {
	readerSegment, err := shm.Create(common.BufSize)
	if err != nil {
		log.Println("Failed to create shm: ", err)
		return
	}
	defer func() {
		if err := readerSegment.Destroy(); err != nil {
			log.Print("Failed to destroy shm ", err)
		}
	}()

	log.Printf("Created reader segment with id: %d", readerSegment.Id)

	writerSegment, err := shm.Create(common.BufSize)
	if err != nil {
		log.Println("Failed to create shm: ", err)
		return
	}
	defer func() {
		if err := writerSegment.Destroy(); err != nil {
			log.Print("Failed to destroy shm ", err)
		}
	}()

	log.Printf("Created writer segment with id: %d", writerSegment.Id)

	readerSegmentAddress, err := readerSegment.Attach()
	if err != nil {
		log.Println("Failed to attach reader to shm")
	}
	defer func() {
		if err := readerSegment.Detach(readerSegmentAddress); err != nil {
			log.Println("Failed to detach readerSegment", err)
		}
	}()

	writerSegmentAddress, err := readerSegment.Attach()
	if err != nil {
		log.Println("Failed to attach writer to shm")
	}
	defer func() {
		if err := readerSegment.Detach(writerSegmentAddress); err != nil {
			log.Println("Failed to detach writerSegment", err)
		}
	}()

	done := make(chan interface{})
	go func() {
		tmp := make(chan os.Signal)
		signal.Notify(tmp, syscall.SIGINT, syscall.SIGTERM)
		<-tmp
		close(done)
	}()

	go common.ReadMsg(readerSegment, done)
	go common.WriteMsg(writerSegment, done)

	<-done
}
