<script lang="ts">
	import Icon from '@iconify/svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Button } from '$lib/components/ui/button';
	import { getContext } from 'svelte';
	import { Nexus, type ImportBootImagesRequest } from '$gen/api/nexus/v1/nexus_pb';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { toast } from 'svelte-sonner';
	import { onMount } from 'svelte';

	let {
		isImportingBootImages = $bindable()
	}: {
		isImportingBootImages: boolean;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {} as ImportBootImagesRequest;
	let importBootImageRequest = $state(DEFAULT_REQUEST);

	function reset() {
		importBootImageRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	async function fetchIsImportingBootImages() {
		await new Promise((resolve) => setTimeout(resolve, 3000)); // Wait 5 seconds between checks

		while (true) {
			const response = await client.isImportingBootImages({});
			console.log(`Checking importing status: ${response.importing}`);
			if (!response.importing) {
				isImportingBootImages = false;
				console.log(`Finished`);
				break;
			} else {
				await new Promise((resolve) => setTimeout(resolve, 5000)); // Wait 5 seconds between checks
			}
		}
	}

	onMount(() => {});
</script>

{#if isImportingBootImages == true}
	{#await fetchIsImportingBootImages()}
		<span class="flex items-center gap-2 text-sm text-muted-foreground">
			<Icon icon="ph:spinner" class="size-5 animate-spin" />
			Importing
		</span>
	{/await}
{:else}
	<AlertDialog.Root bind:open>
		<AlertDialog.Trigger>
			<Button variant="ghost" class="flex items-center gap-2">
				<Icon icon="ph:arrows-clockwise" />
				Import
			</Button>
		</AlertDialog.Trigger>
		<AlertDialog.Content>
			<AlertDialog.Header>
				<AlertDialog.Title>Import Boot Images</AlertDialog.Title>
				<AlertDialog.Description class="rounded-lg bg-muted/50 p-4">
					Are you sure you want to import boot images?
				</AlertDialog.Description>
			</AlertDialog.Header>
			<AlertDialog.Footer>
				<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
				<AlertDialog.Action
					onclick={() => {
						isImportingBootImages = true;
						client
							.importBootImages(importBootImageRequest)
							.then((r) => {
								toast.info(`Import boot images success`);
							})
							.catch((e) => {
								toast.error(`Fail to import boot images: ${e.toString()}`);
							});
						reset();
						close();
					}}
				>
					Import
				</AlertDialog.Action>
			</AlertDialog.Footer>
		</AlertDialog.Content>
	</AlertDialog.Root>
{/if}
