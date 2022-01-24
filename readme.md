# DefAtt - Defence and Attack Platform

<p align="center">
  <img src="http://cybertraining.dk/defatt.png" alt="Defatt's Logo"/>
</p>
The platform is not ready for production usage. If you wish to use it, use your own responsibility. 

## How to run 

The platform is using gRPC communication between client and daemon to create specified games by administrators. 
As the time of writing this readme file, there are some components which need to be completed. 
Missing components are provided on the [missing components section]() below. 

### Run Daemon 

Example config files has been provided under [example-configs](example-configs/)

````bash 
$ go run app/daemon/main/main.go --config=<absolute-path-to-config-file>
````
Keep in mind that, games and information on this development stage are stateless, which means, they will NOT be recorded. (However, administrators information will be recorded to users.yml file. )

### Run Client 
The available functionalities to use are: 

- Listing Scenarios
- Create Game 
- Stopping Game
- Create/Modify/List Administrators 

Example calls to client can be found under [docs](docs/client.md)

## Scenarios

Creating scenarios for the platform is done by adding yml files to the folder scenarios [scenarios](scenarios/) following the structure provided in [example-configs](example-configs/).
Each yml file is a single scenario and an example is provided [here](example-configs/scenario.yml).
This example scenario has two networks where one of the machines is connected to both networks.

## Missing Components

- Scoring
- Web interface
- Administration web interface

