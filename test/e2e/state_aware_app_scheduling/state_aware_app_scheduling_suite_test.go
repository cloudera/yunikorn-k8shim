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

package stateawareappscheduling_test

import (
	"path/filepath"
	"testing"

	v1 "k8s.io/api/core/v1"

	"github.com/apache/incubator-yunikorn-k8shim/test/e2e/framework/helpers/k8s"

	"github.com/onsi/ginkgo/reporters"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"github.com/apache/incubator-yunikorn-k8shim/test/e2e/framework/configmanager"
)

func init() {
	configmanager.YuniKornTestConfig.ParseFlags()
}

var k = k8s.KubeCtl{}
var oldConfigMap *v1.ConfigMap

var _ = BeforeSuite(func() {
	By("Enabling state aware app scheduling config over config maps")
	Ω(k.SetClient()).To(BeNil())
	var c, err = k.GetConfigMaps(configmanager.YuniKornTestConfig.YkNamespace,
		configmanager.DefaultYuniKornConfigMap)
	Ω(err).NotTo(HaveOccurred())
	Ω(c).NotTo(BeNil())

	oldConfigMap = c.DeepCopy()
	Ω(c).Should(BeEquivalentTo(oldConfigMap))
	var stateAwareStr = "partitions:\n  -\n    name: default\n    placementrules:\n      - name: tag\n        " +
		"value: namespace\n        create: true\n    queues:\n      - name: root\n        " +
		"submitacl: '*'\n        properties:\n          application.sort.policy: stateaware"
	c.Data[configmanager.DefaultPolicyGroup] = stateAwareStr
	var d, err2 = k.UpdateConfigMap(c, configmanager.YuniKornTestConfig.YkNamespace)
	Ω(err2).NotTo(HaveOccurred())
	Ω(d).NotTo(BeNil())
})

var _ = AfterSuite(func() {
	By("Restoring the old config maps")
	var c, err = k.GetConfigMaps(configmanager.YuniKornTestConfig.YkNamespace,
		configmanager.DefaultYuniKornConfigMap)
	Ω(err).NotTo(HaveOccurred())
	Ω(c).NotTo(BeNil())
	c.Data = oldConfigMap.Data
	var e, err3 = k.UpdateConfigMap(c, configmanager.YuniKornTestConfig.YkNamespace)
	Ω(err3).NotTo(HaveOccurred())
	Ω(e).NotTo(BeNil())
})

func TestStateAwareAppScheduling(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter(
		filepath.Join(configmanager.YuniKornTestConfig.LogDir, "StateAwareAppScheduling_junit.xml"))
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "TestStateAwareAppScheduling", []ginkgo.Reporter{junitReporter})
}

// Declarations for Ginkgo DSL
var Fail = ginkgo.Fail
var Describe = ginkgo.Describe
var It = ginkgo.It
var PIt = ginkgo.PIt
var By = ginkgo.By
var BeforeEach = ginkgo.BeforeEach
var AfterEach = ginkgo.AfterEach
var BeforeSuite = ginkgo.BeforeSuite
var AfterSuite = ginkgo.AfterSuite

// Declarations for Gomega DSL
var RegisterFailHandler = gomega.RegisterFailHandler

// Declarations for Gomega Matchers
var Equal = gomega.Equal
var Ω = gomega.Expect
var BeNil = gomega.BeNil
var HaveOccurred = gomega.HaveOccurred
var BeEmpty = gomega.BeEmpty
var BeTrue = gomega.BeTrue
var ContainSubstring = gomega.ContainSubstring
var BeEquivalentTo = gomega.BeEquivalentTo
