<script lang="ts">
	import Icon from '@iconify/svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
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
	const CHART_TITLE = m.time_till_full();
	const CHART_DESCRIPTION = m.six_hour_average();

	const query = $derived(
		`
	(
		ceph_pool_max_avail{
			job=~".+",
			juju_application=~".*",
			juju_model=~".*",
			juju_model="${scope}",
			juju_unit=~".*"
		} 
		/ 
		deriv(ceph_pool_stored{
			job=~".+",
			juju_application=~".*",
			juju_model=~".*",
			juju_model="${scope}",
			juju_unit=~".*"
		}[6h])
	) 
	* 
	on(pool_id) group_left(instance, name) ceph_pool_metadata{
		job=~".+",
		juju_application=~".*",
		juju_model=~".*",
		juju_model="${scope}",
		juju_unit=~".*",
		name=~".mgr"
	} > 0
	`
	);

	// Format time duration
	function formatTimeTillFull(days: number): string {
		if (!isFinite(days) || days <= 0) {
			return '∞ years';
		}

		if (days < 1) {
			const hours = Math.round(days * 24);
			return `${hours} hour${hours !== 1 ? 's' : ''}`;
		} else if (days < 30) {
			const roundedDays = Math.round(days);
			return `${roundedDays} day${roundedDays !== 1 ? 's' : ''}`;
		} else if (days < 365) {
			const months = Math.round(days / 30);
			return `${months} month${months !== 1 ? 's' : ''}`;
		} else {
			const years = Math.round(days / 365);
			return `${years} year${years !== 1 ? 's' : ''}`;
		}
	}

	// Auto Update
	let response = $state<string>();
	let isLoading = $state(true);
	const reloadManager = new ReloadManager(fetch);

	// Fetch function
	async function fetch(): Promise<void> {
		try {
			const queryResponse = await client.instantQuery(query);

			if (queryResponse.result && queryResponse.result.length > 0) {
				const days = parseFloat(queryResponse.result[0].value.value);
				response = formatTimeTillFull(days);
			} else {
				response = '∞ years';
			}
		} catch (err) {
			console.error('Failed to fetch cluster health:', err);
			response = 'ERROR';
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

<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
	<Icon
		icon="ph:clock"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header class="items-center">
		<Card.Title>{CHART_TITLE}</Card.Title>
		<Card.Description>{CHART_DESCRIPTION}</Card.Description>
	</Card.Header>
	{#if isLoading}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if response === undefined}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content class="text-3xl">{response}</Card.Content>
	{/if}
</Card.Root>
