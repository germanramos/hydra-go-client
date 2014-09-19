hydra-go-client
===============
[![Build Status](https://travis-ci.org/innotech/hydra-go-client.svg?branch=master)](https://travis-ci.org/innotech/hydra-go-client) [![Coverage Status](https://coveralls.io/repos/innotech/hydra-go-client/badge.png?branch=master)](https://coveralls.io/r/innotech/hydra-go-client?branch=master) [![GoDoc](https://godoc.org/github.com/innotech/hydra-go-client/client?status.png)](https://godoc.org/github.com/innotech/hydra-go-client/client)

Client of Hydra development in go. Hydra is a multi-cloud broker system.Provides a multi-cloud application discovery, management and balancing service. Hydra attempts to ease the routing and balancing burden from servers and delegate it on the client 

For a complete information about the project visit http://innotech.github.io/hydra/.

##Obtain client

```
    go get github.com/innotech/hydra-go-client
```

##Hydra client basic usage

The basic way to connect to hydra using the GO client is:

```go
    ...
    
    import (
  	  . "github.com/innotech/hydra-reverse-proxy/vendors/github.com/innotech/hydra-go-client/client"
  	)
  
    ...

    if err := HydraClientFactory.Config([]string{"http://localhost:7772"}); err != nil {
  		log.Fatal(err.Error())
  	}
    hydraClient := HydraClientFactory.Build()
    serverURLs, err := hydraClient.Get(h.AppId, false)
  
    //Some network call using the first of the candidate servers.

```

The previous code fragment configure the client to search hydra server in localhost.

In this case the rest of config parameters are configured using the following default values.

###Configuration parameters
Name | Default value | Description 
:---  | :--- | :---
appsCacheDuration | 60 seconds | The time period that the cache that store the candidate servers for applications is invalidated.
hydraServersCacheDuration| 20 seconds | The time period that the cache that store the hydra servers is invalidated.
maxNumberOfRetries| 10 | The client try this number of times to connect to all the register hydra servers.
durationBetweenAllServersRetry| 0 milliseconds | The time between all hydra servers are tried and the next retry.


