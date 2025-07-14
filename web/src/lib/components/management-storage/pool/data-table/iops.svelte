<script lang="ts" module>
	import { Sparkline } from '$lib/components/custom/loading/index';
	import { fetchRange, integrateSerieses } from '$lib/components/dashboard/utils';
	import { formatNetworkIO } from '$lib/formatter';
	import Icon from '@iconify/svelte';
	import { format } from 'date-fns';
	import { LineChart, Tooltip } from 'layerchart';
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

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;
	let serieses: Map<string, SampleValue[] | undefined> = $state(new Map());
	let isMounted = $state(false);

	const readQuery = $derived(
		`
        rate(ceph_pool_rd{juju_model_uuid=~"${selectedScope}"}[4m]) * on (pool_id) group_left (instance, name) {juju_model_uuid=~"${selectedScope}",name=~"${selectedPool}"}
		`
	);
	const writeQuery = $derived(
		`
        rate(ceph_pool_wr{juju_model_uuid=~"${selectedScope}"}[4m]) * on (pool_id) group_left (instance, name) {juju_model_uuid=~"${selectedScope}",name=~"${selectedPool}"}
		`
	);

	onMount(async () => {
		try {
			await fetchRange(client, timeRange, step, readQuery).then((response) => {
				if (response && response.length > 0) {
					serieses.set('read', response);
				}
			});
			await fetchRange(client, timeRange, step, writeQuery).then((response) => {
				if (response && response.length > 0) {
					serieses.set('write', response);
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
			series={[
				{ key: 'read', color: 'hsl(var(--color-info))' },
				{ key: 'write', color: 'hsl(var(--color-danger))' }
			]}
			axis={false}
			grid={false}
			{renderContext}
			{debug}
		>
			<svelte:fragment slot="tooltip">
				<Tooltip.Root let:data>
					<Tooltip.Header class="bg-muted flex items-center justify-between rounded p-2">
						<Icon icon="ph:network" />
						{format(data.time, 'yyyy-MM-dd HH:mm')}
					</Tooltip.Header>
					<Tooltip.List class="rounded px-2">
						<Tooltip.Item label="read">
							{@const { value, unit } = formatNetworkIO(data.read)}
							<span class="flex items-center justify-end gap-1">
								<p>{value} {unit}</p>
							</span>
						</Tooltip.Item>
						<Tooltip.Item label="write">
							{@const { value, unit } = formatNetworkIO(data.write)}
							<span class="flex items-center justify-end gap-1">
								<p>{value} {unit}</p>
							</span>
						</Tooltip.Item>
					</Tooltip.List>
				</Tooltip.Root>
			</svelte:fragment>
		</LineChart>
	</div>
{/if}
