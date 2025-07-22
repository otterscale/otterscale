<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import { Arc, Svg, Group, Chart, Text } from 'layerchart';
	import { cn } from '$lib/utils';
	import { metricColor, metricBackgroundColor } from '../../utils';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Empty from '../../utils/empty';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();

	const query = $derived(
		`
avg(
    sum by (instance) (
      irate(
        node_cpu_seconds_total{instance!~".*lxd.*",instance!~"juju.*",job=~".*",juju_application=~".*",juju_model=~".*",juju_model_uuid="${scope.uuid}",juju_unit=~".*",mode!="idle"}[4m]
      )
    )
  / on (instance) group_left ()
    sum by (instance) (
      (
        irate(
          node_cpu_seconds_total{instance!~".*lxd.*",instance!~"juju.*",job=~".*",juju_application=~".*",juju_model=~".*",juju_model_uuid="${scope.uuid}",juju_unit=~".*"}[4m]
        )
      )
    )
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
