// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/thunder-id/thunderid/internal/flow/common"
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// RepresentationNodeInterface extends NodeInterface for representation nodes (START/END).
// These nodes use simple onSuccess navigation for linear flow.
type RepresentationNodeInterface interface {
	NodeInterface
	GetOnSuccess() string
	SetOnSuccess(nodeID string)
}

// representationNode implements the RepresentationNodeInterface
type representationNode struct {
	*node
	onSuccess string
}

// Ensure representationNode implements RepresentationNodeInterface
var _ RepresentationNodeInterface = (*representationNode)(nil)

// newRepresentationNode creates a new representation node
func newRepresentationNode(id string, nodeType common.NodeType, properties map[string]interface{},
	isStartNode bool, isFinalNode bool) NodeInterface {
	if properties == nil {
		properties = make(map[string]interface{})
	}
	return &representationNode{
		node: &node{
			id:               id,
			_type:            nodeType,
			properties:       properties,
			isStartNode:      isStartNode,
			isFinalNode:      isFinalNode,
			nextNodeList:     []string{},
			previousNodeList: []string{},
		},
		onSuccess: "",
	}
}

// Execute executes representation nodes with simple onSuccess navigation
func (n *representationNode) Execute(ctx *providers.NodeContext) (*common.NodeResponse, *tidcommon.ServiceError) {
	response := &common.NodeResponse{
		Status:         common.NodeStatusComplete,
		RuntimeData:    make(map[string]string),
		AdditionalData: make(map[string]string),
	}

	// Set next node using onSuccess property
	if n.onSuccess != "" {
		response.NextNodeID = n.onSuccess
	}

	return response, nil
}

// GetOnSuccess returns the onSuccess node ID for representation nodes
func (n *representationNode) GetOnSuccess() string {
	return n.onSuccess
}

// SetOnSuccess sets the onSuccess node ID for representation nodes
func (n *representationNode) SetOnSuccess(nodeID string) {
	n.onSuccess = nodeID
}
