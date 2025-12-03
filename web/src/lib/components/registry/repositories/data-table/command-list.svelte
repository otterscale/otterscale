<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';

	import { RegistryService } from '$lib/api/registry/v1/registry_pb';
	import CopyButton from '$lib/components/custom/copy-button/copy-button.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Item from '$lib/components/ui/item/index.js';
	import { m } from '$lib/paraglide/messages';

	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const registryClient = createClient(RegistryService, transport);

	let registryURL = $state('');
	async function fetch() {
		try {
			const response = await registryClient.getRegistryURL({
				scope
			});
			registryURL = response.registryUrl;
		} catch (error) {
			console.error('Failed to fetch registry URL:', error);
		}
	}

	let isLoaded = $state(false);
	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
</script>

{#if isLoaded}
	<Dialog.Root>
		<Dialog.Trigger class={buttonVariants({ variant: 'ghost' })}>
			<span class="flex items-center gap-2">
				<Icon icon="ph:code" />
				{m.commands()}
			</span>
		</Dialog.Trigger>
		<Dialog.Content class="min-w-[38vw] overflow-y-auto">
			<Dialog.Header>
				<Dialog.Title class="text-center">{m.commands()}</Dialog.Title>
			</Dialog.Header>

			<div class="flex flex-col gap-2">
				<Item.Root class="w-full">
					{@const command = `docker push ${registryURL}/$REPOSITORY[:TAG]`}
					<Item.Media variant="icon">
						<Icon icon="logos:docker-icon" />
					</Item.Media>
					<Item.Content class="flex flex-col items-start">
						<Item.Description>{m.push_image_description()}</Item.Description>
						<Item.Title><p class="font-mono text-xs">{command}</p></Item.Title>
					</Item.Content>
					<Item.Actions>
						<CopyButton text={command} />
					</Item.Actions>
				</Item.Root>

				<Item.Root class="w-full">
					{@const command = `helm push CHART_PACKAGE oci://${registryURL}`}
					<Item.Media variant="icon">
						<Icon icon="logos:helm" />
					</Item.Media>
					<Item.Content class="flex flex-col items-start">
						<Item.Description>{m.push_chart_description()}</Item.Description>
						<Item.Title><p class="font-mono text-xs">{command}</p></Item.Title>
					</Item.Content>
					<Item.Actions>
						<CopyButton text={command} />
					</Item.Actions>
				</Item.Root>
			</div>
		</Dialog.Content>
	</Dialog.Root>
{/if}
