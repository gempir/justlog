# justlog [![Build Status](https://travis-ci.org/gempir/justlog.svg?branch=ansible)](https://travis-ci.org/gempir/justlog)

#### What is this?
justlog is a bot I maintain for a couple of channels. Its features differ from other bots in that it doesn't support
commands, etc. yet, it focuses on logging and providing an api for the logs.

## API

**Root Domain** `https://justlog.gempir.com`

##### All Requests only with HTTPS!


#### RandomQuote

Responds with a random message from the given user in that channel.


* **URL** `/channel/%channel%/user/%username%/random`

* **Methods** `GET` 

* **Parameter**
  -  `%channel%` Channel from which you want the quote
  - `%username%` Username from whom you want the quote

* **Success Response**

  * **Code:** 200 
    **Content:** `gempir: My Random Message`
 
* **Error Response:** `{"Error":"error finding logs"}`

- **Notice**

  With the Header Content-Type: "application/json" the response will be in json format. Example:
`{
  "channel": "gempbot",
  "username": "gempir",
  "message": "MingLee",
  "timestamp": "2017-03-12 09:43:22"
}`

#### Userlog

Responds with the timestamped log of a user in the current month for the given channel.


* **URL** `/channel/%channel%/user/%username%`

* **Methods** `GET` 

* **Parameter**
  -  `%channel%` Channel that the user wrote in
  - `%username%` Username that typed in the chat

* **Success Response**

  * **Code:** 200 
    **Content:**

[2017-04-1 06:40:24] gempir: HaHaa

[2017-04-1 06:40:36] gempir: KKona

[2017-04-1 06:40:43] gempir: NaM

 
* **Error Response:** `{"message":"Not Found"}`
* **Notice**

  * This endpoint will redirect to the current month and year so the url will say `/channel/%channel%/user/%username%/2017/April`

#### Channellog

Responds with the timestamped log of a channel of the current day for the given channel.


* **URL** `/channel/%channel%`

* **Methods** `GET` 

* **Parameter**
  -  `%channel%` Channel you want the logs for

* **Success Response**

  * **Code:** 200 
    **Content:**

[2017-04-1 06:40:24] gempir: PHPDETECTED

[2017-04-1 06:40:36] pajlada: KKaper

[2017-04-1 06:40:43] fourtf: fuck qt

 
* **Error Response:** `{"message":"Not Found"}`
* **Notice**

  * This endpoint will redirect to the current day, month and year so the url will say `/channel/%channel%/2017/April/20`
