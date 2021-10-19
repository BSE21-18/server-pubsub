package processor


import (
    "fmt"
)

func Process(msg string) (string, error) {
  fmt.Println("processor.Process: Processing ...")
  
  //TODO: create an http client
  //TODO: call the processor endpoint and pass the received string msg
  //TODO: wait for processed message
  //TODO: if response is an error, 
            //return "", resp
  //TODO: return the processed message
  return msg, nil
}
