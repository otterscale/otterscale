<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
	import { PrometheusDriver } from 'prometheus-query';

	// Props
	let {
		client,
		scope,
		isReloading = $bindable()
	}: { client: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	// Constants
	const CHART_TITLE = m.cluster_health();
	const CHART_DESCRIPTION = m.status();

	// Query
	const query = $derived(
		`
		ceph_health_status{juju_model_uuid=~"${scope.uuid}"}
		`
	);

	// Health status mappings
	const HEALTH_STATUS = {
		0: {
			label: 'HEALTHY',
			color: 'text-healthy',
			icon: 'ph:check-bold',
			iconClass: '-right-6 top-4'
		},
		1: {
			label: 'WARNING',
			color: 'text-warning',
			icon: 'ph:exclamation-mark',
			iconClass: '-right-3 top-2'
		},
		2: { label: 'ERROR', color: 'text-error', icon: 'ph:x-bold', iconClass: '-right-3 top-2' }
	} as const;
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	<Card.Root class="relative gap-2 overflow-hidden">
		<Card.Header class="items-center">
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		{@const value = response.result[0].value.value}
		{@const healthStatus = HEALTH_STATUS[value as keyof typeof HEALTH_STATUS]}
		<Card.Content class="flex-1 {healthStatus?.color}">
			{healthStatus?.label}
			<Icon
				icon={healthStatus.icon}
				class="text-primary/5 absolute size-36 text-nowrap text-8xl uppercase tracking-tight group-hover:hidden {healthStatus.iconClass}"
			/>
		</Card.Content>
	</Card.Root>
{:catch error}
	<Card.Root class="gap-2">
		<Card.Header class="items-center">
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		<Card.Content class="flex-1">LOADING ERROR</Card.Content>
	</Card.Root>
{/await}
