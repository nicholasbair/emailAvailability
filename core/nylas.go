package core

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func send(m Message, auth string) (Message, error) {
	payload, marshalErr := json.Marshal(m)

	if marshalErr != nil {
		return Message{}, &RequestError{
			StatusCode: 400,
			Err:        marshalErr,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.nylas.com/send", bytes.NewBuffer(payload))
	req.Header.Add("Authorization", auth)
	res, httpErr := http.DefaultClient.Do(req)

	if httpErr != nil {
		return Message{}, &RequestError{
			StatusCode: res.StatusCode,
			Err:        httpErr,
		}
	}

	if res.StatusCode != 200 {
		return Message{}, &RequestError{
			StatusCode: res.StatusCode,
			Err:        errors.New(res.Status),
		}
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var response Message

	if unmarshalErr := json.Unmarshal(body, &response); unmarshalErr != nil {
		return Message{}, &RequestError{
			StatusCode: 500,
			Err:        unmarshalErr,
		}
	}

	return response, nil
}

func getAvailability(schedulerURL string) ([]Timeslot, error) {
	u, _ := url.Parse(schedulerURL)
	u.Host = "api.schedule.nylas.com"
	u.Path = "/schedule" + u.Path + "/timeslots"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("allow_stale", "true")
	req.URL.RawQuery = q.Encode()
	res, httpErr := http.DefaultClient.Do(req)

	if httpErr != nil {
		return nil, &RequestError{StatusCode: 500, Err: httpErr}
	}

	if res.StatusCode != 200 {
		return nil, &RequestError{StatusCode: res.StatusCode, Err: errors.New(res.Status)}
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var t []Timeslot

	if err := json.Unmarshal(body, &t); err != nil {
		return nil, &RequestError{StatusCode: 500, Err: err}
	}

	return t, nil
}

func (w *WebhookRequest) LogInfo() {
	t := w.Deltas[0].ObjectData.Metadata.LinkData[w.Deltas[0].ObjectData.Metadata.Recents[0].LinkIndex]
	log.Printf("Link clicked: %s\n", t.Url)
}

// CheckSignature - Copied from https://github.com/Teamwork/nylas-go/blob/master/webhook.go
func CheckSignature(secret, signature string, body []byte) error {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write(body)
	if err != nil {
		return err
	}
	if signature != hex.EncodeToString(mac.Sum(nil)) {
		return errors.New("signature mismatch")
	}
	return nil
}
