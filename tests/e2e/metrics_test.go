// Copyright 2017 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"fmt"
	"strings"
	"testing"

	"go.etcd.io/etcd/version"
)

func TestV3MetricsSecure(t *testing.T) {
	cfg := configTLS
	cfg.clusterSize = 1
	cfg.metricsURLScheme = "https"
	testCtl(t, metricsTest)
}

func TestV3MetricsInsecure(t *testing.T) {
	cfg := configTLS
	cfg.clusterSize = 1
	cfg.metricsURLScheme = "http"
	testCtl(t, metricsTest)
}

func TestV3MetricsSecureTLSCertAuth(t *testing.T) {
	testCtl(t, metricsTestCertAuth, withCfg(configMetricsTLS))
}

func metricsTest(cx ctlCtx) {
	if err := ctlV3Put(cx, "k", "v", ""); err != nil {
		cx.t.Fatal(err)
	}
	if err := cURLGet(cx.epc, cURLReq{endpoint: "/metrics", expected: `etcd_debugging_mvcc_keys_total 1`, metricsURLScheme: cx.cfg.metricsURLScheme}); err != nil {
		cx.t.Fatalf("failed get with curl (%v)", err)
	}
	if err := cURLGet(cx.epc, cURLReq{endpoint: "/metrics", expected: fmt.Sprintf(`etcd_server_version{server_version="%s"} 1`, version.Version), metricsURLScheme: cx.cfg.metricsURLScheme}); err != nil {
		cx.t.Fatalf("failed get with curl (%v)", err)
	}
	ver := version.Version
	if strings.HasSuffix(ver, "+git") {
		ver = strings.Replace(ver, "+git", "", 1)
	}
	if err := cURLGet(cx.epc, cURLReq{endpoint: "/metrics", expected: fmt.Sprintf(`etcd_cluster_version{cluster_version="%s"} 1`, ver), metricsURLScheme: cx.cfg.metricsURLScheme}); err != nil {
		cx.t.Fatalf("failed get with curl (%v)", err)
	}
	if err := cURLGet(cx.epc, cURLReq{endpoint: "/health", expected: `{"health":"true"}`, metricsURLScheme: cx.cfg.metricsURLScheme}); err != nil {
		cx.t.Fatalf("failed get with curl (%v)", err)
	}
}

func metricsTestCertAuth(cx ctlCtx) {
	//	fmt.Printf("%#+v\n", cx.epc.procs[0].Config())
	//	fmt.Printf("\n")
	//	return

	if err := ctlV3Put(cx, "k", "v", ""); err != nil {
		cx.t.Fatal(err)
	}
	if err := cURLGet(cx.epc, cURLReq{endpoint: "/metrics", expected: `etcd_debugging_mvcc_keys_total 1`, metricsURLScheme: cx.cfg.metricsURLScheme, useCertAuth: true}); err != nil {
		cx.t.Fatalf("failed get with curl (%v)", err)
	}
	if err := cURLGet(cx.epc, cURLReq{endpoint: "/metrics", expected: fmt.Sprintf(`etcd_server_version{server_version="%s"} 1`, version.Version), metricsURLScheme: cx.cfg.metricsURLScheme, useCertAuth: true}); err != nil {
		cx.t.Fatalf("failed get with curl (%v)", err)
	}
	ver := version.Version
	if strings.HasSuffix(ver, "+git") {
		ver = strings.Replace(ver, "+git", "", 1)
	}
	if err := cURLGet(cx.epc, cURLReq{endpoint: "/metrics", expected: fmt.Sprintf(`etcd_cluster_version{cluster_version="%s"} 1`, ver), metricsURLScheme: cx.cfg.metricsURLScheme, useCertAuth: true}); err != nil {
		cx.t.Fatalf("failed get with curl (%v)", err)
	}
	if err := cURLGet(cx.epc, cURLReq{endpoint: "/health", expected: `{"health":"true"}`, metricsURLScheme: cx.cfg.metricsURLScheme, useCertAuth: true}); err != nil {
		cx.t.Fatalf("failed get with curl (%v)", err)
	}

	req := cURLReq{endpoint: "/metrics", metricsURLScheme: cx.cfg.metricsURLScheme}

	expectErr := []string{
		"curl: (60) SSL certificate problem: unable to get local issuer certificate",
		"More details here: https://curl.haxx.se/docs/sslcerts.html",
		"",
		"curl failed to verify the legitimacy of the server and therefore could not",
		"establish a secure connection to it. To learn more about this situation and",
		"how to fix it, please visit the web page mentioned above.",
	}

	if err := spawnWithExpects(cURLPrefixArgs(cx.epc, "GET", req), expectErr...); err != nil {
		cx.t.Fatalf("failed get with curl (%v)", err)
	}
}
