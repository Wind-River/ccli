package pulsar

import (
	"log/slog"

	bus "bitbucket.wrs.com/scm/weststar/communication-bus.git"
	pulsar "bitbucket.wrs.com/scm/weststar/communication-bus.git/pulsar"
	pb "bitbucket.wrs.com/scm/weststar/pulsar-schemas-go.git"
	apachepulsar "github.com/apache/pulsar-client-go/pulsar"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// PulsarPartProducer handles parts and puts them onto part bus topic
type PulsarPartProducer struct {
	logger       *slog.Logger
	partProducer bus.ProtoProducer[*pb.PartMessage]
}

func NewPulsarPartProducer(logger *slog.Logger, host string, partTopic string) (*PulsarPartProducer, error) {
	ret := new(PulsarPartProducer)
	if logger != nil {
		ret.logger = logger
	} else {
		ret.logger = slog.Default()
	}

	// Create Pulsar Client
	client, err := pulsar.NewPulsarClient(nil, host, 6650, "", false)
	if err != nil {
		return nil, errors.Wrapf(err, "could not instantiate Pulsar client")
	}

	partProducer, err := pulsar.NewPulsarProtoProducer[*pb.PartMessage](client, partTopic)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create proto product line producer")

	}
	ret.partProducer = partProducer

	ret.logger.Info("Created Pulsar Part Producer", slog.String("host", host), slog.String("topic", partTopic))

	return ret, nil
}

// SendSchemaValue sends a Protobuf message onto the part topic
func (partProducer *PulsarPartProducer) SendPartSchemaValue(actionType pb.PartAction, partID string) error {

	partUUID, err := uuid.Parse(partID)
	if err != nil {
		return err
	}
	id, err := partUUID.MarshalBinary()
	if err != nil {
		return err
	}
	value := pb.PartMessage{
		Part: &pb.Part{
			Id: id,
		},
		Action: actionType,
	}

	// Send to part topic
	if messageID, err := partProducer.partProducer.Send(&value); err != nil {
		return errors.Wrapf(err, "error producing part message")
	} else {
		switch t := messageID.(type) {
		case apachepulsar.MessageID:
			partProducer.logger.Debug("produced part message", slog.String("messageID", string(t.Serialize())))
		}
	}

	return nil
}
