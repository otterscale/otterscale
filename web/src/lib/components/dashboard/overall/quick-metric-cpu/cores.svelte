<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Empty from '../../utils/empty';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
count(
  count by (instance, cpu) (
    node_cpu_seconds_total{instance!~".*lxd.*",instance!~"juju.*",job=~".*",juju_application=~".*",juju_model=~".*",juju_model_uuid="${scope.uuid}",juju_unit=~".*"}
  )
)
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const result = response.result}
	{#if result.length === 0}
		<Empty.Text />
	{:else}
		{@const cores = result[0].value.value}
		<p class="text-3xl">{cores} Cores</p>
	{/if}
{:catch error}
	Error
{/await}
