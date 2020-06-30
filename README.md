<div align="center">
    <h2>Shorten Bulky Links</h2>
    <h1>Short-URL</h1>
    <h4>Build simple and reliable short links.</h4>
</div>

<p align="center">
    <a href = "#about">About</a> |
    <a href = "#api">API</a> |
    <a href = "#installation">Installation</a> |
    <a href = "#license">License</a> 
</p>

## About
> Short-URL is used to make six characters URL and it can be retrived at lightning speed.

It has an amazing API for getting the click details of each shorturl that has been clicked. 


<p align="center">
  <img src="assets/shorturl.gif" />
</p>

## API

#### Get all click details for particular shorturl.
```
https://shrt-url.xyz/api/v1/{shorturl}?skip=0&limit=100
```
Example:  {shorturl} = 52ea82r. Try it out [here.](https://shrt-url.xyz/#/Click%20Details/getAllClickDetails)

## Installation
To run the application you need to type below mentioned command.
```bash
go run main.go
```
To successfully run the application you need to set below mentioned environment variables on your system.

```
primaryDB_name = 
primaryDB_host = 
primaryDB_port = 
cacheDB_name =  
cacheDB_host = 
cacheDB_port =
host = 
ipstack_apiKey = Ex: API Key for ipstack
env = DEV
fullchain= Ex: /path/fullchain.pem
privkey= Ex: /path/privkey.pem
privateToken=Ex: Captcha Token.
```