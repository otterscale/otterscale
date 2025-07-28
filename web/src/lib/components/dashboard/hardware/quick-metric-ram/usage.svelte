<script lang="ts">
	import type { Machine } from '$gen/api/machine/v1/machine_pb';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { cn } from '$lib/utils';
	import { Arc, Chart, Group, Svg } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';
	import { metricBackgroundColor, metricColor } from '../../utils';
	import * as Empty from '../../utils/empty';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();

	const query = $derived(
		`
		1
		-
		(
			(node_memory_MemAvailable_bytes{instance=~"${machine.fqdn}"})
			/
			node_memory_MemTotal_bytes{instance=~"${machine.fqdn}"}
		)
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const results = response.result}
	{#if results.length === 0}
		<Empty.Gauge />
	{:else}
		{@const [result] = results}
		{@const usage = result.value.value * 100}
		{@const radius = 100}
		{@const border = radius * 2}
		<div class="flex h-full w-full items-center justify-center">
			<div class={cn(`h-[${border}px] w-[${border}px]`)}>
				<Chart>
					<Svg center>
						<Group>
							<Arc
								value={usage}
								domain={[0, 100]}
								outerRadius={100}
								innerRadius={-13}
								cornerRadius={13}
								range={[-120, 120]}
								class={metricColor(usage)}
								track={{ class: metricBackgroundColor(usage) }}
							/>
						</Group>
					</Svg>
				</Chart>
			</div>
			<div class="absolute">
				<p class="text-xl">{!isNaN(usage) ? `${usage.toFixed(2)}%` : 'NaN'}</p>
			</div>
		</div>
	{/if}
{:catch error}
	Error
{/await}
