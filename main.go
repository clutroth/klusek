package main

/*
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
#include <X11/extensions/XInput.h>
// XEventType returns a XEvent's type field. XEvent is union. Unions are not supported
XEventClass XEventType(XEvent *event) {
	return event->type;
}
unsigned int keyCode(XEvent *event) {
	return event->xkey.keycode;
}
unsigned int state(XEvent *event) {
	return event->xkey.state;
}
*/
import "C"
import (
	"context"
	"fmt"
	"log"
	"time"
)

type Event interface {
	name() string
}
type KeyPressed struct {
	keycode     uint
	pressedCtrl bool
}

func (k KeyPressed) name() string { return "KeyPressed" }

type KeyReleased struct {
	keycode     uint
	pressedCtrl bool
}

func (k KeyReleased) name() string { return "KeyReleased" }

type OtherEvent struct{}

func (k OtherEvent) name() string { return "OtherEvent" }

func main() {
	ctrlC := KeyReleased{keycode: 54, pressedCtrl: true}

	display := C.XOpenDisplay(nil)
	if display == nil {
		log.Fatalln(fmt.Errorf("unable to open display"))
	}
	defer C.XCloseDisplay(display)

	screen := C.XDefaultScreen(display)
	rootWin := C.XRootWindow(display, screen)
	if int(rootWin) == 0 {
		log.Fatalln(fmt.Errorf("unable to open rootWin"))
	}

	C.XGrabKeyboard(display, rootWin, C.False, C.GrabModeAsync, C.GrabModeAsync, C.CurrentTime)
	defer C.XUngrabKeyboard(display, C.CurrentTime)

	eventChan := make(chan Event, 128)
	defer close(eventChan)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go fetchEvents(ctx, display, eventChan)
	runLoop(eventChan, ctrlC)
	log.Printf("bye bye")
}

func runLoop(events chan Event, breakEvent Event) {
	for {
		event := <-events
		log.Printf("event read %s %+v\n", event.name(), event)
		if event == breakEvent {
			log.Printf("breaking loop by %s %+v\n", event.name(), event)
			return
		}
	}
}

func fetchEvents(ctx context.Context, display *C.Display, eventChan chan Event) {
	for {
		select {
		case <-ctx.Done(): // if cancel() execute
			return
		default:
			if eventQueueSize(display) > 0 {
				eventChan <- readKeyEventBlocking(display)
			} else {
				time.Sleep(300 * time.Millisecond)
			}

		}
	}
}

func readKeyEventBlocking(display *C.Display) Event {
	var xEvent C.XEvent

	C.XNextEvent(display, &xEvent)
	switch C.XEventType(&xEvent) {
	case C.KeyPress:
		return KeyPressed{keycode: keyCode(&xEvent), pressedCtrl: isCtrl(&xEvent)}
	case C.KeyRelease:
		return KeyReleased{keycode: keyCode(&xEvent), pressedCtrl: isCtrl(&xEvent)}
	default:
		return OtherEvent{}
	}
}

func keyCode(xEvent *C.XEvent) uint {
	return uint(C.keyCode(xEvent))
}
func isCtrl(xEvent *C.XEvent) bool {
	return (state(xEvent) & uint(C.ControlMask)) > 0
}
func state(xEvent *C.XEvent) uint {
	return uint(C.state(xEvent))
}
func eventQueueSize(display *C.Display) int {
	return int(C.XEventsQueued(display, C.QueuedAfterReading))
}
