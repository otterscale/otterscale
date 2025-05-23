<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import Dashboard from './index.svelte';
	import PageLoading from '$lib/components/otterscale/ui/page-loading.svelte';

	const juju_model_uuid = 'b62d195e-3905-4960-85ee-7673f71eb21e';
	const prometheus = new PrometheusDriver({
		endpoint: 'http://10.102.197.175/microk8s-cos-prometheus-0',
		baseURL: '/api/v1'
	});

	const query = `node_time_seconds{juju_model_uuid=~"${juju_model_uuid}"}`;
</script>

{#await prometheus.instantQuery(query)}
	<PageLoading />
{:then response}
	{@const results = response.result.sort((p, n) =>
		p.metric.labels.instance.localeCompare(n.metric.labels.instance)
	)}
	{@const instances = results.map((result) => result.metric.labels.instance)}
	<Dashboard client={prometheus} {juju_model_uuid} {instances} />
{:catch error}
	Error
{/await}
