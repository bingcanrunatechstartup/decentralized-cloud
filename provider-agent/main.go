package main

import (
	"context"
    "os"
    "fmt"
    "os/exec"
    "encoding/json"
    zmq "github.com/pebbe/zmq4"
	"github.com/hashicorp/terraform-exec/tfexec"
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
        json.Unmarshal(msgBytes, &msgData)

        // execute terraform script with provided data
        //cmd := exec.Command("terraform", "apply", "-var", fmt.Sprintf("vcpus=%d", msgData.VCPUs), "-var", fmt.Sprintf("memory=%d", msgData.Memory))
        //cmd.Run()
		tf, err := tfexec.NewTerraform("<TerraformDir>", exec.LookPath("terraform"))
		if err != nil {
			// handle error
		}
	
		err = tf.Init(context.Background(), tfexec.Upgrade(true))
		if err != nil {
			// handle error
		}
	
		err = tf.Apply(context.Background())
		if err != nil {
			// handle error
		}
		output,err := tf.Output(context.Background(),"ip_address")
		if err != nil{
			//handle error 
		}
		fmt.Println(output.Value)

        // execute ansible script against provisioned VM
        cmd = exec.Command("ansible-playbook", "playbook.yml")
        cmd.Run()

        // return relevant information such as IP address by placing it on ZeroMQ queue
        respData := MessageData{VCPUs: 2, Memory: 2048}
        respBytes, _ := json.Marshal(respData)
        sender.SendBytes(respBytes, 0)
    }
}