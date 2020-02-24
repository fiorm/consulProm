package discovery

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/testutil"
	common_config "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	sd_config "github.com/prometheus/prometheus/discovery/config"
	"github.com/prometheus/prometheus/discovery/consul"
	"github.com/prometheus/prometheus/discovery/file"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"gopkg.in/yaml.v2"
)


func TestGaugeFailedConfigs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	discoveryManager := NewManager(ctx, log.NewNopLogger())
	discoveryManager.updatert = 100 * time.Millisecond
	go discoveryManager.Run()

	c := map[string]sd_config.ServiceDiscoveryConfig{
		"prometheus": sd_config.ServiceDiscoveryConfig{
			ConsulSDConfigs: []*consul.SDConfig{
				&consul.SDConfig{
					Server: "foo:8500",
					TLSConfig: common_config.TLSConfig{
						CertFile: "/tmp/non_existent",
					},
				},
				&consul.SDConfig{
					Server: "bar:8500",
					TLSConfig: common_config.TLSConfig{
						CertFile: "/tmp/non_existent",
					},
				},
				&consul.SDConfig{
					Server: "foo2:8500",
					TLSConfig: common_config.TLSConfig{
						CertFile: "/tmp/non_existent",
					},
				},
			},
		},
	}
	discoveryManager.ApplyConfig(c)
	<-discoveryManager.SyncCh()

	failedCount := testutil.ToFloat64(failedConfigs)
	if failedCount != 3 {
		t.Fatalf("Expected to have 3 failed configs, got: %v", failedCount)
	}

	c["prometheus"] = sd_config.ServiceDiscoveryConfig{
		StaticConfigs: []*targetgroup.Group{
			&targetgroup.Group{
				Source: "0",
				Targets: []model.LabelSet{
					model.LabelSet{
						model.AddressLabel: "foo:9090",
					},
				},
			},
		},
	}
	discoveryManager.ApplyConfig(c)
	<-discoveryManager.SyncCh()

	failedCount = testutil.ToFloat64(failedConfigs)
	if failedCount != 0 {
		t.Fatalf("Expected to get no failed config, got: %v", failedCount)
	}

}
