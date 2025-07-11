<script lang="ts" module>
	import { fetchRange } from '$lib/components/dashboard/utils';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { LineChart } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onMount } from 'svelte';

	const client = new PrometheusDriver({
		endpoint: 'http://10.102.197.18/cos-dev-prometheus-0',
		baseURL: '/api/v1'
	});
	const step = 10 * 60;
	const timeRange = {
		start: new Date(Date.now() - 60 * 60 * 1000),
		end: new Date()
	};
</script>

<script lang="ts">
	let { selectedScope, selectedPool }: { selectedScope: string; selectedPool: string } = $props();

	let renderContext: 'svg' | 'canvas' = 'canvas';
	let debug = false;

	let series: SampleValue[] | undefined = $state([]);
	let mounted = $state(false);

	const query = $derived(
		`
        rate(ceph_pool_rd{juju_model_uuid=~"${selectedScope}"}[4m]) * on (pool_id) group_left (instance, name) {juju_model_uuid=~"${selectedScope}",name=~"${selectedPool}"}
		`
	);

	onMount(async () => {
		try {
			series = await fetchRange(client, timeRange, step, query);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !mounted}
	<ComponentLoading />
{:else if series && series.length > 0}
	<div class="h-[50px] w-[100px]">
		<LineChart data={series} x="time" y="value" {renderContext} {debug} />
	</div>
{/if}
