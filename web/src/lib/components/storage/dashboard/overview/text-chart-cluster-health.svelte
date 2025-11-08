<script lang="ts">
	import Icon from '@iconify/svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

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
		2: { label: 'ERROR', color: 'text-error', icon: 'ph:x-bold', iconClass: '-right-3 top-2' },
		null: {
			label: 'ERROR',
			color: 'text-muted-foreground',
			icon: 'ph:question-bold',
			iconClass: '-right-3 top-2'
		}
	} as const;

	// Auto Update
	let response = $state<number | null>(null);
	let isLoading = $state(true);
	const reloadManager = new ReloadManager(fetch);

	// Fetch function
	async function fetch(): Promise<void> {
		try {
			const queryResponse = await client.instantQuery(query);

			if (queryResponse.result && queryResponse.result.length > 0) {
				response = Number(queryResponse.result[0].value.value);
			} else {
				response = null;
			}
		} catch (err) {
			console.error('Failed to fetch cluster health:', err);
			response = null;
		}
	}

	// Effects
	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	onMount(() => {
		fetch();
		isLoading = false;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

{#if isLoading}
	<ComponentLoading />
{:else}
	<Card.Root class="relative h-full gap-2 overflow-hidden">
		<Card.Header class="items-center">
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		{@const healthStatus = HEALTH_STATUS[response as keyof typeof HEALTH_STATUS]}
		<Card.Content class="flex-1 {healthStatus?.color}">
			{healthStatus?.label}
			<Icon
				icon={healthStatus.icon}
				class="absolute size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden {healthStatus.iconClass}"
			/>
		</Card.Content>
	</Card.Root>
{/if}
