package bird

import (
	"bytes"
	"io"
	"net"
	"time"
)

var UnixSocketPath = "/run/bird/bird.ctl"

func Command(cmd string) (string, error) {
	conn, err := net.DialTimeout("unix", UnixSocketPath, 30*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_, err = conn.Write(append([]byte(cmd), []byte("\n")...))
	if err != nil {
		return "", err
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	var buf bytes.Buffer

	ec := make(chan error, 1)
	tm := time.NewTimer(3 * time.Second)

	go func() {
		conn.SetReadDeadline(time.Now().Add(time.Second * 10))

		var tmp [128]byte
		for {
			size, err := conn.Read(tmp[:])
			if err != nil && err != io.EOF {
				tm.Stop()
				ec <- err
				return
			}
			if size == 0 {
				tm.Stop()
				ec <- nil
				return
			}
			buf.Write(tmp[:size])
			tm.Reset(time.Second)
		}
	}()

	select {
	case <-tm.C:
	case err := <-ec:
		if err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}
