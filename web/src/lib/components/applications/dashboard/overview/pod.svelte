<script lang="ts">
	import Icon from '@iconify/svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();
	let runningPods = $state(0);

	async function fetch() {
		prometheusDriver
			.instantQuery(
				`
						sum(
							kubelet_running_pods{job="kubelet",juju_model_uuid="${scope.uuid}",metrics_path="/metrics"}
						)
						or
						sum(
							kubelet_running_pod_count{job="kubelet",juju_model_uuid="${scope.uuid}",metrics_path="/metrics"}
						)
						`
			)
			.then((response) => {
				runningPods = response.result[0].value.value;
			});
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		await fetch();
		isLoading = false;
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
</script>

{#if isLoading}
	Loading
{:else}
	<Card.Root class="relative h-full gap-2 overflow-hidden">
		<Icon
			icon="ph:squares-four"
			class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
		/>
		<Card.Header>
			<Card.Title>{m.pods()}</Card.Title>
			<Card.Description>{m.running()}</Card.Description>
		</Card.Header>
		<Card.Content class="text-3xl">
			{runningPods}
		</Card.Content>
	</Card.Root>
{/if}
