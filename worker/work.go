package worker

import (
	"fmt"
	"golang.org/x/sys/unix"
)

func (w *Worker) Work() error {
	defer func() {
		err := unix.EpollCtl(w.EpollFD, unix.EPOLL_CTL_DEL, w.ClientFD, nil)
		if err != nil {
			fmt.Printf("server start: delete client fd from epoll: %v\n", err)
			return
		}

		err = unix.Close(w.ClientFD)
		if err != nil {
			fmt.Printf("server start: close client fd: %v\n", err)
			return
		}

		fmt.Println("work done")
	}()

	buf := make([]byte, 1024)
	_, err := unix.Read(w.ClientFD, buf)
	if err != nil {
		return fmt.Errorf("worker work: read from client fd:%d %v\n", w.ClientFD, err)
	}

	var resp string
	switch string(buf[0]) {
	case "1":
		resp = "first response"
	case "2":
		resp = "second response"
	case "3":
		resp = "third response"
	default:
		resp = "unknown response"
	}

	err = unix.Send(
		w.ClientFD,
		[]byte(resp),
		unix.MSG_DONTWAIT,
	)
	if err != nil {
		return fmt.Errorf("worker work: send message to socket: %v\n", err)
	}

	return nil
}
