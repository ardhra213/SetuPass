// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/flow/common"
)

type NodeInterfaceTestSuite struct {
	suite.Suite
}

func TestNodeInterfaceTestSuite(t *testing.T) {
	suite.Run(t, new(NodeInterfaceTestSuite))
}

func (s *NodeInterfaceTestSuite) TestStartNodeImplementsRepresentationNodeInterface() {
	node := newRepresentationNode("start", common.NodeTypeStart, nil, true, false)
	_, ok := node.(RepresentationNodeInterface)
	s.True(ok, "START node should implement RepresentationNodeInterface")
}

func (s *NodeInterfaceTestSuite) TestEndNodeImplementsRepresentationNodeInterface() {
	node := newRepresentationNode("end", common.NodeTypeEnd, nil, false, true)
	_, ok := node.(RepresentationNodeInterface)
	s.True(ok, "END node should implement RepresentationNodeInterface")
}

func (s *NodeInterfaceTestSuite) TestPromptNodeDoesNotImplementRepresentationNodeInterface() {
	node := newPromptNode("prompt", nil, false, false)
	_, ok := node.(RepresentationNodeInterface)
	s.False(ok, "PROMPT node should NOT implement RepresentationNodeInterface")
}

func (s *NodeInterfaceTestSuite) TestPromptNodeDoesNotImplementExecutorBackedNodeInterface() {
	node := newPromptNode("prompt", nil, false, false)
	_, ok := node.(ExecutorBackedNodeInterface)
	s.False(ok, "PROMPT node should NOT implement ExecutorBackedNodeInterface")
}

func (s *NodeInterfaceTestSuite) TestTaskExecutionNodeImplementsExecutorBackedNodeInterface() {
	node := newTaskExecutionNode("task", nil, false, false)
	_, ok := node.(ExecutorBackedNodeInterface)
	s.True(ok, "TASK_EXECUTION node should implement ExecutorBackedNodeInterface")
}

func (s *NodeInterfaceTestSuite) TestPromptNodeImplementsPromptNodeInterface() {
	node := newPromptNode("prompt", nil, false, false)
	_, ok := node.(PromptNodeInterface)
	s.True(ok, "PROMPT node should implement PromptNodeInterface")
}
