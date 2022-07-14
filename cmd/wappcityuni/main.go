package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"golang.captainalm.com/cityuni-webserver/conf"
	"golang.captainalm.com/cityuni-webserver/pageHandler"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	buildVersion = "develop"
	buildDate    = ""
)

func main() {
	log.Printf("[Main] Starting up GO Package Header Server #%s (%s)\n", buildVersion, buildDate)
	y := time.Now()

	//Hold main thread till safe shutdown exit:
	wg := &sync.WaitGroup{}
	wg.Add(1)

	//Get working directory:
	cwdDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	//Load environment file:
	err = godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	//Data directory processing:
	dataDir := os.Getenv("DIR_DATA")
	if dataDir == "" {
		dataDir = path.Join(cwdDir, ".data")
	}

	check(os.MkdirAll(dataDir, 0777))

	//Config file processing:
	configLocation := os.Getenv("CONFIG_FILE")
	if configLocation == "" {
		configLocation = path.Join(dataDir, "config.yml")
	} else {
		if !filepath.IsAbs(configLocation) {
			configLocation = path.Join(dataDir, configLocation)
		}
	}

	//Config loading:
	configFile, err := os.Open(configLocation)
	if err != nil {
		log.Fatalln("Failed to open config.yml")
	}

	var configYml conf.ConfigYaml
	groupsDecoder := yaml.NewDecoder(configFile)
	err = groupsDecoder.Decode(&configYml)
	if err != nil {
		log.Fatalln("Failed to parse config.yml:", err)
	}

	//Server definitions:
	var webServer *http.Server
	var fcgiListen net.Listener
	switch strings.ToLower(configYml.Listen.WebMethod) {
	case "http":
		webServer = &http.Server{
			Handler:           pageHandler.GetRouter(configYml),
			ReadTimeout:       configYml.Listen.ReadTimeout,
			ReadHeaderTimeout: configYml.Listen.WriteTimeout,
		}
		go runBackgroundHttp(webServer, getListener(configYml, cwdDir), false)
	case "fcgi":
		fcgiListen = getListener(configYml, cwdDir)
		if fcgiListen == nil {
			log.Fatalln("Listener Nil")
		} else {
			go runBackgroundFCgi(pageHandler.GetRouter(configYml), fcgiListen)
		}
	default:
		log.Fatalln("Unknown Web Method.")
	}

	//=====================
	// Safe shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//Startup complete:
	z := time.Now().Sub(y)
	log.Printf("[Main] Took '%s' to fully initialize modules\n", z.String())

	go func() {
		<-sigs
		fmt.Printf("\n")

		log.Printf("[Main] Attempting safe shutdown\n")
		a := time.Now()

		if webServer != nil {
			log.Printf("[Main] Shutting down HTTP server...\n")
			err := webServer.Close()
			if err != nil {
				log.Println(err)
			}
		}

		if fcgiListen != nil {
			log.Printf("[Main] Shutting down FCGI server...\n")
			err := fcgiListen.Close()
			if err != nil {
				log.Println(err)
			}
		}

		log.Printf("[Main] Signalling program exit...\n")
		b := time.Now().Sub(a)
		log.Printf("[Main] Took '%s' to fully shutdown modules\n", b.String())
		wg.Done()
	}()
	//
	//=====================
	wg.Wait()
	log.Println("[Main] Goodbye")
	//os.Exit(0)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getListener(config conf.ConfigYaml, cwd string) net.Listener {
	split := strings.Split(strings.ToLower(config.Listen.WebNetwork), ":")
	if len(split) == 0 {
		log.Fatalln("Invalid Web Network")
		return nil
	} else {
		var theListener net.Listener
		var theError error
		log.Println("[Main] Socket Network Type:" + split[0])
		log.Printf("[Main] Starting up %s server on %s...\n", config.Listen.WebMethod, config.Listen.Web)
		switch split[0] {
		case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6", "ip", "ip4", "ip6":
			theListener, theError = net.Listen(strings.ToLower(config.Listen.WebNetwork), config.Listen.Web)
		case "unix", "unixgram", "unixpacket":
			socketPath := config.Listen.Web
			if !filepath.IsAbs(socketPath) {
				if !filepath.IsAbs(cwd) {
					log.Fatalln("Web Path Not Absolute And No Working Directory.")
					return nil
				}
				socketPath = path.Join(cwd, socketPath)
			}
			log.Println("[Main] Removing old socket.")
			if err := os.RemoveAll(socketPath); err != nil {
				log.Fatalln("Could Not Remove Old Socket.")
				return nil
			}
			theListener, theError = net.Listen(strings.ToLower(config.Listen.WebNetwork), config.Listen.Web)
		default:
			log.Fatalln("Unknown Web Network.")
			return nil
		}
		if theError != nil {
			log.Fatalln("Failed to listen due to:", theError)
			return nil
		}
		return theListener
	}
}

func runBackgroundHttp(s *http.Server, l net.Listener, tlsEnabled bool) {
	var err error
	if tlsEnabled {
		err = s.ServeTLS(l, "", "")
	} else {
		err = s.Serve(l)
	}
	if err != nil {
		if err == http.ErrServerClosed {
			log.Println("The http server shutdown successfully")
		} else {
			log.Fatalf("[Http] Error trying to host the http server: %s\n", err.Error())
		}
	}
}

func runBackgroundFCgi(h http.Handler, l net.Listener) {
	err := fcgi.Serve(l, h)
	if err != nil {
		if err == net.ErrClosed {
			log.Println("The fcgi server shutdown successfully")
		} else {
			log.Fatalf("[Http] Error trying to host the fcgi server: %s\n", err.Error())
		}
	}
}
