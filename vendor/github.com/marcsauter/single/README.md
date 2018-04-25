# single

`single` provides a mechanism to ensure, that only one instance of a program is running.

    package main

    import (
        "log"
        "time"

        "github.com/marcsauter/single"
    )

    func main() {
        s := single.New("name")
        s.Lock()
        defer s.Unlock()
        log.Println("working")
        time.Sleep(60 * time.Second)
        log.Println("finished")
    }

The package currently supports `linux`, `solaris` and `windows`.
