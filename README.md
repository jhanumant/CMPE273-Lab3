# CMPE273-Lab3
The purpose of this lab is to understand how consistent hashing works and to implement a simple RESTful key-value cache data store.  

## Usage
Build REST endpoints to store and retrieve key-Value data.  
This application uses consistent hashing at client-side to hash hostnames and keys.  
### Install

If the package ("github.com/julienschmidt/httprouter") is not installed, go to the code file's folder and do  
```
go get 
```
```
go build

```
Start the  server:
```
go run Server.go
```
####Execute the Client
Run the client:
```
go run Client.go
```
```
Client Output:
1.PUT values  
2.GET a keyValue pair  
3.GET all KeyValues  
4.GET port specific keyValues  
5.Show Node Count  
6.Show Mapping  
7.Remove a Node  
Enter your choice:  
```
Enter a choice(1-7) to perform the realted actions.    
```
1. To PUT a value to any server instance.    
2. GET a Key-Value pair from the server (Requires to input a key).  
3. GET all key-Value pairs from all server instances.  
4. GET port specific key-Value (requires port no.{3000/3001/3002}).  
5. Get the node count on each server (no of key-Value pairs on each server).  
6. Get the server-key mapping information.   
7. Remove a node temporarily (to know how a consistent hashing actually works).  
 ----After this step please execute steps 6/5 to see the re-mapping of nodes.  
```
####Comments
Run the Server.go file first using "go run Server.go" (without double quotes) and in separate command prompt run the Client.go using "go run Client.go"  
There are three server instances(http://localhost:3000,http://localhost:3001,http://localhost:3002).  
Client is already aware of the instances and has added the instance nodes to consistent hashing.
To GET port specific keyValue pairs, you can also place separate cURL requests to the Server.  
e.g  
```
Request:
curl http://localhost:3000/keys
```
```
Response:
[{"key":1,"value":"a"},{"key":3,"value":"c"},{"key":4,"value":"d"},{"key":5,"value":"e"},{"key":7,"value":"g"},{"key":9,"value":"i"},{"key":10,"value":"j"}]
```  
Similarly, requests can be place for remaining ports i.e. 3001 & 3002.  
Two Keys can be same. 
####Refer to OutputImages folder to understand complete execution of the program.
