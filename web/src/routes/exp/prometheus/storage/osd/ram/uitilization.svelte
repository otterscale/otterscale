<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../../utils/empty.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		avg(
			(
				(
					node_memory_MemTotal{juju_model_uuid=~"${scope.uuid}"}
				or
					node_memory_MemTotal_bytes{juju_model_uuid=~"${scope.uuid}"}
				)
			-
				(
						(
							node_memory_MemFree{juju_model_uuid=~"${scope.uuid}"}
						or
							node_memory_MemFree_bytes{juju_model_uuid=~"${scope.uuid}"}
						)
					+
						(
							node_memory_Cached{juju_model_uuid=~"${scope.uuid}"}
						or
							node_memory_Cached_bytes{juju_model_uuid=~"${scope.uuid}"}
						)
					+
					(
						node_memory_Buffers{juju_model_uuid=~"${scope.uuid}"}
						or
						node_memory_Buffers_bytes{juju_model_uuid=~"${scope.uuid}"}
					)
				+
					(
						node_memory_Slab{juju_model_uuid=~"${scope.uuid}"}
					or
						node_memory_Slab_bytes{juju_model_uuid=~"${scope.uuid}"}
					)
				)
			)
		/
			(
				node_memory_MemTotal{juju_model_uuid=~"${scope.uuid}"}
			or
				node_memory_MemTotal_bytes{juju_model_uuid=~"${scope.uuid}"}
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
		<NoData />
	{:else}
		{@const [result] = results}
		{@const value = result.value.value}
		{@const utilization = value * 100}
		<p class="text-5xl">{utilization.toFixed(2)}%</p>
	{/if}
{:catch error}
	Error
{/await}
