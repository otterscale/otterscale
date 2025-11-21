<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { type Application, ApplicationService } from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { Data } from './data';
</script>

<script lang="ts">
	let {
		scope,
		namespace,
		applicationName
	}: { scope: string; namespace: string; applicationName: string } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);

	const application = writable<Application>();
	async function fetch() {
		try {
			const response = await client.getApplication({
				scope: scope,
				namespace: namespace,
				name: applicationName
			});
			application.set(response);
		} catch (error) {
			console.error('Failed to fetch application:', error);
		}
	}
	const reloadManager = new ReloadManager(fetch, false);

	let isMounted = $state(false);
	onMount(async () => {
		await fetch();
		isMounted = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<main>
	{#if isMounted}
		<Data {application} {scope} {namespace} {reloadManager} />
	{:else}
		<Loading.Data />
	{/if}
</main>
