package client

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ghetzel/shmtool/shm"

	"github.com/kolya59/shm-chat/pkg/common"
)

func StartClient(readerID, writerID int) {
	readerSegment, err := shm.Open(readerID)
	if err != nil {
		log.Println("Failed to reader shm: ", err)
		return
	}
	log.Printf("Opened reader segment with id: %d", readerSegment.Id)

	writerSegment, err := shm.Open(writerID)
	if err != nil {
		log.Println("Failed to open writer shm: ", err)
		return
	}
	log.Printf("Opened writer segment with id: %d", writerSegment.Id)

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
