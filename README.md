# Swanson Service
***

[[_TOC_]]

***

## Purpose

This project shows an example of Kubernetes cronjobs. The cronjob will call out to a service and send the content to a Slack channel at a configured scheduled time.

## Components

The main components of this example are:

- Golang service used to call out to a set of public APIs.
- ConfigMap for setting the Slack webhook URL, target service URLs, and target Slack channel/user ID. 
- CronJob configs for setting the schedule and for including a `curl` command for triggering the service to send info to Slack.
- Service config with NodePort for making it easier to test the service from my local machine.


### Service

The service uses the following format: `http://host:30080/{API}`

| API | Description |
|:---:|---|
| nasa | Public API for getting NASA's Astronomy Picture of the Day. |
| swanson | Public API for getting random Ron Swanson quotes. |
| xkcd | Public API for getting the current day's XKCD. |

### ConfigMaps

| Name | Description |
|:---|---|
| slacker-configmap | Contains the URL to the Slack incoming webhook. |
| swanson-configmap | Contains the URL to the Ron Swanson Quotes API and the target channel for the messages. |
| nasa-configmap | Contains the URL to the NASA API and the target channel for the messages. |
| xkcd-configmap | Contains the URL to the XKCD API and the target channel for the messages. |

### CronJobs

The `*-schedule.yaml` files will create cronjobs for the configured open APIs.

### Service

The service defines a nodeport.

#### Data Retrievers

##### NASA

Target URL: `https://api.nasa.gov/planetary/apod?api_key=NASA_API_KEY`

```json
{
   "copyright":"Shi Huan",
   "date":"2021-07-16",
   "explanation":"Venus, named for the Roman goddess of love, and Mars, the war god's namesake, come together by moonlight in this serene skyview, recorded on July 11 from Lualaba province, Democratic Republic of Congo, planet Earth. Taken in the western twilight sky shortly after sunset the exposure also records earthshine illuminating the otherwise dark surface of the young crescent Moon. Of course the Moon has moved on. Venus still shines in the west though as the evening star, third brightest object in Earth's sky, after the Sun and the Moon itself. Seen here above a brilliant Venus, Mars moved even closer to the brighter planet and by July 13 could be seen only about a Moon's width away. Mars has since slowly wandered away from much brighter Venus in the twilight, but both are sliding toward bright star Regulus. Alpha star of the constellation Leo, Regulus lies off the top of this frame and anticipates a visit from Venus and then Mars in twilight skies of the coming days.",
   "hdurl":"https://apod.nasa.gov/apod/image/2107/2021Jul11MarsVenusMoon_ShiHuan.jpg",
   "media_type":"image",
   "service_version":"v1",
   "title":"Love and War by Moonlight",
   "url":"https://apod.nasa.gov/apod/image/2107/2021Jul11MarsVenusMoon_ShiHuan1024.jpg"
}
```

The NASA data retriever will use the copyright, date, explanation, title, and URL to build the Slack message.
	
##### XKCD

Target URL: `https://xkcd.com/info.0.json`

```json
{
  "month":"7",
  "num":2490,
  "link":"",
  "year":"2021",
  "news":"",
  "safe_title":"Pre-Pandemic Ketchup",
  "transcript":"",
  "alt":"I wonder what year I'll discard the last weird food item that I bought online in early 2020.",
  "img":"https://imgs.xkcd.com/comics/pre_pandemic_ketchup.png",
  "title":"Pre-Pandemic Ketchup",
  "day":"16"
}
```

The XKCD data retriever will use the alt, title, img, year, month, day, and num to build the Slack message.

##### Ron Swanson

Target URL: `http://ron-swanson-quotes.herokuapp.com/v2/quotes`

```json
["Every two weeks I need to sand down my toe nails. They're too strong for clippers."]
```

The Swanson data retriever will extract the first string from the returned JSON string array to build the Slack message.

of three services - [Ron Swanson Quote service](https://github.com/jamesseanwright/ron-swanson-quotes) to get a random quote, and then will post the quote to a Slack channel.



## Update the Schedule

Here is an example of how you can update the schedule:

```bash
kubectl patch cronjob get-swanson-quote-cron \
  -p '{"spec":{"schedule": "30 9 * * *"}}'
```

