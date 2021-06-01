package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/sirupsen/logrus"
)

type Parser struct {
	Messages map[string]*desc.MessageDescriptor
	Methods  map[string]*desc.MethodDescriptor
}

func NewParser(protoPath string) (*Parser, error) {
	parser := &Parser{
		Messages: make(map[string]*desc.MessageDescriptor),
		Methods:  make(map[string]*desc.MethodDescriptor),
	}

	protoparser := protoparse.Parser{
		ImportPaths:      []string{protoPath},
		InferImportPaths: true,
	}

	err := filepath.Walk(protoPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".proto" {
			_, err := protoparser.ParseFilesButDoNotLink(path)
			if err != nil {
				logrus.WithError(err).Warn("Skipping %s", path)
				return nil
			}

			relativePath, err := filepath.Rel(protoPath, path)
			if err != nil {
				return err
			}

			descriptions, err := protoparser.ParseFiles(relativePath)
			if err != nil {
				logrus.WithError(err).Warn("Failed to parse %s", path)
				return err
			}

			for _, desc := range descriptions {
				for _, messageType := range desc.GetMessageTypes() {
					parser.Messages[messageType.GetFullyQualifiedName()] = messageType
				}
				for _, serviceType := range desc.GetServices() {
					for _, methodType := range serviceType.GetMethods() {
						fullyQualifiedName := fmt.Sprintf("/%s/%s", serviceType.GetFullyQualifiedName(), methodType.GetName())
						parser.Methods[fullyQualifiedName] = methodType
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return parser, nil
}

func (p *Parser) ParseMessage(fullMethod string, origin MessageOrigin, m interface{}) (*dynamic.Message, error) {
	messageType, err := p.findMessageType(fullMethod, origin)
	if err != nil {
		return nil, err
	}

	data := []byte{}
	switch message := m.(type) {
	case []byte:
		data = message
	case *[]byte:
		data = *message
	}

	message := dynamic.NewMessage(messageType)
	err = proto.Unmarshal(data, message)
	return message, err
}

func (p *Parser) findMessageType(fullMethod string, origin MessageOrigin) (*desc.MessageDescriptor, error) {
	if methodType, ok := p.Methods[fullMethod]; ok {
		switch origin {
		case HeadMessage:
			return methodType.GetInputType(), nil
		case RuntimeMessage:
			return methodType.GetOutputType(), nil
		}
	}
	return nil, errors.New("method not found")
}
