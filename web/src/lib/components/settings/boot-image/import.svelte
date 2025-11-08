<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		type Configuration,
		ConfigurationService	} from '$lib/api/configuration/v1/configuration_pb';
	import { Button } from '$lib/components/ui/button';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { configuration }: { configuration: Writable<Configuration> } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ConfigurationService, transport);

	let isImportingBootImages = $state(false);

	async function checkImportingStatus() {
		while (true) {
			const response = await client.isImportingBootImages({});
			if (response.importing) {
				isImportingBootImages = true;
				await new Promise((resolve) => setTimeout(resolve, 5000));
			} else {
				isImportingBootImages = false;
				break;
			}
		}
	}

	onMount(async () => {
		await checkImportingStatus();
	});
</script>

<Button
	variant="ghost"
	disabled={isImportingBootImages}
	onclick={() => {
		toast.promise(() => client.importBootImages({}), {
			loading: 'Loading...',
			success: () => {
				isImportingBootImages = true;
				client.getConfiguration({}).then((response) => {
					configuration.set(response);
				});
				checkImportingStatus();
				return `Import boot images success`;
			},
			error: (error) => {
				let message = `Fail to import boot images`;
				toast.error(message, {
					description: (error as ConnectError).message.toString(),
					duration: Number.POSITIVE_INFINITY
				});
				return message;
			}
		});
	}}
	class="flex items-center gap-2"
>
	{#if isImportingBootImages == true}
		<Icon icon="ph:spinner" class="size-5 animate-spin text-muted-foreground" />
		{m.importing()}
	{:else}
		<Icon icon="ph:arrows-clockwise" />
		{m.import()}
	{/if}
</Button>
