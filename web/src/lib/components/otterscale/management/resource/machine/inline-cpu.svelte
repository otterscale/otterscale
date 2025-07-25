<script lang="ts" module>
	import { Sparkline } from '$lib/components/custom/loading/index';
	import { fetchRange, integrateSerieses } from '$lib/components/dashboard/utils';
	import Icon from '@iconify/svelte';
	import { format } from 'date-fns';
	import { LineChart, Tooltip } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onMount } from 'svelte';

	const step = 10 * 60;
	const timeRange = {
		start: new Date(Date.now() - 60 * 60 * 1000),
		end: new Date()
	};
</script>

<script lang="ts">
	let {
		client,
		selectedMachine
	}: { client: PrometheusDriver; selectedMachine: string } = $props();

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;
	let serieses: Map<string, SampleValue[] | undefined> = $state(new Map());
	let isMounted = $state(false);

	const query = $derived(
		`
        avg by (instance) (rate(node_cpu_seconds_total{instance="${selectedMachine}", mode!="idle"}[5m]))
		`
	);

	onMount(async () => {
		try {
			await fetchRange(client, timeRange, step, query).then((response) => {
				if (response && response.length > 0) {
					serieses.set('cpu', response);
				}
			});

			isMounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !isMounted}
	<Sparkline class="w-[100px]" />
{:else if serieses.size > 0}
	{@const data = integrateSerieses(serieses)}
	<div class="h-[50px] w-[100px]">
		<LineChart
			{data}
			x="time"
			y="cpu"
			yDomain={[0, 1]}
			series={[
				{ key: 'cpu', color: 'hsl(var(--color-info))' }
			]}
			axis={false}
			grid={false}
			{renderContext}
			{debug}
		>
			<svelte:fragment slot="tooltip">
				<Tooltip.Root let:data>
					<Tooltip.Header>
						<span><Icon icon="ph:cpu" /></span>
						<span>{format(data.time, 'yyyy-MM-dd HH:mm')}</span>
					</Tooltip.Header>
					<Tooltip.List class="rounded px-2">
						<Tooltip.Item label="CPU" color="hsl(var(--color-info))">
							<p>{(data.cpu * 100).toFixed(2)}%</p>
						</Tooltip.Item>
					</Tooltip.List>
				</Tooltip.Root>
			</svelte:fragment>
		</LineChart>
	</div>
{/if}
