<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let version: string | undefined = $state(undefined);
	let platform: string | undefined = $state(undefined);
	async function fetchBuildInfo() {
		const response = await prometheusDriver.instantQuery(
			`kubernetes_build_info{juju_model="${scope}", job="apiserver", service="kubernetes"}`
		);
		version = response.result[0]?.metric?.labels?.git_version ?? undefined;
		platform = response.result[0]?.metric?.labels?.platform ?? undefined;
	}

	async function fetch() {
		try {
			await fetchBuildInfo();
		} catch (error) {
			console.error('Failed to fetch build info:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	let isLoaded = $state(false);
	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
	<Icon
		icon="ph:cube"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header>
		<Card.Title>Information</Card.Title>
		<Card.Description class="flex h-6 items-center">Cluster Version and Platform</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if !version}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content class="flex items-baseline gap-2">
			<span class="text-3xl">{version ?? 'N/A'}</span>
			<span class="text-1xl font-medium tracking-wider text-muted-foreground uppercase"
				>{platform ?? 'N/A'}</span
			>
		</Card.Content>
	{/if}
</Card.Root>
