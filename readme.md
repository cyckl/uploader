# Uploader
"Simple" webserver which handles file uploads

## Features
* Supports HTTP Basic Auth
* Random file name generator
* Can specify port, save location, and max file size
* bcrypt for authentication secret
* Written in Go? (I guess)

## Setting up
```
$ ./uploader -u <username> -s <secret>
$ ./uploader -w <domain> -d <path to webroot>
```

## How to upload
```
$ curl -u username:secret -F "data=@foo.txt" http://127.0.0.1:8080/upload 
```

## Accepted flags
```
Usage of uploader:
  -d string
    	Location to save files in
  -m int
    	The max file size in MB (default 10)
  -p string
    	The port to bind to (default "8080")
  -s string
    	Set a new auth secret
  -u string
    	Set a new auth username
  -w string
    	Public-facing URL for server
```

## Todo
* Multi-user support
	* How should I implement this without adding a whole bunch of code in order to deal with user management
* Auto-deletion
* Dedicated web interface with auth
* Add runit service
