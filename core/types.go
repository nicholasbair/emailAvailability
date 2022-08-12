package core

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}

type MessageRequest struct {
	Message         Message `json:"message"`
	Scheduler       string  `json:"scheduler"`
	MaxTimeslots    int     `json:"maxTimeslots"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Timezone        string  `json:"timezone"`
	UseLinkTracking bool    `json:"useLinkTracking"`
}

type MessageResponse struct {
	Success      bool    `json:"success"`
	Data         Message `json:"data,omitempty"`
	ErrorMessage string  `json:"errorMessage,omitempty"`
}

type Message struct {
	AccountId string             `json:"account_id,omitempty"`
	Bcc       []EmailParticipant `json:"bcc,omitempty"`
	Body      string             `json:"body,omitempty"`
	Cc        []EmailParticipant `json:"cc,omitempty"`
	Date      int                `json:"date,omitempty"`
	Files     []File             `json:"files,omitempty"`
	From      []EmailParticipant `json:"from,omitempty"`
	Id        string             `json:"id,omitempty"`
	Labels    []OrganizationUnit `json:"labels,omitempty"`
	Folder    *OrganizationUnit  `json:"folder,omitempty"`
	Object    string             `json:"object,omitempty"`
	ReplyTo   []EmailParticipant `json:"reply_to,omitempty"`
	Snippet   string             `json:"snippet,omitempty"`
	Starred   bool               `json:"starred,omitempty"`
	Subject   string             `json:"subject,omitempty"`
	ThreadId  string             `json:"thread_id,omitempty"`
	To        []EmailParticipant `json:"to,omitempty"`
	Unread    bool               `json:"unread,omitempty"`
	Tracking  *Tracking          `json:"tracking,omitempty"`
}

type Tracking struct {
	Links bool `json:"links,omitempty"`
}

type EmailParticipant struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type OrganizationUnit struct {
	DisplayName string `json:"display_name,omitempty"`
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
}

type File struct {
	ContentType        string `json:"content_type,omitempty"`
	Filename           string `json:"filename,omitempty"`
	Object             string `json:"object,omitempty"`
	Size               int    `json:"size,omitempty"`
	ContentDisposition string `json:"content_disposition,omitempty"`
}

type Timeslot struct {
	AccountId  string   `json:"account_id"`
	CalendarId string   `json:"calendar_id"`
	Emails     []string `json:"emails"`
	End        int64    `json:"end"`
	HostName   string   `json:"host_name"`
	Start      int64    `json:"start"`
}

type WebhookRequest struct {
	Deltas []struct {
		Date       int    `json:"date"`
		Type       string `json:"type"`
		Object     string `json:"object"`
		ObjectData struct {
			Object    string `json:"object"`
			Id        string `json:"id"`
			AccountId string `json:"account_id"`
			Metadata  struct {
				//SenderAppId string `json:"sender_app_id"` ----- NOTE - Nylas is sending the internal ID here, filed bug report
				LinkData []struct {
					Url   string `json:"url"`
					Count int    `json:"count"`
				} `json:"link_data"`
				Payload   string `json:"payload"`
				MessageId string `json:"message_id"`
				Recents   []struct {
					Id        int    `json:"id"`
					Ip        string `json:"ip"`
					UserAgent string `json:"user_agent"`
					LinkIndex int    `json:"link_index"`
				} `json:"recents"`
			} `json:"metadata"`
		} `json:"object_data"`
	} `json:"deltas"`
}
