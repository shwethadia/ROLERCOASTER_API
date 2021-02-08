# Rolercoaster REST API

* Implementation of REST API for ROLERCOASTER GAME
* Tested Using Postamn and Curl as an HTTP Clinet
* No Third Party Packages and Dependencies


## Requirements

* ` GET /coasters ` returns list of coasters as JSON
* ` GET /coasters/{id} ` returns details of specific coaster as JSON
* ` POST /coasters ` accepts a new coaster to be added
* ` POST /coasters ` returns status 415 if content is not `application/json`
* ` GET /admin ` requires auth
* ` GET /coasters/random ` redirects to a random coaster


### Data Types

ROALERCOASTER OBJECT 

``` json 
{
    "id": "ID (int)",
    "name":"NAME OF THE ROLERCOASTER(string)",
    "inPark":"Name OF THE PARK (string) ",
    "manufacturer": "NAME OF THE MANUFACTURER (string)",
    "height":"27(int)",
}
``` 

### COMMANDS

- Get 
    * localhost:8080/coasters
- POST 
    * curl localhost:8080/coasters -X POST -d '{"name":"Taron","inpark":"Phantasialand","height":30,"manufacturer":"Intamin"}' -H "Content-Type:application/json"
- GET 
    * curl localhost:8080/coasters/1612775957075876000
- GET 
    * curl localhost:8080/admin -u admin:secret
- GET 
    * curl localhost:8080/coasters/random -L