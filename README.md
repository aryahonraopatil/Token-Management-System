# Token Management System

## Introduction
This is a client-server system for a token management system. In this project, Remote Procedure Call (RPC) functionality is implemented across different operating systems. The system, triggered by specific flags in the command line, reads data from a YAML file and spawns client-side processes to perform operations such as creating, writing, reading, or dropping tokens. These operations are executed through RPC calls, allowing for flexible communication between client and server processes. Additionally, the system incorporates fail-silent behavior by setting conditions for specific tokens on designated servers, simulating failure scenarios for robustness testing. 

## System Implementation

The execution begins at the client.go file. But the major chunk of the code and system logic exists in the server.go file. 

When the client program comes across a Create flag for a token, it reads the yaml_final.yml file and retrieves the data for that token. After retrieving the data, the client spawns the token writer, who then becomes the receiver of a following Create RPC call. After receiving the Create RPC Call, a token is generated if it is not already present. It then becomes responsible for creating the reader servers and sending Create RPC calls to them in order to replicate the token. The RPC calls are adjusted to be adaptable in identifying the origin of the call, enabling servers to decide on accessing data from the YAML file.

The Write RPC call operates based on similar concepts as the Create RPC method. When the Write flag is encountered, the client program will read the yaml_final.yml file and retrieve the data of the token that matches. After retrieving the data, the token writer receives the Write RPC call. After receiving the Write RPC call, the writer compares the timestamp on the token to the timestamp of the incoming call. The write operation is executed if the incoming call has a newer timestamp. The state information of the token is revised, including the updated timestamp. Now, it is responsible for updating the token's state information by sending the relevant readers the Write RPC call to inform them. Similar to Create, the Write RPC call has the flexibility to identify the origin of the call, distinguishing between the writer and the reader within the same function.

In the  Read RPC call, the client program reads the yaml.yml file and fetches the corresponding token's data. It gets the list of readers and chooses a reader randomly to send the Read RPC call to. Once the corresponding reader receives the RPC call, it checks whether the token exists or not. After this check is passed, the Read-Impose Write-Majority comes into play, where all the nodes with the token are contacted by dispatching RPC calls as goroutines. Here the RIWMTest RPC Call is utilized. Upon receiving the call, the server verifies the token, retrieves the requested information, and returns the timestamp and final value. The reader initiating this RPC call starts gathering the results from the goroutines. During the collection process, the timestamps are compared to determine the most recent value. Nevertheless, the reader simply waits for most of the calls to be completed and returned before providing the client with the updated information and writing back to the remaining nodes.

Similarly to the Create and Write RPC calls, in the Drop RPC call, the client program retrieves the token's details from the yaml_final.yml file and acquires the writer. Next, the writer receives a Drop RPC call that must be carried out. If the Drop RPC call occurred after the most recent operation on the token, the token will be removed from the server.
Next, the writer sends Drop RPC requests to each reader to guarantee the removal of the duplicated token. This is achieved by executing the Drop RPC calls as goroutines, making sure uninterrupted deletion and minimizing client wait time.

The requirement of emulating fail-silent behaviour is implemented by hardcoding values to check for the conditions for a specific token in a specific server. In this project, the fail-silent behaviour is emulated for the token with the ID: 1000 in the server running on the port: 5000. The server running on this port will fail to respond to queries on the token 1000 after a time duration of 10 seconds since the token's creation.

## File descriptions

### 1. tokens.proto
The file proto/tokens.proto contains the definition of the messages and lists the RPC calls of the service. This file is automatically generated and hence should not be tampered with. Also, changing the proto file might require changing the client and server files. Whenever you make any changes to this file, recompile the proto file, which can be done with the following command:

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative .\proto\tokens.proto

### 2. yaml.yml
This file contains the replication schemes of all the tokens that are present for this project. This file is accessible by both the client and the server programs.

### 3. logs files
This folder contains the logs of each server. Each file corresponds to the log of a server. The file log_5000.log contains the results of fail-silent emulations. 

### Other files
go.mod and go.sum are files that are automatically generated during the compilation of the project. Similar to the proto file, do not tamper with this file. 

## Project Folder Structure
The structure of the project folder is as follows: 

-client.go  

-server.go  

-proto    
  -tokens_grpc.pb.go  
  -tokens.pb.gp  
  -tokens.proto  

-logs  

-yaml.yml  
-go.mod  
-go.sum

## Commands to run the code
client  -create -id 1500  
client  -write -id 1500 -name token1 -low 0 -mid 10 -high 100  
client  -read -id 1500  
client  -drop -id 1500  
