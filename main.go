package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"time"

	"github.com/pwaller/barrier"
)

var once sync.Once

const (
	EOT   = 0x4
	ESC   = 0x1b
	UP    = 'A'
	DOWN  = 'B'
	RIGHT = 'C'
	LEFT  = 'D'
)

func main() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	r := bufio.NewReader(os.Stdin)

	quit := barrier.Barrier{}

	cc := make(chan rune)
	// esc := make(chan rune)
	errc := make(chan error)
	go func() {
		ru, _, err := r.ReadRune()
		if err != nil {
			errc <- err
		} else {
			cc <- ru
		}
	}()

	go func() {
		note := time.Now()
		t := 10 * time.Microsecond
		for {
			start := time.Now()
			for time.Since(start) < 300*time.Microsecond {
			}
			time.Sleep(t)
			if time.Since(note) > 250*time.Millisecond {
				note = time.Now()
				t *= 2
				if t > 1*time.Millisecond {
					t = time.Microsecond
				}
			}
		}
	}()

	go func() {
		for c := range cc {
			switch c {
			case EOT, 'q':
				quit.Fall()
				return
			case ESC:
				if c = <-cc; c != '[' {
					r.UnreadRune()
					continue
				}
				switch <-cc {
				case UP:
					log.Println("up")

				case DOWN:
					log.Println("down")
				}
				continue
			}
			log.Printf("c = %c", c)
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	select {
	case <-sig:
	case <-quit.Barrier():
	}
}
