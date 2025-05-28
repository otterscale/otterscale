<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import { Arc, Svg, Group, Chart, Text } from 'layerchart';
	import { cn } from '$lib/utils';
	import { metricColor, metricBackgroundColor } from '../..';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../utils/empty.svelte';

	let {
		client,
		scope: scope,
		instance
	}: { client: PrometheusDriver; scope: Scope; instance: string } = $props();

	const query = $derived(
		`
		(
			(
				node_memory_SwapTotal_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
			-
				node_memory_SwapFree_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
			)
		/
			(node_memory_SwapTotal_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"})
		)
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const result = response.result}
	{#if result.length === 0}
		<NoData type="gauge" />
	{:else}
		{@const rawSwapUsage = result[0].value.value}
		<div class="flex h-full w-full items-center justify-center">
			<div class={cn(`h-[173px] w-[173px]`)}>
				<Chart>
					<Svg center>
						<Group y={100 / 4}>
							<Arc
								value={rawSwapUsage * 100}
								domain={[0, 100]}
								outerRadius={100}
								innerRadius={-13}
								cornerRadius={13}
								range={[-120, 120]}
								class={metricColor(rawSwapUsage * 100)}
								track={{ class: metricBackgroundColor(rawSwapUsage * 100) }}
								let:value
							>
								<Text
									value={!isNaN(value) ? `${value.toFixed(2)}%` : 'NaN'}
									textAnchor="middle"
									verticalAnchor="middle"
									class="text-xl tabular-nums"
								/>
							</Arc>
						</Group>
					</Svg>
				</Chart>
			</div>
		</div>
	{/if}
{:catch error}
	Error
{/await}
