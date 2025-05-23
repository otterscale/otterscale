<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import { Arc, Svg, Group, Chart, Text } from 'layerchart';
	import { cn } from '$lib/utils';
	import { metricColor, metricBackgroundColor } from '..';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';

	let {
		client,
		juju_model_uuid,
		instance
	}: { client: PrometheusDriver; juju_model_uuid: string; instance: string } = $props();

	const query = $derived(
		`
		1
		-
		(
			(
				node_filesystem_avail_bytes{fstype!="rootfs",instance="${instance}",juju_model_uuid=~"${juju_model_uuid}",mountpoint="/"}
			)
			/
			node_filesystem_size_bytes{fstype!="rootfs",instance="${instance}",juju_model_uuid=~"${juju_model_uuid}",mountpoint="/"}
		)
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const rawSwapUsage = response.result[0].value.value}
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
{:catch error}
	Error
{/await}
