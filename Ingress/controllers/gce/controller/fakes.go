/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package controller

import (
	"k8s.io/contrib/Ingress/controllers/gce/backends"
	"k8s.io/contrib/Ingress/controllers/gce/healthchecks"
	"k8s.io/contrib/Ingress/controllers/gce/instances"
	"k8s.io/contrib/Ingress/controllers/gce/loadbalancers"
	"k8s.io/contrib/Ingress/controllers/gce/utils"
	"k8s.io/kubernetes/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/util/sets"
)

const testDefaultBeNodePort = int64(3000)

var testBackendPort = intstr.IntOrString{Type: intstr.Int, IntVal: 80}

// ClusterManager fake
type fakeClusterManager struct {
	*ClusterManager
	fakeLbs      *loadbalancers.FakeLoadBalancers
	fakeBackends *backends.FakeBackendServices
	fakeIGs      *instances.FakeInstanceGroups
}

// NewFakeClusterManager creates a new fake ClusterManager.
func NewFakeClusterManager(clusterName string) *fakeClusterManager {
	fakeLbs := loadbalancers.NewFakeLoadBalancers(clusterName)
	fakeBackends := backends.NewFakeBackendServices()
	fakeIGs := instances.NewFakeInstanceGroups(sets.NewString())
	fakeHCs := healthchecks.NewFakeHealthChecks()
	namer := utils.Namer{clusterName}
	nodePool := instances.NewNodePool(fakeIGs)
	healthChecker := healthchecks.NewHealthChecker(fakeHCs, "/", namer)
	backendPool := backends.NewBackendPool(
		fakeBackends,
		healthChecker, nodePool, namer)
	l7Pool := loadbalancers.NewLoadBalancerPool(
		fakeLbs,
		// TODO: change this
		backendPool,
		testDefaultBeNodePort,
		namer,
	)
	cm := &ClusterManager{
		ClusterNamer: namer,
		instancePool: nodePool,
		backendPool:  backendPool,
		l7Pool:       l7Pool,
	}
	return &fakeClusterManager{cm, fakeLbs, fakeBackends, fakeIGs}
}
