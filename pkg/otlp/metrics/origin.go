// Copyright The OpenTelemetry Authors
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

package metrics

import (
	"fmt"
	"strings"
)

// OriginProduct defines the origin product.
type OriginProduct int32

const (
	OriginProductUnknown OriginProduct = iota
	// OriginProductDatadogAgent is the origin for metrics coming from the Datadog Agent OTLP Ingest.
	OriginProductDatadogAgent
	// OriginProductDatadogExporter is the origin for metrics coming from the OTel Collector Datadog Exporter.
	OriginProductDatadogExporter
)

func (o OriginProduct) String() string {
	switch o {
	case OriginProductUnknown:
		return "unknown"
	case OriginProductDatadogAgent:
		return "datadog-agent"
	case OriginProductDatadogExporter:
		return "datadog-exporter"
	default:
		return fmt.Sprintf("OriginProduct(%d)", o)
	}
}

// OriginCategory defines the origin category.
type OriginCategory int32

// OriginCategoryOTLP is the origin category for all metrics coming from OTLP.
// All metrics produced by the translator MUST have origin category set to OTLP.
const OriginCategoryOTLP OriginCategory = 0

func (o OriginCategory) String() string {
	switch o {
	case OriginCategoryOTLP:
		return "otlp"
	default:
		return fmt.Sprintf("OriginCategory(%d)", o)
	}
}

// OriginService defines the origin service.
type OriginService int32

// List all receivers that set the scope name.
const (
	OriginServiceUnknown OriginService = iota
	OriginServiceActiveDirectoryDSReceiver
	OriginServiceAerospikeReceiver
	OriginServiceApacheReceiver
	OriginServiceApacheSparkReceiver
	OriginServiceAzureMonitorReceiver
	OriginServiceBigIPReceiver
	OriginServiceChronyReceiver
	OriginServiceCouchDBReceiver
	OriginServiceDockerStatsReceiver
	OriginServiceElasticsearchReceiver
	OriginServiceExpVarReceiver
	OriginServiceFileStatsReceiver
	OriginServiceFlinkMetricsReceiver
	OriginServiceGitProviderReceiver
	OriginServiceHAProxyReceiver
	OriginServiceHostMetricsReceiver
	OriginServiceHTTPCheckReceiver
	OriginServiceIISReceiver
	OriginServiceK8SClusterReceiver
	OriginServiceKafkaMetricsReceiver
	OriginServiceKubeletStatsReceiver
	OriginServiceMemcachedReceiver
	OriginServiceMongoDBAtlasReceiver
	OriginServiceMongoDBReceiver
	OriginServiceMySQLReceiver
	OriginServiceNginxReceiver
	OriginServiceNSXTReceiver
	OriginServiceOracleDBReceiver
	OriginServicePostgreSQLReceiver
	OriginServicePrometheusReceiver
	OriginServiceRabbitMQReceiver
	OriginServiceRedisReceiver
	OriginServiceRiakReceiver
	OriginServiceSAPHANAReceiver
	OriginServiceSNMPReceiver
	OriginServiceSnowflakeReceiver
	OriginServiceSplunkEnterpriseReceiver
	OriginServiceSQLServerReceiver
	OriginServiceSSHCheckReceiver
	OriginServiceStatsDReceiver
	OriginServiceVCenterReceiver
	OriginServiceZookeeperReceiver
)

func originServiceFromScopeName(scopeName string) OriginService {
	const collectorPrefix = "otelcol/"
	if !strings.HasPrefix(scopeName, collectorPrefix) {
		return OriginServiceUnknown
	}

	// otelcol/kubeletstatsreceiver -> kubeletstatsreceiver
	// otelcol/hostmetricsreceiver/disk -> hostmetricsreceiver
	receiverName := strings.Split(scopeName, "/")[1]

	// otelcol
	switch receiverName {
	case "activedirectorydsreceiver":
		return OriginServiceActiveDirectoryDSReceiver
	case "aerospikereceiver":
		return OriginServiceAerospikeReceiver
	case "apachereceiver":
		return OriginServiceApacheReceiver
	case "apachesparkreceiver":
		return OriginServiceApacheSparkReceiver
	case "azuremonitorreceiver":
		return OriginServiceAzureMonitorReceiver
	case "bigipreceiver":
		return OriginServiceBigIPReceiver
	case "chronyreceiver":
		return OriginServiceChronyReceiver
	case "couchdbreceiver":
		return OriginServiceCouchDBReceiver
	case "dockerstatsreceiver":
		return OriginServiceDockerStatsReceiver
	case "elasticsearchreceiver":
		return OriginServiceElasticsearchReceiver
	case "expvarreceiver":
		return OriginServiceExpVarReceiver
	case "filestatsreceiver":
		return OriginServiceFileStatsReceiver
	case "flinkmetricsreceiver":
		return OriginServiceFlinkMetricsReceiver
	case "gitproviderreceiver":
		return OriginServiceGitProviderReceiver
	case "haproxyreceiver":
		return OriginServiceHAProxyReceiver
	case "hostmetricsreceiver":
		return OriginServiceHostMetricsReceiver
	case "httpcheckreceiver":
		return OriginServiceHTTPCheckReceiver
	case "iisreceiver":
		return OriginServiceIISReceiver
	case "k8sclusterreceiver":
		return OriginServiceK8SClusterReceiver
	case "kafkametricsreceiver":
		return OriginServiceKafkaMetricsReceiver
	case "kubeletstatsreceiver":
		return OriginServiceKubeletStatsReceiver
	case "memcachedreceiver":
		return OriginServiceMemcachedReceiver
	case "mongodbatlasreceiver":
		return OriginServiceMongoDBAtlasReceiver
	case "mongodbreceiver":
		return OriginServiceMongoDBReceiver
	case "mysqlreceiver":
		return OriginServiceMySQLReceiver
	case "nginxreceiver":
		return OriginServiceNginxReceiver
	case "nsxtreceiver":
		return OriginServiceNSXTReceiver
	case "oracledbreceiver":
		return OriginServiceOracleDBReceiver
	case "postgresqlreceiver":
		return OriginServicePostgreSQLReceiver
	case "prometheusreceiver":
		return OriginServicePrometheusReceiver
	case "rabbitmqreceiver":
		return OriginServiceRabbitMQReceiver
	case "redisreceiver":
		return OriginServiceRedisReceiver
	case "riakreceiver":
		return OriginServiceRiakReceiver
	case "saphanareceiver":
		return OriginServiceSAPHANAReceiver
	case "snmpreceiver":
		return OriginServiceSNMPReceiver
	case "snowflakereceiver":
		return OriginServiceSnowflakeReceiver
	case "splunkenterprisereceiver":
		return OriginServiceSplunkEnterpriseReceiver
	case "sqlserverreceiver":
		return OriginServiceSQLServerReceiver
	case "sshcheckreceiver":
		return OriginServiceSSHCheckReceiver
	case "statsdreceiver":
		return OriginServiceStatsDReceiver
	case "vcenterreceiver":
		return OriginServiceVCenterReceiver
	case "zookeeperreceiver":
		return OriginServiceZookeeperReceiver
	}

	return OriginServiceUnknown
}
