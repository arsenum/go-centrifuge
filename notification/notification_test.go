// +build unit

package notification

import (
	"os"
	"testing"
	"time"

	"github.com/centrifuge/centrifuge-protobufs/documenttypes"
	"github.com/centrifuge/centrifuge-protobufs/gen/go/coredocument"
	"github.com/centrifuge/centrifuge-protobufs/gen/go/notification"
	"github.com/centrifuge/go-centrifuge/bootstrap"
	"github.com/centrifuge/go-centrifuge/config"
	"github.com/centrifuge/go-centrifuge/context/testlogging"
	"github.com/centrifuge/go-centrifuge/storage"
	"github.com/centrifuge/go-centrifuge/testingutils"
	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	ibootstappers := []bootstrap.TestBootstrapper{
		&testlogging.TestLoggingBootstrapper{},
		&config.Bootstrapper{},
		&storage.Bootstrapper{},
	}
	bootstrap.RunTestBootstrappers(ibootstappers, nil)
	result := m.Run()
	bootstrap.RunTestTeardown(ibootstappers)
	os.Exit(result)
}

func TestWebhookConstructPayload(t *testing.T) {
	documentIdentifier := utils.RandomSlice(32)
	coredoc := &coredocumentpb.CoreDocument{DocumentIdentifier: documentIdentifier}
	cid := testingutils.Rand32Bytes()

	ts, err := ptypes.TimestampProto(time.Now().UTC())
	assert.Nil(t, err, "Should not error out")
	notificationMessage := &notificationpb.NotificationMessage{
		DocumentIdentifier: coredoc.DocumentIdentifier,
		DocumentType:       documenttypes.InvoiceDataTypeUrl,
		CentrifugeId:       cid,
		EventType:          uint32(ReceivedPayload),
		Recorded:           ts,
	}

	whs := WebhookSender{}
	bresult, err := whs.constructPayload(notificationMessage)
	assert.Nil(t, err, "Should not error out")

	unmarshaledNotificationMessage := &notificationpb.NotificationMessage{}

	jsonpb.UnmarshalString(string(bresult), unmarshaledNotificationMessage)

	assert.Equal(t, notificationMessage.Recorded, unmarshaledNotificationMessage.Recorded, "Recorder Timestamp should be equal")
	assert.Equal(t, notificationMessage.DocumentType, unmarshaledNotificationMessage.DocumentType, "DocumentType should be equal")
	assert.Equal(t, notificationMessage.DocumentIdentifier, unmarshaledNotificationMessage.DocumentIdentifier, "DocumentIdentifier should be equal")
	assert.Equal(t, notificationMessage.CentrifugeId, unmarshaledNotificationMessage.CentrifugeId, "CentrifugeID should be equal")
	assert.Equal(t, notificationMessage.EventType, unmarshaledNotificationMessage.EventType, "EventType should be equal")
}
