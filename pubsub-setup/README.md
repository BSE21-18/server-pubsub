# Usage
Please follow the instructions below to use the listener.

## Files
There are 4 binary files, 2 for windows and 2 for linux.
1. processor-win.exe (processor mockup on windows)
2. processor (processor mockup on linux)
3. pubsub-listener-win.exe (the listener on windows)
4. pubsub-listener (the listener on linux)
Execute both files for your operating system and use the following endpoints to make requests

## PORT: 7000
The listener port
## PORT: 7500
The processor mock-up port. 
(You do not need to call this one at all. The listener calls it internally)

## Endpoints
1. To register a new device:
```javascript
var endpoint = "http://localhost:7000/subscribe";
```
2. To publish data after detecting using the sniffer :
```javascript
var endpoint = "ws://localhost:7000/publish"; 
//method: POST
```
NOTE: The expected data which is to be sent via a websocket message should be a stringified json of the following format
```javascript
var mockData = `{"gas1":"345PH","gas2":"345PH","gas3":"345PH","gas4":"345PH","gas5":"345PH","gas6":"345PH"}`;
var requestData = {"topic":activeDeviceLabel, "message": mockData };  
```
3. To listen and get updates from a device or all your registered devices:
```javascript
var endpoint = "ws://localhost:7000/getupdates"; 
```
NOTE: The expected data which is to be sent via a websocket message should be a stringified json of the following format
```javascript
var requestData = `{"device": "DV00023", "phone": "+256706123303"}` ;
```

or, to receive updates from all devices registered by the current user:
```javascript
var requestData = `{"device": "All", "phone": "+256706123303"}`;
```

## Database config
```go
    dbUser := "root"
	dbPswd := ""
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "datavoc"
```

## code
The code is available on github:
https://github.com/BSE21-18/server-pubsub
https://github.com/BSE21-18/server-processor






