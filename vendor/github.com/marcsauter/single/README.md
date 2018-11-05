# single

`single` provides a mechanism to ensure, that only one instance of a program is running.

    package main

    import (
        "log"
        "time"

        "github.com/marcsauter/single"
    )

    func main() {
        s := single.New("your-app-name")
        if err := s.CheckLock(); err != nil && err == single.ErrAlreadyRunning {
            log.Fatal("another instance of the app is already running, exiting")
        } else if err != nil {
            // Another error occurred, might be worth handling it as well
            log.Fatalf("failed to acquire exclusive app lock: %v", err)
        }
        defer s.TryUnlock()

        log.Println("working")
        time.Sleep(60 * time.Second)
        log.Println("finished")
    }

The package currently supports `linux`, `solaris` and `windows`.
