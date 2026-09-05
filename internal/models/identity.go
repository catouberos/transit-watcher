package models

type ExternalIdentifier string

const (
	ExternalIdentifierEBMS    ExternalIdentifier = "EXTERNAL_IDENTIFIER_EBMS"
	ExternalIdentifierMultiGo ExternalIdentifier = "EXTERNAL_IDENTIFIER_MULTIGO"
)

type Identity struct {
	Identifier ExternalIdentifier `json:"identifier" redis:"identifier"`
	Value      string             `json:"value" redis:"value"`
}
