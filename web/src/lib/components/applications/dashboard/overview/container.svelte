<script lang="ts">
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
	import Icon from '@iconify/svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { onMount } from 'svelte';

	let {
		prometheusDriver,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; isReloading: boolean } = $props();

	let runningContainers = $state(0);

	async function fetch() {
		prometheusDriver
			.instantQuery(
				`
						sum(
							kubelet_running_containers{job="kubelet",juju_model_uuid="${$currentKubernetes?.scopeUuid}",metrics_path="/metrics"}
						)
						or
						sum(
							kubelet_running_container_count{job="kubelet",juju_model_uuid="${$currentKubernetes?.scopeUuid}",metrics_path="/metrics"}
						)
						`
			)
			.then((response) => {
				runningContainers = response.result[0].value.value;
			});
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		await fetch();
		isLoading = false;
	});

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});
</script>

{#if isLoading}
	Loading
{:else}
	<Card.Root class="relative h-full gap-2 overflow-hidden">
		<Icon
			icon="ph:shipping-container"
			class="text-primary/5 absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
		/>
		<Card.Header>
			<Card.Title>{m.containers()}</Card.Title>
			<Card.Description>{m.running()}</Card.Description>
		</Card.Header>
		<Card.Content class="text-3xl">
			{runningContainers}
		</Card.Content>
	</Card.Root>
{/if}
