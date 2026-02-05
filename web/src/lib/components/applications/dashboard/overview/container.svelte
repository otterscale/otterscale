<script lang="ts">
	import Icon from '@iconify/svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
	import { activeNamespace } from '$lib/stores';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let runningContainers = $state(0);

	async function fetch() {
		try {
			const response = await prometheusDriver.instantQuery(
				`
				count(kube_pod_container_status_running{juju_model="${scope}",namespace="${$activeNamespace}"} == 1)
				`
			);
			runningContainers = response.result[0]?.value?.value ?? 0;
		} catch (error) {
			console.error('Failed to fetch running containers:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoaded = $state(false);
	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	$effect(() => {
		if ($activeNamespace !== undefined) {
			fetch();
		}
	});
</script>

<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
	<Icon
		icon="ph:shipping-container"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header>
		<Card.Title>{m.containers()}</Card.Title>
		<Card.Description>{m.running()}</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else}
		<Card.Content class="text-3xl">
			{runningContainers}
		</Card.Content>
	{/if}
</Card.Root>
