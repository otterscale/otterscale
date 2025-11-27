package model

import (
	"strconv"
	"strings"
)

// v1.2.0-rc.1
func convertGAIEValuesMap(mode Mode, name, releaseName string, port int32) map[string]string {
	ret := map[string]string{
		"inferenceExtension.image.name":                                                    "llm-d-inference-scheduler",
		"inferenceExtension.image.hub":                                                     "ghcr.io/llm-d",
		"inferenceExtension.image.tag":                                                     "v0.4.0-rc.1",
		"inferenceExtension.monitoring.prometheus.enabled":                                 "true",
		"inferenceExtension.monitoring.prometheus.auth.secretName":                         name + "-sa-metrics-reader-secret",
		"inferencePool.targetPortNumber":                                                   strconv.FormatInt(int64(port), 10),
		"inferencePool.modelServers.matchLabels." + escapeDot("llm-d.ai/inferenceServing"): "true",
		"inferencePool.modelServers.matchLabels." + escapeDot("llm-d.ai/model"):            releaseName,
		"provider.name": "istio",
		"provider.istio.destinationRule.trafficPolicy.tls.insecureSkipVerify": "true",
		"provider.istio.destinationRule.trafficPolicy.tls.mode":               "SIMPLE",
	}

	switch mode {
	case ModeIntelligentInferenceScheduling:
		ret["inferenceExtension.pluginsConfigFile"] = "default-plugins.yaml"

	case ModePrefillDecodeDisaggregation:
		ret["inferenceExtension.pluginsConfigFile"] = "pd-config.yaml"
		ret["inferenceExtension.pluginsCustomConfig."+escapeDot("pd-config.yaml")] = `# ALWAYS DO PD IN THIS EXAMPLE (THRESHOLD 0)
apiVersion: inference.networking.x-k8s.io/v1alpha1
kind: EndpointPickerConfig
plugins:
- type: prefill-header-handler
- type: prefill-filter
- type: decode-filter
- type: max-score-picker
- type: queue-scorer
  parameters:
    hashBlockSize: 5
    maxPrefixBlocksToMatch: 256
    lruCapacityPerServer: 31250
- type: pd-profile-handler
  parameters:
    threshold: 0
    hashBlockSize: 5
schedulingProfiles:
- name: prefill
  plugins:
  - pluginRef: prefill-filter
  - pluginRef: queue-scorer
    weight: 1.0
  - pluginRef: max-score-picker
- name: decode
  plugins:
  - pluginRef: decode-filter
  - pluginRef: queue-scorer
    weight: 1.0
  - pluginRef: max-score-picker`
	}

	return ret
}

