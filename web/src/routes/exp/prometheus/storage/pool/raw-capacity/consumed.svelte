<script lang="ts">
	import { PrometheusDriver, InstantVector } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../../utils/empty.svelte';
	import { Arc, Svg, Group, Chart, Text } from 'layerchart';

	import { formatCapacity } from '$lib/formatter';
	import { onMount } from 'svelte';
	import { metricBackgroundColor, metricColor } from '../../..';
	import { cn } from '$lib/utils';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();

	const totalQuery = $derived(
		`
		sum(ceph_osd_stat_bytes{juju_model_uuid=~"${scope.uuid}"})
		`
	);
	const consumedQuery = $derived(
		`
		sum(ceph_pool_bytes_used{juju_model_uuid=~"${scope.uuid}"})
		`
	);

	async function fetch(query: string): Promise<InstantVector | undefined | null> {
		try {
			const response = await client.instantQuery(query);
			const results = response.result;

			if (results.length === 0) {
				return undefined;
			}

			const [result] = results;
			return result;
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let totalResponse: InstantVector | undefined | null = $state();
	let consumedResponse: InstantVector | undefined | null = $state();

	let mounted = $state(false);
	onMount(async () => {
		try {
			totalResponse = await fetch(totalQuery);
			consumedResponse = await fetch(consumedQuery);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !mounted}
	<ComponentLoading />
{:else if !totalResponse || !consumedResponse}
	<NoData />
{:else}
	{@const total = formatCapacity(totalResponse.value.value / 1024 / 1024)}
	{@const consumed = formatCapacity(consumedResponse.value.value / 1024 / 1024)}
	{@const value = (consumedResponse.value.value / totalResponse.value.value) * 100}
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
				{`${consumed.value}${consumed.unit} / ${total.value}${total.unit}`}
			</p>
		</div>
	</div>
{/if}
