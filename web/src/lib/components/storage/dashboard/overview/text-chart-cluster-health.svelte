<script lang="ts">
	import Icon from '@iconify/svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	// Props
	let {
		client,
		scope,
		isReloading = $bindable()
	}: { client: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	// Constants
	const CHART_TITLE = m.cluster_health();
	const CHART_DESCRIPTION = m.status();

	// Query
	const query = $derived(
		`
		ceph_health_status{juju_model="${scope}"}
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

	onMount(async () => {
		await fetch();
		isLoading = false;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

{#if isLoading}
	<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
		<Card.Header>
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		<Card.Content class="text-3xl">
			<div class="flex h-9 w-full items-center justify-center">
				<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
			</div>
		</Card.Content>
	</Card.Root>
{:else if response === null}
	<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
		<Card.Header>
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		<Card.Content class="text-3xl">
			<div class="flex h-full w-full flex-col items-center justify-center">
				<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
				<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
			</div>
		</Card.Content>
	</Card.Root>
{:else}
	<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
		{@const healthStatus = HEALTH_STATUS[response as keyof typeof HEALTH_STATUS]}
		<Icon
			icon={healthStatus.icon}
			class="absolute size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden {healthStatus.iconClass}"
		/>
		<Card.Header class="items-center">
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		<Card.Content class="text-3xl {healthStatus?.color}">
			{healthStatus?.label}
		</Card.Content>
	</Card.Root>
{/if}