// v0.3.10
//
//nolint:funlen // ignore
func convertModelServiceValuesMap(mode Mode, releaseName, modelName string, sizeBytes uint64, prefill *Prefill, decode *Decode, maxModelLength uint32) map[string]string {
	ret := map[string]string{
		"modelArtifacts.name": modelName,
		"modelArtifacts.labels." + escapeDot("llm-d.ai/model"): releaseName,
		"modelArtifacts.size":                                      "[str]" + strconv.FormatUint(sizeBytes, 10),
		"accelerator.type":                                         "nvidia",
		"accelerator.resources.nvidia":                             vgpuResource,
		"routing.servicePort":                                      "8000",
		"routing.proxy.enable":                                     "true",
		"routing.proxy.image":                                      "ghcr.io/llm-d/llm-d-routing-sidecar:v0.4.0-rc.1",
		"routing.proxy.connector":                                  "nixlv2",
		"routing.proxy.secure":                                     "false",
		"decode.create":                                            "true",
		"decode.autoscaling.enabled":                               "true",
		"decode.replicas":                                          strconv.FormatUint(uint64(decode.Replica), 10),
		"decode.containers[0].image":                               "ghcr.io/llm-d/llm-d-cuda:v0.3.1",
		"decode.containers[0].modelCommand":                        "vllmServe",
		"decode.containers[0].args[0]":                             "--kv-transfer-config",
		"decode.containers[0].args[1]":                             "\\{\"kv_connector\":\"NixlConnector\"\\, \"kv_role\":\"kv_both\"\\}",
		"decode.containers[0].args[2]":                             "--disable-uvicorn-access-log",
		"decode.containers[0].args[3]":                             "--max-model-len",
		"decode.containers[0].args[4]":                             "[str]" + strconv.FormatUint(uint64(maxModelLength), 10),
		"decode.containers[0].env[0].name":                         "VLLM_NIXL_SIDE_CHANNEL_HOST",
		"decode.containers[0].env[0].valueFrom.fieldRef.fieldPath": "status.podIP",
		"decode.containers[0].env[1].name":                         "VLLM_LOGGING_LEVEL",
		"decode.containers[0].env[1].value":                        "INFO",
		"decode.containers[0].ports[0].containerPort":              "8200",
		"decode.containers[0].ports[0].name":                       "metrics",
		"decode.containers[0].ports[0].protocol":                   "TCP",
		"decode.containers[0].resources.limits." + escapeDot(vgpuMemPercentageResource):   strconv.FormatUint(uint64(decode.VGPUMemory), 10),
		"decode.containers[0].resources.requests." + escapeDot(vgpuMemPercentageResource): strconv.FormatUint(uint64(decode.VGPUMemory), 10),
		"decode.containers[0].mountModelVolume":                                           "true",
		"decode.containers[0].volumeMounts[0].name":                                       "metrics-volume",
		"decode.containers[0].volumeMounts[0].mountPath":                                  "/.config",
		"decode.containers[0].volumeMounts[1].name":                                       "torch-compile-cache",
		"decode.containers[0].volumeMounts[1].mountPath":                                  "/.cache",
		"decode.containers[0].startupProbe.httpGet.path":                                  "/v1/models",
		"decode.containers[0].startupProbe.httpGet.port":                                  "8200",
		"decode.containers[0].startupProbe.initialDelaySeconds":                           "15",
		"decode.containers[0].startupProbe.periodSeconds":                                 "30",
		"decode.containers[0].startupProbe.timeoutSeconds":                                "5",
		"decode.containers[0].startupProbe.failureThreshold":                              "60",
		"decode.containers[0].livenessProbe.httpGet.path":                                 "/health",
		"decode.containers[0].livenessProbe.httpGet.port":                                 "8200",
		"decode.containers[0].livenessProbe.periodSeconds":                                "10",
		"decode.containers[0].livenessProbe.timeoutSeconds":                               "5",
		"decode.containers[0].livenessProbe.failureThreshold":                             "3",
		"decode.containers[0].readinessProbe.httpGet.path":                                "/v1/models",
		"decode.containers[0].readinessProbe.httpGet.port":                                "8200",
		"decode.containers[0].readinessProbe.periodSeconds":                               "5",
		"decode.containers[0].readinessProbe.timeoutSeconds":                              "2",
		"decode.containers[0].readinessProbe.failureThreshold":                            "3",
		"decode.volumes[0].name":                                                          "metrics-volume",
		"decode.volumes[1].name":                                                          "torch-compile-cache",
		"decode.monitoring.podmonitor.enabled":                                            "true",
		"decode.monitoring.podmonitor.portName":                                           "metrics",
		"decode.monitoring.podmonitor.path":                                               "/metrics",
		"decode.monitoring.podmonitor.interval":                                           "30s",
	}

	switch mode {
	case ModeIntelligentInferenceScheduling:
		ret["decode.containers[0].env[2].name"] = "CUDA_VISIBLE_DEVICES"
		ret["decode.containers[0].env[2].value"] = "[str]" + strconv.FormatUint(0, 10)
		ret["decode.containers[0].env[3].name"] = "UCX_TLS"
		ret["decode.containers[0].env[3].value"] = "cuda_ipc\\,cuda_copy\\,tcp"
		ret["decode.containers[0].env[4].name"] = "VLLM_NIXL_SIDE_CHANNEL_PORT"
		ret["decode.containers[0].env[4].value"] = "[str]" + strconv.FormatUint(5557, 10)
		ret["decode.containers[0].ports[1].containerPort"] = "5557"
		ret["decode.containers[0].ports[1].protocol"] = "TCP"
		ret["prefill.create"] = "false"

	case ModePrefillDecodeDisaggregation:
		ret["decode.parallelism.tensor"] = strconv.FormatUint(uint64(decode.Tensor), 10)
		ret["decode.containers[0].volumeMounts[2].name"] = "shm"
		ret["decode.containers[0].volumeMounts[2].mountPath"] = "/dev/shm"
		ret["decode.volumes[2].name"] = "shm"
		ret["decode.volumes[2].emptyDir.medium"] = "Memory"
		ret["decode.volumes[2].emptyDir.sizeLimit"] = "16Gi"
		ret["prefill.replicas"] = strconv.FormatUint(uint64(prefill.Replica), 10)
		ret["prefill.containers[0].image"] = "ghcr.io/llm-d/llm-d-cuda:v0.3.1"
		ret["prefill.containers[0].modelCommand"] = "vllmServe"
		ret["prefill.containers[0].args[0]"] = "--kv-transfer-config"
		ret["prefill.containers[0].args[1]"] = "\\{\"kv_connector\":\"NixlConnector\"\\, \"kv_role\":\"kv_both\"\\}"
		ret["prefill.containers[0].args[2]"] = "--disable-uvicorn-access-log"
		ret["prefill.containers[0].args[3]"] = "--max-model-len"
		ret["prefill.containers[0].args[4]"] = "[str]" + strconv.FormatUint(uint64(maxModelLength), 10)
		ret["prefill.containers[0].env[0].name"] = "VLLM_NIXL_SIDE_CHANNEL_HOST"
		ret["prefill.containers[0].env[0].valueFrom.fieldRef.fieldPath"] = "status.podIP"
		ret["prefill.containers[0].env[1].name"] = "VLLM_LOGGING_LEVEL"
		ret["prefill.containers[0].env[1].value"] = "INFO"
		ret["prefill.containers[0].ports[0].containerPort"] = "8000"
		ret["prefill.containers[0].ports[0].name"] = "metrics"
		ret["prefill.containers[0].ports[0].protocol"] = "TCP"
		ret["prefill.containers[0].resources.limits."+escapeDot(vgpuMemPercentageResource)] = strconv.FormatUint(uint64(prefill.VGPUMemory), 10)
		ret["prefill.containers[0].resources.requests."+escapeDot(vgpuMemPercentageResource)] = strconv.FormatUint(uint64(prefill.VGPUMemory), 10)
		ret["prefill.containers[0].mountModelVolume"] = "true"
		ret["prefill.containers[0].volumeMounts[0].name"] = "metrics-volume"
		ret["prefill.containers[0].volumeMounts[0].mountPath"] = "/.config"
		ret["prefill.containers[0].volumeMounts[1].name"] = "torch-compile-cache"
		ret["prefill.containers[0].volumeMounts[1].mountPath"] = "/.cache"
		ret["prefill.containers[0].volumeMounts[2].name"] = "shm"
		ret["prefill.containers[0].volumeMounts[2].mountPath"] = "/dev/shm"
		ret["prefill.containers[0].startupProbe.httpGet.path"] = "/v1/models"
		ret["prefill.containers[0].startupProbe.httpGet.port"] = "8000"
		ret["prefill.containers[0].startupProbe.initialDelaySeconds"] = "15"
		ret["prefill.containers[0].startupProbe.periodSeconds"] = "30"
		ret["prefill.containers[0].startupProbe.timeoutSeconds"] = "5"
		ret["prefill.containers[0].startupProbe.failureThreshold"] = "60"
		ret["prefill.containers[0].livenessProbe.httpGet.path"] = "/health"
		ret["prefill.containers[0].livenessProbe.httpGet.port"] = "8000"
		ret["prefill.containers[0].livenessProbe.periodSeconds"] = "10"
		ret["prefill.containers[0].livenessProbe.timeoutSeconds"] = "5"
		ret["prefill.containers[0].livenessProbe.failureThreshold"] = "3"
		ret["prefill.containers[0].readinessProbe.httpGet.path"] = "/v1/models"
		ret["prefill.containers[0].readinessProbe.httpGet.port"] = "8000"
		ret["prefill.containers[0].readinessProbe.periodSeconds"] = "5"
		ret["prefill.containers[0].readinessProbe.timeoutSeconds"] = "2"
		ret["prefill.containers[0].readinessProbe.failureThreshold"] = "3"
		ret["prefill.volumes[0].name"] = "metrics-volume"
		ret["prefill.volumes[1].name"] = "torch-compile-cache"
		ret["prefill.volumes[2].name"] = "shm"
		ret["prefill.volumes[2].emptyDir.medium"] = "Memory"
		ret["prefill.volumes[2].emptyDir.sizeLimit"] = "16Gi"
		ret["prefill.monitoring.podmonitor.enabled"] = "true"
		ret["prefill.monitoring.podmonitor.portName"] = "metrics"
		ret["prefill.monitoring.podmonitor.path"] = "/metrics"
		ret["prefill.monitoring.podmonitor.interval"] = "30s"
	}

	return ret
}

func escapeDot(key string) string {
	return strings.ReplaceAll(key, ".", "\\.")
}
