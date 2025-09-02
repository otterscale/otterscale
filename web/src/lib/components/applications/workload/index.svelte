<script lang="ts" module>
	import { ApplicationService, type Application } from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { Data } from './data';
</script>

<script lang="ts">
	let {
		scopeUuid,
		facilityName,
		namespace,
		applicationName,
	}: { scopeUuid: string; facilityName: string; namespace: string; applicationName: string } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);

	const application = writable<Application>();

	let isMounted = $state(false);
	onMount(async () => {
		try {
			client
				.getApplication({
					scopeUuid: scopeUuid,
					facilityName: facilityName,
					namespace: namespace,
					name: applicationName,
				})
				.then((response) => {
					application.set(response);
					isMounted = true;
				});
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<main>
	{#if isMounted}
		<Data {application} />
	{:else}
		<Loading.Data />
	{/if}
</main>
