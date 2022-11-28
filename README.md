# Peloton To Garmin

This is a simple CLI tool that will pull the last x Peloton workouts and upload them to Garmin Connect. To use this tool you need a Peloton username and password as well as a Garmin Connect username and password. 

Currently this tool only supports Peloton Cycling and stretching workouts. 

Please note, Peloton do not publicly publish their API documentation, so things could break if Peloton decide to change how their API performs. If you notice issues, please create a github issue and I'll take a look when I can. If you like this project and use it please do watch and star the repo. 

Example usage: 

```
peloton-to-garmin.exe sync --pelotonUsername joeblogs@hotmail.com --pelotonPassword 'toSecretPassword' --garminEmail joeblogs@hotmail.com -garminPassword 'ToSecretToTellAnyone'
```


## Default Options
By default, this cli will lookup your last 30 workouts from Peloton and attempt to upload them. It will not overwrite existing workouts. Re-running this tool again will simply output the workout already exists in Garmin.  You can also set your datapoint granularity from Peloton. The default is set to a datapoint per second, but this could be changed if needed to something less ganular if required. 

To see optional options, you can run `peloton-to-garmin.exe sync --help`


## Still To Do

This is a work in progress project and some of the things I'd like to do as I get time are:

* Implement a http.client logger so trace logging can show request and response
* Improve my usage of the http.client and try and reduce duplicate code
* Support additional workout types from Peloton
