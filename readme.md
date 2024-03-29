# DefAtt - Defence and Attack Platform

<p align="center">
  <img src="http://cybertraining.dk/defatt.png" alt="Defatt's Logo"/>
</p>
The platform is not ready for production usage. If you wish to use it, use your own responsibility. 

## How to run 

The platform is using gRPC communication between client and daemon to create specified games by administrators. As the time of writing this readme file, there are 
some components which need to be completed. Missing components are provided on the [missing components section]() below. 

### Run Daemon 

Example config file has been provided under [config dir](config/)

````bash 
$ go run app/daemon/main/main.go --config=<absolute-path-to-config-file>
10:00AM INF Started daemon
10:00AM INF gRPC daemon has been started  ! on port :5454
10:00AM INF Reflection Registration is called....
````
Keep in mind that, games and information on this development stage are stateless, which means, they will NOT be recorded. (However, administrators information will be recorded to users.yml file. )

Example config file is [here](./config/config.yml)

### Run Client 
The available functionalities to use are: 

- Listing Scenarios
- Create Game 
- Create/Modify/List Administrators 

Example calls to client can be found under [docs](docs/client.md)

## Missing Components

- Scoring
- Web interface
- Administration web interface
- Scenarios

