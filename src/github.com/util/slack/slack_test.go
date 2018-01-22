package slack

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	mockInfo     WebhookInfo
	mockChannel  = "bot_home"
	mockText     = "droipkg unit test"
	mockUsername = "merei"
	mockProxy    = "http://10.10.40.2:8080/"
	mockString   = "mock proxyClient"
	testURL      = "http://droibaas.com"
	mockPayload  = Payload{
		Username: mockUsername,
		Text:     mockText,
	}
	mockAttachment Attachment
)

func BeforeTest() {
	mockInfo = WebhookInfo{
		Payload: &mockPayload,
	}
}

func Test_SendMessageToSlack(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	// missing Payload
	assert.Error(t, SendMessageToSlack(WebhookInfo{}))

	client := &http.Client{
		Timeout: time.Duration(1) * time.Hour,
	}
	tmpInfo := WebhookInfo{
		CustomClient: client,
		Payload:      &mockPayload,
	}
	httpmock.ActivateNonDefault(client)
	httpmock.RegisterResponder("POST", GetSlackWebhookURL(),
		httpmock.NewStringResponder(http.StatusOK, mockString))
	assert.NoError(t, SendMessageToSlack(tmpInfo))
	// assert slack webhook is called exactly once
	httpmockInfo := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, httpmockInfo["POST "+GetSlackWebhookURL()])
}

func Test_getPayload(t *testing.T) {
	payload, err := getPayload(&mockInfo)
	assert.Nil(t, err)
	assert.Equal(t, "{\"username\":\"merei\",\"channel\":\"bot_home\",\"text\":\"droipkg unit test\"}", payload)
}

func Test_AddField(t *testing.T) {
	mockAttachment := Attachment{}
	fields := map[string]string{
		"b":   "",
		"123": "",
		"a":   "",
	}
	for k, v := range fields {
		field := Field{
			Title: k,
			Value: v,
		}
		mockAttachment.AddField(field)
	}
	assert.Equal(t, len(fields), len(mockAttachment.Fields))
	assert.Equal(t, "123", mockAttachment.Fields[0].Title)
	assert.Equal(t, "a", mockAttachment.Fields[1].Title)
	assert.Equal(t, "b", mockAttachment.Fields[2].Title)
}

func Test_fillDefaultValues(t *testing.T) {
	fillDefaultValues(&mockInfo)
	assert.Equal(t, defaultChannel, mockInfo.Payload.Channel)
}

// Do somethings after all test cases
func AfterTest() {
	httpmock.DeactivateAndReset()
}

func TestMain(m *testing.M) {
	BeforeTest()
	retCode := m.Run()
	AfterTest()
	os.Exit(retCode)
}
