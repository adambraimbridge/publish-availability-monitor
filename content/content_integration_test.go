package content

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"github.com/stretchr/testify/assert"
)

func TestThat_ContentPayloadFromKafkaMsg_MatchesTheContentPayloadSentToValidationServices(t *testing.T) {
	var testCases = []struct {
		name string
		msg  consumer.Message
	}{
		{
			"Methode List",
			consumer.Message{
				Headers: map[string]string{
					"Origin-System-Id": "http://cmdb.ft.com/systems/methode-web-pub",
				},
				Body: string(loadBytesForFile(t, "methode_list.json")),
			},
		},
		{
			"Methode Empty List",
			consumer.Message{
				Headers: map[string]string{
					"Origin-System-Id": "http://cmdb.ft.com/systems/methode-web-pub",
				},
				Body: string(loadBytesForFile(t, "methode_empty_list.json")),
			},
		},
		{
			"Methode Article",
			consumer.Message{
				Headers: map[string]string{
					"Origin-System-Id": "http://cmdb.ft.com/systems/methode-web-pub",
				},
				Body: string(loadBytesForFile(t, "methode_article.json")),
			},
		},
	}

	for _, tc := range testCases {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, strings.TrimSpace(tc.msg.Body), strings.TrimSpace(string(body)), "[%s] testcase failure", tc.name)
		}))
		defer ts.Close()
		content, err := UnmarshalContent(tc.msg)
		if err != nil {
			t.Errorf("Expected success, but error occured [%v]", err)
			return
		}
		//don't care about the actual result
		content.IsValid(ts.URL, "tid_foobar", "", "")
	}
}
