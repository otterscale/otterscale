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

	const totalQuery = $derived(
		`
		ceph_cluster_total_bytes{juju_model_uuid=~"${scope.uuid}"}
		`
	);
	const usedQuery = $derived(
		`
		ceph_cluster_total_used_bytes{juju_model_uuid=~"${scope.uuid}"}
		`
	);

	let totalResponse: InstantVector | undefined | null = $state();
	let usedResponse: InstantVector | undefined | null = $state();

	let mounted = $state(false);
	onMount(async () => {
		try {
			totalResponse = await fetchInstance(client, totalQuery);
			usedResponse = await fetchInstance(client, usedQuery);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !mounted}
	<ComponentLoading />
{:else if !totalResponse || !usedResponse}
	<NoData />
{:else}
	{@const total = formatCapacity(totalResponse.value.value / 1024 / 1024)}
	{@const used = formatCapacity(usedResponse.value.value / 1024 / 1024)}
	{@const value = (usedResponse.value.value / totalResponse.value.value) * 100}
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
		<div class="absolute">
			<p class="mt-10 text-xs font-extralight">
				{`${used.value}${used.unit} / ${total.value}${total.unit}`}
			</p>
		</div>
	</div>
{/if}
