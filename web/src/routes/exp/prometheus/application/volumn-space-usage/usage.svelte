<script lang="ts">
	import { PrometheusDriver, InstantVector } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../utils/empty.svelte';
	import { Arc, Svg, Group, Chart, Text } from 'layerchart';

	import { formatCapacity } from '$lib/formatter';
	import { onMount } from 'svelte';
	import { metricBackgroundColor, metricColor, fetchInstance } from '../../utils';
	import { cn } from '$lib/utils';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();

	const usageQuery = $derived(
		`
		max without (instance, node) (
			(
				topk(
					1,
					kubelet_volume_stats_capacity_bytes{job="kubelet",juju_model_uuid=~"${scope.uuid}",metrics_path="/metrics"}
				)
				-
				topk(
					1,
					kubelet_volume_stats_available_bytes{job="kubelet",juju_model_uuid=~"${scope.uuid}",metrics_path="/metrics"}
				)
			)
			/
			topk(
				1,
				kubelet_volume_stats_capacity_bytes{job="kubelet",juju_model_uuid=~"${scope.uuid}",metrics_path="/metrics"}
			)
		*
			100
		)
		`
	);

	let usageResponse: InstantVector | undefined | null = $state();

	let mounted = $state(false);
	onMount(async () => {
		try {
			usageResponse = await fetchInstance(client, usageQuery);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !mounted}
	<ComponentLoading />
{:else if !usageResponse}
	<NoData />
{:else}
	{@const value = usageResponse.value.value * 100}
	{@const radius = 100}
	{@const border = radius * 2}
	<div class="flex h-full w-full items-center justify-center">
		<div class={cn(`h-[${border}px] w-[${border}px]`)}>
			<Chart>
				<Svg center>
					<Group>
						<Arc
							{value}
							domain={[0, 100]}
							outerRadius={radius}
							innerRadius={-13}
							cornerRadius={13}
							range={[-120, 120]}
							class={metricColor(value)}
							track={{ class: metricBackgroundColor(value) }}
						/>
					</Group>
				</Svg>
			</Chart>
		</div>
		<div class="absolute">
			<p class="text-xl">{`${value.toFixed(2)}%`}</p>
		</div>
	</div>
{/if}
