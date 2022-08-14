package wathcer

import (
	fs "github.com/fsnotify/fsnotify"
	"log"
)

var watcherFiles = [2]string{"config/group.toml", "config/mapping.toml"}

func WatchFileChanged(callback func()) {
	// Create new watcher.
	watcher, err := fs.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer func(watcher *fs.Watcher) {
			err := watcher.Close()
			if err != nil {
				log.Println("Close watcher failed:", err)
			}
		}(watcher)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Println("Event watch failed!")
					return
				}
				if has(event.Op, fs.Write) {
					log.Printf("Watched  event changed: %v", event)
					callback()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Println("Event watch err!")
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	for _, p := range watcherFiles {
		err = watcher.Add(p)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func has(now fs.Op, want fs.Op) bool { return now&want == want }
