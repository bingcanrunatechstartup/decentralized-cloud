package main

import (
    "fmt"
    "os/exec"
    zmq "github.com/pebbe/zmq4"
)

type MessageData struct {
    VCPUs int
    Memory int
    // other fields...
}

func main() {
    // setup ZeroMQ socket to receive messages
    receiver, _ := zmq.NewSocket(zmq.PULL)
    defer receiver.Close()
    receiver.Connect("tcp://localhost:5557")

    for {
        // receive message from queue
        msgBytes, _ := receiver.RecvBytes(0)
        var msgData MessageData
        // unmarshal message data into struct
        // ...

        // execute terraform script with provided data
        cmd := exec.Command("terraform", "apply", "-var", fmt.Sprintf("vcpus=%d", msgData.VCPUs), "-var", fmt.Sprintf("memory=%d", msgData.Memory))
        cmd.Run()

        // execute ansible script against provisioned VM
        cmd = exec.Command("ansible-playbook", "playbook.yml")
        cmd.Run()

        // return relevant information such as IP address by placing it on ZeroMQ queue
        // ...
    }
}