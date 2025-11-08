<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { type Application,ApplicationService } from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';

	import { Data } from './data';
</script>

<script lang="ts">
	let {
		scope,
		facility,
		namespace,
		applicationName
	}: { scope: string; facility: string; namespace: string; applicationName: string } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);

	const application = writable<Application>();

	let isMounted = $state(false);
	onMount(async () => {
		try {
			client
				.getApplication({
					scope: scope,
					facility: facility,
					namespace: namespace,
					name: applicationName
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
