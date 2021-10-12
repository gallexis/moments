package main

import (
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "steelseries/back"
)

func main(){
    server := back.NewApiServer()

    go func() {
        ch := make(chan os.Signal)
        signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGHUP)
        <-ch

        err := server.Close()
        if err != nil {
            fmt.Println("Error when closing server")
        }
    }()

    err := server.Start()
    if err != nil && err != http.ErrServerClosed{
        fmt.Println("Error when starting server", err)
    }
}