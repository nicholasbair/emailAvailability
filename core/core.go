package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const availabilityTag = "{availability}"

func InjectAvailabilityAndSendMessage(r []byte, headers map[string]string) (Message, error) {
	m := new(MessageRequest)

	if unmarshalErr := json.Unmarshal(r, m); unmarshalErr != nil {
		return Message{}, &RequestError{
			StatusCode: 400,
			Err:        unmarshalErr,
		}
	}

	if strings.Contains(m.Message.Body, availabilityTag) {
		timeslots, availErr := getAvailability(m.Scheduler)

		if availErr != nil {
			return Message{}, availErr
		}

		URLs, buildErr := buildAvailabilityUrls(timeslots, *m)

		if buildErr != nil {
			return Message{}, &RequestError{
				StatusCode: 422,
				Err:        buildErr,
			}
		}

		m.Message.Body = strings.Replace(m.Message.Body, availabilityTag, URLs, -1)

		if m.UseLinkTracking {
			m.Message.Tracking = &Tracking{Links: true}
		}

		res, sendErr := send(m.Message, getAuthHeader(headers))

		if sendErr != nil {
			return Message{}, sendErr
		}

		return res, nil
	} else {
		return Message{}, &RequestError{
			StatusCode: 400,
			Err:        errors.New("missing availability tag in email body"),
		}
	}
}

func buildAvailabilityUrls(timeslots []Timeslot, r MessageRequest) (string, error) {
	output := fmt.Sprintf("<br><br><b>Book a meeting with me (%s):</b><br>", r.Timezone)

	if len(timeslots) < 1 {
		return output, errors.New("not enough timeslots")
	}

	for i, t := range timeslots {
		// Create max of N links for email
		if i >= r.MaxTimeslots {
			break
		}
		start := strconv.Itoa(int(t.Start))
		loc, _ := time.LoadLocation(r.Timezone)
		datetime := time.Unix(t.Start, 0).In(loc)

		output = output + fmt.Sprintf("<a href='%s'>%s</a><br>", buildBookingUrl(r, start), datetime.Format("Monday Jan 2, 2006 at 3:04 PM"))
	}

	output = output + fmt.Sprintf("<br><i>None of these slots work for you? Book directly on my calendar <a href='%s'>here</a>.</i><br>", r.Scheduler)

	return output, nil
}

func buildBookingUrl(r MessageRequest, time string) string {
	v := url.Values{}
	v.Add("name", r.Name)
	v.Add("email", r.Email)
	params := v.Encode()

	uri := fmt.Sprintf("%s/book/%s?%s", r.Scheduler, time, params)

	return uri
}

// NOTE - based on testing, no need to also check for lowercase authorization
func getAuthHeader(headers map[string]string) string {
	return headers["Authorization"]
}
