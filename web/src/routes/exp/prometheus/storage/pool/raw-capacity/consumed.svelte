<script lang="ts" module>
	import { InstantVector } from 'prometheus-query';
</script>

<script lang="ts">
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../../utils/empty.svelte';
	import { Arc, Svg, Group, Chart, Text } from 'layerchart';

	import { formatCapacity } from '$lib/formatter';
	import Badge from '$lib/components/ui/badge/badge.svelte';
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
			console.log(totalResponse);
			console.log(consumedResponse);

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
	<div class="flex h-full w-full items-center justify-center">
		<div class={cn(`h-[173px] w-[173px]`)}>
			<Chart>
				<Svg center>
					<Group y={radius / 4}>
						<Arc
							{value}
							domain={[0, 100]}
							outerRadius={radius}
							innerRadius={-13}
							cornerRadius={13}
							range={[-120, 120]}
							class={metricColor(value)}
							track={{ class: metricBackgroundColor(value) }}
						>
							<Text
								value={`${value.toFixed(2)}%`}
								textAnchor="middle"
								verticalAnchor="middle"
								class="text-xl tabular-nums"
							/>
							<Text
								value={`${consumed.value}${consumed.unit} / ${total.value}${total.unit}`}
								y={radius / 4}
								textAnchor="middle"
								verticalAnchor="middle"
								class="text-xs font-extralight tabular-nums"
							/>
						</Arc>
					</Group>
				</Svg>
			</Chart>
		</div>
	</div>
{/if}
