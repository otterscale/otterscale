<script lang="ts">
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import Dashboard from './index.svelte';
	import { getContext, onMount } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import PageLoading from '$lib/components/otterscale/ui/page-loading.svelte';
	import { ScopeService, type Scope } from '$gen/api/scope/v1/scope_pb';

	const juju_model_uuid = 'b62d195e-3905-4960-85ee-7673f71eb21e';
	const prometheus = new PrometheusDriver({
		endpoint: 'http://10.102.197.175/microk8s-cos-prometheus-0',
		baseURL: '/api/v1'
	});

	const transport: Transport = getContext('transport');
	const client = createClient(ScopeService, transport);

	let scopes: Scope[] = $state([]);
	async function fetchScopes() {
		// try {
		// 	const response = await client.listScopes({});
		// 	scopes = response.scopes;
		// } catch (error) {
		// 	console.error('Error fetching:', error);
		// }
		scopes = [
			{ uuid: 'b62d195e-3905-4960-85ee-7673f71eb21e', name: 'one' } as Scope,
			{ uuid: '66c3ec8b-1052-4586-8b5b-db0b2cb141ea', name: 'another' } as Scope
		];
	}

	let instances: string[] = $state([]);
	async function fetchInstances() {
		const query = `node_time_seconds{juju_model_uuid=~"${juju_model_uuid}"}`;

		try {
			const response = await prometheus.instantQuery(query);

			instances = response.result
				.sort((p, n) => p.metric.labels.instance.localeCompare(n.metric.labels.instance))
				.map((result) => result.metric.labels.instance);
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchScopes();
			await fetchInstances();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if !mounted}
	<PageLoading />
{:else}
	<Dashboard client={prometheus} {scopes} {instances} />
{/if}
