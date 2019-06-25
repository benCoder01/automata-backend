# Raspberry-Home-Automation

This is the repository for the raspberry pi home automation application. It is driven by a raspberry pi 1 model B, which hosts a api based on the http protocol with the routes for activating certain GPIO-pins. (work in progress)

Additionally the raspberry hosts a website for activating the requests to the REST-api. (not implemented)

As the backend of the application is written in GO, this folder should be located in your gopath.

## How to use it

http get-request for activating a switch:

```
/activate?id=1
```