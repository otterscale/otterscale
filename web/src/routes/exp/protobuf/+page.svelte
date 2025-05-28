<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Nexus, type Application_Chart_Metadata } from '$gen/api/nexus/v1/nexus_pb';
	import { getContext, onMount } from 'svelte';
	import { get, writable } from 'svelte/store';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Unstruct } from '$lib/components/otterscale/ui/index';
	import type { JsonObject } from '@bufbuild/protobuf';
	import Display from './display.svelte';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const Metadata = writable<Application_Chart_Metadata>();
	const isLoading = writable(true);

	async function fetchMetadata() {
		try {
			const response = await client.getApplication({
				scopeUuid: 'db23d197-9178-4202-874c-b9374bc9987e',
				facilityName: 'kubernetes-worker',
				namespace: 'argocd',
				name: 'argocd-applicationset-controller'
			});
			if (response.metadata) {
				Metadata.set(response.metadata);
			}
		} catch (error) {
			console.error('Error fetching machines:', error);
		} finally {
			isLoading.set(false);
		}
	}

	let parsedMetadata = $state({} as JsonObject);
	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchMetadata();
			if ($Metadata?.customization?.values) {
				parsedMetadata = JSON.parse(JSON.stringify($Metadata.customization.values));
			}
			console.log($Metadata);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
		mounted = true;
	});
</script>

{#if mounted}
	<div class="p-4">
		<Unstruct bind:data={parsedMetadata} />
	</div>
{:else}
	<p>Loading...</p>
{/if}
