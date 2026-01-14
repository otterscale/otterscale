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

	type PVCData = {
		namespace: string;
		pvc: string;
		storageClass: string;
		capacity: number;
		used: number;
		available: number;
		status: string;
		usedPercent: number;
	};

	let pvcData: PVCData[] = $state([]);

	async function fetchPVCData() {
		try {
			// Fetch PVC info
			const infoResponse = await prometheusDriver.instantQuery(
				`kube_persistentvolumeclaim_info{juju_model="${scope}", juju_unit=""}`
			);

			// Fetch capacity
			const capacityResponse = await prometheusDriver.instantQuery(
				`kubelet_volume_stats_capacity_bytes{juju_model="${scope}", juju_unit=""}`
			);

			// Fetch used
			const usedResponse = await prometheusDriver.instantQuery(
				`kubelet_volume_stats_used_bytes{juju_model="${scope}", juju_unit=""}`
			);

			// Fetch available
			const availableResponse = await prometheusDriver.instantQuery(
				`kubelet_volume_stats_available_bytes{juju_model="${scope}", juju_unit=""}`
			);

			// Fetch status
			const statusResponse = await prometheusDriver.instantQuery(
				`kube_persistentvolumeclaim_status_phase{juju_model="${scope}", juju_unit=""}`
			);

			// Process data
			const pvcMap = new Map<string, PVCData>();

			// Initialize from info
			infoResponse.result.forEach((item: any) => {
				const labels = item.metric.labels;
				const key = `${labels.namespace}-${labels.persistentvolumeclaim}`;
				pvcMap.set(key, {
					namespace: labels.namespace,
					pvc: labels.persistentvolumeclaim,
					storageClass: labels.storageclass || '',
					capacity: 0,
					used: 0,
					available: 0,
					status: 'Bound',
					usedPercent: 0
				});
			});

			// Add capacity
			capacityResponse.result.forEach((item: any) => {
				const labels = item.metric.labels;
				const key = `${labels.namespace}-${labels.persistentvolumeclaim}`;
				if (pvcMap.has(key)) {
					pvcMap.get(key)!.capacity = item.value.value;
				}
			});

			// Add used
			usedResponse.result.forEach((item: any) => {
				const labels = item.metric.labels;
				const key = `${labels.namespace}-${labels.persistentvolumeclaim}`;
				if (pvcMap.has(key)) {
					pvcMap.get(key)!.used = item.value.value;
				}
			});

			// Add available
			availableResponse.result.forEach((item: any) => {
				const labels = item.metric.labels;
				const key = `${labels.namespace}-${labels.persistentvolumeclaim}`;
				if (pvcMap.has(key)) {
					pvcMap.get(key)!.available = item.value.value;
				}
			});

			// Add status
			statusResponse.result.forEach((item: any) => {
				const labels = item.metric.labels;
				const key = `${labels.namespace}-${labels.persistentvolumeclaim}`;
				if (pvcMap.has(key) && Number(item.value[1]) === 1) {
					pvcMap.get(key)!.status = labels.phase;
				}
			});

			// Calculate used percent
			pvcMap.forEach((pvc) => {
				if (pvc.capacity > 0) {
					pvc.usedPercent = (pvc.used / pvc.capacity) * 100;
				}
			});

			pvcData = Array.from(pvcMap.values())
				.sort((a, b) => b.usedPercent - a.usedPercent)
				.slice(0, 10);
		} catch (error) {
			console.error('Failed to fetch PVC data:', error);
		}
	}

	async function fetch() {
		await fetchPVCData();
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

	function formatBytes(bytes: number): string {
		const gb = bytes / (1024 ** 3);
		return `${gb.toFixed(1)} GB`;
	}

	function getStatusColor(status: string): string {
		switch (status) {
			case 'Bound':
				return 'text-green-600';
			case 'Pending':
				return 'text-yellow-600';
			case 'Lost':
				return 'text-red-600';
			default:
				return 'text-gray-600';
		}
	}
</script>

<Card.Root class="relative overflow-hidden">
	<Icon
		icon="ph:hard-drive"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header>
		<Card.Title>Persistent Volume Claims</Card.Title>
		<Card.Description>Storage usage across namespaces</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-32 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if pvcData.length === 0}
		<div class="flex h-32 w-full flex-col items-center justify-center">
			<Icon icon="ph:hard-drive-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content class="p-0">
			<div class="overflow-x-auto">
				<table class="w-full text-sm">
					<thead class="bg-muted/50">
						<tr>
							<th class="px-4 py-2 text-left font-medium">Namespace</th>
							<th class="px-4 py-2 text-left font-medium">PVC</th>
							<th class="px-4 py-2 text-left font-medium">StorageClass</th>
							<th class="px-4 py-2 text-right font-medium">Capacity</th>
							<th class="px-4 py-2 text-right font-medium">Used</th>
							<th class="px-4 py-2 text-right font-medium">Available</th>
							<th class="px-4 py-2 text-center font-medium">Status</th>
							<th class="px-4 py-2 text-right font-medium">Used %</th>
						</tr>
					</thead>
					<tbody>
						{#each pvcData as pvc}
							<tr class="border-t">
								<td class="px-4 py-2">{pvc.namespace}</td>
								<td class="px-4 py-2 font-mono text-xs">{pvc.pvc}</td>
								<td class="px-4 py-2">{pvc.storageClass}</td>
								<td class="px-4 py-2 text-right">{formatBytes(pvc.capacity)}</td>
								<td class="px-4 py-2 text-right">{formatBytes(pvc.used)}</td>
								<td class="px-4 py-2 text-right">{formatBytes(pvc.available)}</td>
								<td class="px-4 py-2 text-center">
									<span class="inline-flex items-center rounded-full px-2 py-1 text-xs font-medium {getStatusColor(pvc.status)} bg-current/10">
										{pvc.status}
									</span>
								</td>
								<td class="px-4 py-2 text-right">{pvc.usedPercent.toFixed(1)}%</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</Card.Content>
	{/if}
</Card.Root>