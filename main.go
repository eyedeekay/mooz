package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
	"github.com/atotto/clipboard"
	"github.com/christianhujer/isheadless"
	goi2pbrowser "github.com/eyedeekay/go-i2pbrowser"
	server "github.com/yuukanoo/rtchat/cmd"
)

func launch() bool {
	present := commandExists("firefox")
	log.Println("Firefox presence test indicates: ", present)
	env := os.Getenv("mooz_no_launch_app")
	log.Println("mooz_no_launch_app value", env)
	switch env {
	case "false":
		return false
	case "0":
		return false
	case "f":
		return false
	case "":
		return true
	case "true":
		return true
	case "1":
		return true
	case "t":
		return true
	default:
		if i, err := strconv.Atoi(env); err == nil {
			if i > 0 {
				return true
			}
			if i <= 0 {
				return false
			}
		}
		return commandExists("firefox")
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

var addr = ""
var dir = ""

func main() {
	e := server.Flags{
		Turn: server.TurnFlags{
			RealmString:    flag.String("realm", "mooz.i2p", "Realm used by the turn server."),
			PublicIPString: flag.String("turn-ip", "127.0.0.1", "IP Address that TURN can be contacted on. Should be publicly available."),
			PortInt:        flag.Int("turn-port", 3478, "Listening port for the TURN/STUN endpoint."),
			I2p: server.I2pFlags{
				SamIP:   flag.String("sam-ip", "127.0.0.1", "IP address on which the Simple Anonymous Messaging bridge can be reached"),
				SamPort: flag.Int("sam-port", 7656, "Port on which the Simple Anonymous Messaging bridge can be reached"),
			},
		},
		Web: server.WebFlags{
			Port: flag.Int("http-port", 5000, "Web server listening port."),
		},
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	d := flag.String("dir", wd, "directory to store application state in")
	host := flag.String("upstream", "https://moam7ks26jxodox6orfvcc5ypazvfmluramym6435pfechdd4sdq.b32.i2p", "Third-party WebRTC chat host")
	hosted := flag.Bool("hosted", false, "Use third-party WebRTC chat host")
	launch := flag.Bool("app", launch(), "Start the application")
	tray := flag.Bool("tray", !isheadless.IsHeadless(), "Show the application running in the system tray")
	flag.Parse()
	dir = *d
	if !*hosted {
		go func() {
			addr = server.Serve(e, *e.Turn.RealmString)
			log.Println("Server started")
		}()
		for {
			if addr != "" {
				break
			}
			time.Sleep(time.Second * 2)
		}
		defer server.Close()
		log.Println(addr)
	}
	if *hosted {
		addr = *host
	}

	go func() {
		if *tray {
			systray.Run(onReady, onExit)
		}
	}()
	if *launch {
		go goi2pbrowser.BrowseApp(dir, addr)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Shutting down, goodbye 👋")
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("mooz")
	systray.SetTooltip("Video calls over I2P")
	mBrowse := systray.AddMenuItem("Launch call", "Open a window where a call can happen")
	mCopy := systray.AddMenuItem("Copy URL", "Copy the URL of your Voice Chat")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()
	go func() {
		for {
			select {
			case <-mBrowse.ClickedCh:
				go goi2pbrowser.BrowseApp(dir, addr)
			case <-mCopy.ClickedCh:
				clipboard.WriteAll(addr)
			}
			time.Sleep(time.Second)
		}
	}()
	// Sets the icon of a menu item.
	mQuit.SetIcon(icon.Data)
}

func onExit() {
	// clean up here
	if runtime.GOOS == "windows" {
		syscall.Kill(syscall.Getpid(), syscall.SIGKILL)
	} else {
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}
}
