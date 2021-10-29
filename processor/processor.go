package processor


import (
  "fmt"
  "net/http"
  "encoding/json"
  "bytes"
   "io/ioutil"
)

func Process(msg string) (string, error) {
  fmt.Println("processor.Process: Processing ...")
  
   //Encode the data
   postBody, _ := json.Marshal(map[string]string{"data":  msg })
   responseBody := bytes.NewBuffer(postBody)
   
    //Leverage Go's HTTP Post function to make request
   resp, err := http.Post("http://localhost:7500/processor", "application/json", responseBody)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   
    //Read the response body
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }
   
   //convert response to string format
   sb := string(body)
   
  return sb, nil
}




