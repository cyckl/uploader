# Uploader
"Simple" webserver which handles file uploads

## Features
* Supports HTTP Basic Auth
* Random file name generator
* Can specify port, save location, and max file size
* bcrypt for authentication secret
* Only one user unfortunately
* Written in Go? (I guess)

## How to use
```
curl -u username:password -F "data=@example.text" http://127.0.0.1:8080/upload 
```

## Setting new credentials
The default credentials are `username:password`. I recommend changing them.
1. `./uploader -u <username>`
2. `./uploader -s <secret>`

## Help dialog
```
Usage of uploader:
  -d string
    	Location to save files in (default "files/")
  -m int
    	The max file size in MB (default 10)
  -p string
    	The port to bind to (default "8080")
  -s string
    	Set a new auth secret
  -u string
    	Set a new auth username
```

## Todo
* Add the random three word filename thing
* Multi-user support
* Auto-deletion
* Domain-name / hostname support in HTTP response
* Dedicated web interface with auth
