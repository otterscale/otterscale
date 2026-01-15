<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';

	import { page } from '$app/state';
	import {
		type APIResource,
		type DiscoveryRequest,
		ResourceService
	} from '$lib/api/resource/v1/resource_pb';

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);

	let apiResources = $state<APIResource[]>([]);
	async function fetchAPIResources() {
		try {
			const response = await resourceClient.discovery({
				cluster: page.params.cluster
			} as DiscoveryRequest);
			apiResources = response.apiResources;
		} catch (error) {
			console.error('Failed to fetch discoveries:', error);
		}
	}
	const resources = $derived(
		apiResources.filter(
			(resource) =>
				resource.group === page.url.searchParams.get('group') &&
				resource.version === page.url.searchParams.get('version') &&
				resource.kind === page.params.kind
		)
	);

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchAPIResources();
			isMounted = true;
		} catch (error) {
			console.error('Error in fetcing api resources:', error);
		}
	});
</script>

{#if isMounted}
	<pre>{JSON.stringify(resources, null, 2)}</pre>
{/if}
