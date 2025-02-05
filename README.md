# golang-ipc
Golang Inter-process communication library for Window, Mac and Linux.

:fork_and_knife: This is a fork of [james-barrow/golang-ipc](https://github.com/james-barrow/golang-ipc) with some extra TCP sprinkled in.


### Overview

A simple to use package that uses unix sockets on MacOS/Linux and named pipes on Windows or TCP endpoint for all OS's to create a communication channel between two go processes.


## Usage

Create a server with the default configuration and start listening for the client:

```go

    sc, err := ipc.StartServer("<name of socket or pipe>", nil)
    if err != nil {
        log.Println(err)
        return
    }

```

Create a client and connect to the server:

```go

    cc, err := ipc.StartClient("<name of socket or pipe>", nil)
    if err != nil {
        log.Println(err)
        return
    }

```
Read and write data to the connection:

```go

    // Write data
    _ = sc.Write(1, []byte("Message from server"))

    _ = cc.Write(5, []byte("Message from client"))


    // Read data
    for {
        dataType, data, err := sc.Read()

        if err == nil {
            log.Println("Server recieved: "+string(data)+" - Message type: ", dataType)
        } else {
            log.Println(err)
            break
        }
    }

    for {
        dataType, data, err := cc.Read()

        if err == nil {
            log.Println("Client recieved: "+string(data)+" - Message type: ", dataType)
        } else {
            log.Println(err)
            break
        }
    }

```

### Encryption

By default the connection established will be encrypted, ECDH384 is used for the key exchange and AES 256 GCM is used for the cipher.

Encryption can be switched off by passing in a custom configuration to the server & client start functions.

```go

    config := &ipc.ServerConfig{Encryption: false}
    sc, err := ipc.StartServer("<name of socket or pipe>", config)

```

### Unix Socket Permissions

Under most configurations, a socket created by a user will by default not be writable by another user, making it impossible for the client and server to communicate if being run by separate users.

The permission mask can be dropped during socket creation by passing custom configuration to the server start function.  **This will make the socket writable by any user.**

```go

    config := &ipc.ServerConfig{UnmaskPermissions: true}
    sc, err := ipc.StartServer("<name of socket or pipe>", config)

```
Note: Tested on Linux, not tested on Mac, not implemented on Windows.

### Network mode

By default the connection will use socket/named pipes, but when network mode is enabled connection will use TCP.

It's possible to set the port used by passing in a custom configuration to the server & client start functions.

Both server and client must specify `Network: true` to enable network mode. If only `NetworkListen` or `NetworkServer` is specified they are ignored and socket/named pipe is used.

Server should specify `NetworkListen`, which is the interface to listen on. The value `0.0.0.0` binds to all interfaces.

Client should specify `NetworkServer`, which is the IP address of the server to connect to.

```go

    config := &ipc.ServerConfig{Network: true, NetworkPort: 9494}
    sc, err := ipc.StartServer("<name of server>", config)

```


### Testing

The package has been tested on Mac, Windows and Linux.


### Licence

MIT
