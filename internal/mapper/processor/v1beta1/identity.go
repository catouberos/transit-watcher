package processorv1beta1mapper

import (
	"errors"

	processorv1beta1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/processor/v1beta1"
	"codeberg.org/transit-radar/transit-watcher/internal/models"
)

func MapIdentity(identity models.Identity) (*processorv1beta1.Identity, error) {
	var identifier processorv1beta1.ExternalIdentifier
	switch identity.Identifier {
	case models.ExternalIdentifierEBMS:
		identifier = processorv1beta1.ExternalIdentifier_EXTERNAL_IDENTIFIER_EBMS
	case models.ExternalIdentifierMultiGo:
		identifier = processorv1beta1.ExternalIdentifier_EXTERNAL_IDENTIFIER_MULTIGO
	default:
		return nil, errors.New("unsupported identifier")
	}

	return &processorv1beta1.Identity{
		Identifier: identifier,
		Value:      identity.Value,
	}, nil
}
