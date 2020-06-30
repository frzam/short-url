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

* #### Get all click details for particular shorturl.
    ```
    https://shrt-url.xyz/api/v1/{shorturl}?skip=0&limit=100
    ```
    **Example:**  {shorturl} = 52ea82r. Try it out [here.](https://shrt-url.xyz)

* #### Get all click details for particular shorturl for past N days.
    ```
    https://shrt-url.xyz/api/v1/{shorturl}/{days}?skip=0&limit=100
    ```
    **Example:**  {shorturl} = 52ea82r *and* {days} = 7. Try it out [here.](https://shrt-url.xyz)

* #### Delete all the details for a url.
    ```
    https://shrt-url.xyz/api/v1/{shorturl}
    ```
    **Example:**  Method:  *DELETE*. Try it out [here.](https://shrt-url.xyz)

* #### Get all the click details by a country for one shorturl.
    ```
    https://shrt-url.xyz/api/v1/{shorturl}/country/{country}
    ```
    **Example:**  {shorturl} = 52ea82r *and* {country} = India. Try it out [here.](https://shrt-url.xyz)


* #### Get all the click details by a city for one shorturl.
    ```
    https://shrt-url.xyz/api/v1/{shorturl}/city/{city}
    ```
    **Example:**  {shorturl} = 52ea82r *and* {city} = Powai. Try it out [here.](https://shrt-url.xyz)


* #### Get all the clickk details by an IP for one shorturl.
    ```
    https://shrt-url.xyz/api/v1/{shorturl}/ip/{ip}
    ```
    **Example:**  {shorturl} = 52ea82r. Try it out [here.](https://shrt-url.xyz)


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