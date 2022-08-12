# Email Availability

## Purpose
Originally built as part of a customer request, now largely academic exercise in using the Nylas API to inject timeslots into an email using the scheduler. 

## Design
Webserver accepting API calls to POST `/send`, accepting the Nylas message body and authorization header, with additional information needed to call the scheduler `/timeslots` API.  The results of the timeslots API are then injected into the body of the message and the webserver then calls the Nylas `/send` API.  The scheduler is used to provide an easy method for fetching timeslots that match the user's preferences and provide an easy fallback option if timeslots do not meet the end user's needs.  The availability API could also be used instead of scheduler.

## Items to Note
- No support for EU DC 
- No auth for calls to server 
- Send requests are passed through to Nylas using the credentials provided

## Request example
```
curl --request POST \
  --url http://localhost:8000/send \
  --header 'Authorization: Bearer xxxx' \
  --header 'Content-Type: application/json' \
  --data '{
	"scheduler": "https://schedule.nylas.com/newtestnylas-30min-9",
	"timezone": "America/Denver",
	"name": "Nick",
	"email": "nickbair344@gmail.com",
	"maxTimeslots": 3,
	"useLinkTracking" true,
	"message": {
		"subject": "Insert availability demo",
		"from": [
			{
				"email": "nick.b@nylas.com"
			}
		],
		"to": [
			{
				"email": "nickbair344@gmail.com"
			}
		],
		"body": "Hey Nick, I think we should meet on this. {availability}"
	}
}'
```

## Response example
```
{
    "success": true,
    "data": {
        "account_id": "czf0dybf2x7hgzgaoqxqv9wir",
        "body": "Hey Nick, I think we should meet on this. <br><br><b>Book a meeting with me (America/Denver):</b><br><a href=https://schedule.nylas.com/newtestnylas-30min-9/book/1657911600?email=nickbair344%40gmail.com&name=Nick>Friday Jul 15, 2022 at 1:00 PM</a><br><a href=https://schedule.nylas.com/newtestnylas-30min-9/book/1657913400?email=nickbair344%40gmail.com&name=Nick>Friday Jul 15, 2022 at 1:30 PM</a><br><a href=https://schedule.nylas.com/newtestnylas-30min-9/book/1657915200?email=nickbair344%40gmail.com&name=Nick>Friday Jul 15, 2022 at 2:00 PM</a><br><br><i>None of these slots work for you? Book directly on my calendar <a href=https://schedule.nylas.com/newtestnylas-30min-9>here</a>.</i><br>",
        "date": 1657898709,
        "from": [
            {
                "email": "nick.b@nylas.com"
            }
        ],
        "id": "a5z6vyt5lsx2riw0ubqc2t0li",
        "labels": [
            {
                "display_name": "SENT",
                "id": "e86c9l611zsyv0hdo4qtvnqx",
                "name": "sent"
            }
        ],
        "object": "message",
        "snippet": "Hey Nick, I think we should meet on this. Book a meeting with me (America/Denver): Friday Jul 15, 2022 at 1:00 PM Friday Jul 15, 2022 at 1:30 PM Friday Jul 15, 2022 at 2:00 PM None of these s",
        "subject": "Insert availability demo",
        "thread_id": "4uddkkwda83nyp8tc55mqxqlk",
        "to": [
            {
                "email": "nickbair344@gmail.com"
            }
        ]
	}
}
```

## Error response example
```
{
	"success": "false",
	"errorMessage": "401 Unauthorized"
}
```