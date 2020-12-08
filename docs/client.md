# Client Information

This document includes client management for administrators of the platform, examples regarding to available gRPC calls and more information. 

## How to run

Codes for a client, are located under [../app/client](../app/client). As time of writing this readme file, it has following functionalities; 

- Create/Modify/List Administrators 
- Listing Scenarios
- Create Game 


Assistance regarding available commands can be taken with `--help` or `-h` flag to client code execution. For instance; 

```bash 
$ DEFAT_HOST=localhost DEFAT_SSL_OFF=true go run app/client/main.go --help 

Usage:
  defat [command]

Available Commands:
  challenges  List challenges in given scenario
  help        Help about any command
  scenarios   List available scenarios
  start       Start game with given scenario number
  user        Actions to perform on users

Flags:
  -h, --help   help for defat
```

In order to create game on the platform, at least one administrator should be signed up. For development purposes, the code can be executed as following; 

```bash 
$ DEFAT_HOST=localhost DEFAT_SSL_OFF=true go run app/client/main.go user signup 
```
The environment variables in front of `go run` is used for local development and where you do not provide any certificate. In case of production level, there will
be no need to provide any environment variables for running the code. 
Executed code above, will prompt to enter signup key which is provided by another administrator user, if there was no administrator user before, then daemon will log the signup key 
to the console. 

After successful login/signup, you need to know which scenario you would like to use for a game, in such a situation, available scenarios could be found with; 

```bash 
$  DEFAT_HOST=localhost DEFAT_SSL_OFF=true go run app/client/main.go scenarios list
   
   SCENARIO ID   DIFFICULTY   DURATION   NUMBER OF NETWORKS   STORY
   1             Easy         2          2                    Scenario 1 Storyy
   2             Moderate     3          4                    Scenario 2 Storyy
```

Let's create a game with scenario 1, 

```bash
$  DEFAT_HOST=localhost DEFAT_SSL_OFF=true go run app/client/main.go start -n "Test Game" -t "testgame" -s 1 
```

The game has name of "Test Game" and tag of "testgame" , tag will be the subdomain of the website where red/blue teams will signup.

Demo


