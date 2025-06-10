<script lang="ts">
	import PageLoading from '$lib/components/otterscale/ui/page-loading.svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { onMount } from 'svelte';
	import Dashboard from './index.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';

	const prometheus = new PrometheusDriver({
		endpoint: 'http://10.102.197.175/microk8s-cos-prometheus-0',
		baseURL: '/api/v1'
	});

	let scopes: Scope[] = $state([]);
	async function fetchJujuModelUuids() {
		const query = `group by (juju_model_uuid) (node_exporter_build_info{})`;

		try {
			const response = await prometheus.instantQuery(query);

			scopes = response.result
				.map((result) => {
					return {
						uuid: result.metric.labels.juju_model_uuid,
						name: result.metric.labels.juju_model_uuid
					} as Scope;
				})
				.reverse();
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchJujuModelUuids();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if !mounted}
	<PageLoading />
{:else}
	<Dashboard client={prometheus} {scopes} />
{/if}
