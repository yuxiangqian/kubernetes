/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"sync"
	"time"
)

// operationTimestamp stores the start time of an operation by a plugin
type operationTimestamp struct {
	pluginName string
	operation  string
	startTs    time.Time
}

func newOperationTimestamp(pluginName, operationName string) *operationTimestamp {
	return &operationTimestamp{
		pluginName: pluginName,
		operation:  operationName,
		startTs:    time.Now(),
	}
}

// OperationStartTimeCache concurrent safe cache for operation start timestamps
type OperationStartTimeCache struct {
	cache sync.Map
}

// NewOperationStartTimeCache creates a operation timestamp cache
func NewOperationStartTimeCache() OperationStartTimeCache {
	return OperationStartTimeCache{
		cache: sync.Map{},
	}
}

// AddIfNotExist returns directly if there exists an entry with the key. Otherwise, it
// creates a new operation timestamp using operationName, pluginName, and current timestamp
// and stores the operation timestamp with the key
func (c *OperationStartTimeCache) AddIfNotExist(key, pluginName, operationName string) {
	ts := newOperationTimestamp(pluginName, operationName)
	c.cache.LoadOrStore(key, ts)
}

// Delete deletes a value for a key.
func (c *OperationStartTimeCache) Delete(key string) {
	c.cache.Delete(key)
}

// Has returns a bool value indicates the existence of a key in the cache
func (c *OperationStartTimeCache) Has(key string) bool {
	_, exists := c.cache.Load(key)
	return exists
}

func (c *OperationStartTimeCache) Load(key string) (pluginName, operationName string, startTime time.Time, ok bool) {
	obj, exists := c.cache.Load(key)
	if !exists {
		return "", "", time.Time{}, false
	}
	ts, ok := obj.(*operationTimestamp)
	if !ok {
		return "", "", time.Time{}, false
	}
	return ts.pluginName, ts.operation, ts.startTs, true
}
