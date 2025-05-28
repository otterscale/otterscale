<script lang="ts" module>
	const healthMap: Record<string, string> = {
		0: 'HEALTH',
		1: 'WARN',
		2: 'ERROR'
	};
	const healthColorMap: Record<string, string> = {
		0: 'text-green-900 dark:text-green-800',
		1: 'text-yellow-900 dark:text-yellow-800',
		2: 'text-red-900 dark:text-red-800'
	};
</script>

<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../../utils/empty.svelte';
	import { cn } from '$lib/utils';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		ceph_health_status{juju_model_uuid=~"${scope.uuid}"}
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
		{@const health = String(result.value.value)}
		<p class={cn('text-5xl', healthColorMap[health])}>{healthMap[health]}</p>
	{/if}
{:catch error}
	Error
{/await}
