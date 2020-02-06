/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/
package cache

import (
	"fmt"

	"github.com/apache/incubator-yunikorn-k8shim/pkg/common/events"
)

// implements ApplicationManagementProtocol
type MockedAMProtocol struct {
	applications map[string]*Application
}

func NewMockedAMProtocol() *MockedAMProtocol {
	return &MockedAMProtocol{
		applications: make(map[string]*Application)}
}

func (m *MockedAMProtocol) GetApplication(appID string) (*Application, bool) {
	if app, ok := m.applications[appID]; ok {
		return app, true
	}
	return nil, false
}

func (m *MockedAMProtocol) AddApplication(request *AddApplicationRequest) (*Application, bool) {
	if app, ok := m.GetApplication(request.Metadata.ApplicationID); ok {
		return app, false
	}

	app := NewApplication(
		request.Metadata.ApplicationID,
		request.Metadata.QueueName,
		request.Metadata.User,
		request.Metadata.Tags,
		nil)

	// add into cache
	m.applications[app.GetApplicationID()] = app

	switch request.Recovery {
	case true:
		app.SetState(events.States().Application.Recovering)
	case false:
		app.SetState(events.States().Application.New)
	}

	return app, true
}

func (m *MockedAMProtocol) RemoveApplication(appID string) error {
	if _, ok := m.GetApplication(appID); ok {
		delete(m.applications, appID)
		return nil
	}
	return fmt.Errorf("application doesn't exist")
}

func (m *MockedAMProtocol) AddTask(request *AddTaskRequest) (*Task, bool) {
	if app, ok := m.applications[request.Metadata.ApplicationID]; ok {
		if existingTask, err := app.GetTask(request.Metadata.TaskID); err != nil {
			task := NewTask(request.Metadata.TaskID, app, nil, request.Metadata.Pod)
			app.addTask(&task)
			return &task, true
		} else {
			return existingTask, false
		}
	} else {
		return nil, false
	}
}

func (m *MockedAMProtocol) RemoveTask(appID, taskID string) error {
	if app, ok := m.applications[appID]; ok {
		return app.removeTask(taskID)
	} else {
		return fmt.Errorf("app not found")
	}
}

func (m *MockedAMProtocol) NotifyApplicationComplete(appID string) {
	if app, ok := m.GetApplication(appID); ok {
		app.SetState(events.States().Application.Completed)
	}
}

func (m *MockedAMProtocol) NotifyTaskComplete(appID, taskID string) {
	if app, ok := m.GetApplication(appID); ok {
		if task, err := app.GetTask(taskID); err == nil {
			task.sm.SetState(events.States().Task.Completed)
		}
	}
}