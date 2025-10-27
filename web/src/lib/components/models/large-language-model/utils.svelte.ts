import type { RangeVector, SampleValue } from 'prometheus-query';
import { SvelteMap } from 'svelte/reactivity';

import type { Application } from '$lib/api/application/v1/application_pb';

function getGatewayURL(applications: Application[]) {
	const vllmGateway = applications.find(
		(application) => application.labels['app.kubernetes.io/name'] == 'llm-d-infra',
	);
	if (!vllmGateway) {
		throw new Error('No Gateway');
	}

	const publicAddress = vllmGateway.publicAddress;
	const NodePortService = vllmGateway.services.find((service) => service.type == 'NodePort');
	if (!publicAddress || !NodePortService) {
		throw new Error('No NodePort Service');
	}

	const exposedPort = NodePortService.ports.find((port) => port.name == 'default');
	if (!exposedPort) {
		throw new Error('No Port');
	}

	return `http://10.102.197.145:10880/v1/models`;
	// return `${publicAddress}:${exposedPort.port}/llm-d/v1/models`;
}

function getMetricsMap(vectors: RangeVector[]) {
	return new SvelteMap(
		vectors.map((vector) => [
			(vector.metric.labels as { model_name?: string }).model_name,
			vector.values as SampleValue[],
		]),
	);
}

export { getGatewayURL, getMetricsMap };
