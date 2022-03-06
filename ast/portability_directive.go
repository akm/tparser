package ast

type PortabilityDirective string

const (
	PdPlatform   PortabilityDirective = "platform"
	PdDeprecated PortabilityDirective = "deprecated"
	PdLibrary    PortabilityDirective = "library"
)
