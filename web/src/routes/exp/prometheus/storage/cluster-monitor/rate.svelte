<script lang="ts">
	import { InstantVector, PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../utils/empty.svelte';
	import { onMount } from 'svelte';
	import { fetchInstance, metricBackgroundColor, metricColor } from '../../utils';
	import { cn } from '$lib/utils';
	import { Arc, Chart, Group, Svg } from 'layerchart';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const totalQuery = $derived(
		`
        count(ceph_mon_quorum_status{juju_model_uuid=~"${scope.uuid}"})
		`
	);
	const upQuery = $derived(
		`
        sum(ceph_mon_quorum_status{juju_model_uuid=~"${scope.uuid}"})
		`
	);

	let totalResponse: InstantVector | undefined | null = $state();
	let upResponse: InstantVector | undefined | null = $state();

	let mounted = $state(false);
	onMount(async () => {
		try {
			totalResponse = await fetchInstance(client, totalQuery);
			upResponse = await fetchInstance(client, upQuery);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !mounted}
	<ComponentLoading />
{:else if !totalResponse || !upResponse}
	<NoData />
{:else}
	{@const totalNumber = Number(totalResponse.value.value)}
	{@const upNumber = Number(upResponse.value.value)}
	{@const value = (upNumber * 100) / totalNumber}
	{@const reversedValue = 100 - value}
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
							class={metricColor(reversedValue)}
							track={{ class: metricBackgroundColor(reversedValue) }}
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
				{upNumber} / {totalNumber}
			</p>
		</div>
	</div>
{/if}
