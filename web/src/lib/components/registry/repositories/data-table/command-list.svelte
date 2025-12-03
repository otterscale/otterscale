<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';

	import { RegistryService } from '$lib/api/registry/v1/registry_pb';
	import CopyButton from '$lib/components/custom/copy-button/copy-button.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
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

	let sourceImage = $state('');
	const sourceImageToken = $derived(sourceImage || 'SOURCE_IMAGE');
	let repository = $state('');
	const repositoryToken = $derived(repository || 'REPOSITORY');
	let tag = $state('');
	const tagToken = $derived(tag.trim() ? `:${tag}` : tag || '[:TAG]');
	function reset() {
		sourceImage = '';
		repository = '';
		tag = '';
	}

	let isLoaded = $state(false);
	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
</script>

{#if isLoaded}
	<Dialog.Root
		onOpenChange={(isOpen) => {
			if (isOpen) {
				reset();
			}
		}}
	>
		<Dialog.Trigger class={buttonVariants({ variant: 'ghost' })}>
			<span class="flex items-center gap-2">
				<Icon icon="ph:code" />
				{m.commands()}
			</span>
		</Dialog.Trigger>
		<Dialog.Content class="min-w-[50vw] overflow-y-auto">
			<Dialog.Header>
				<Dialog.Title class="text-center">{m.commands()}</Dialog.Title>
			</Dialog.Header>

			<div class="flex items-center">
				<div class="flex w-full flex-col gap-0">
					<Item.Root class="w-full">
						{@const command = `docker tag ${sourceImageToken}${tagToken} ${registryURL}/${repositoryToken}${tagToken}`}
						<Item.Media variant="icon">
							<Icon icon="ph:tag" />
						</Item.Media>
						<Item.Content class="flex flex-col items-start">
							<Item.Description>Tag an image for a project.</Item.Description>
							<Item.Title><p class="font-mono text-xs">{command}</p></Item.Title>
						</Item.Content>
						<Item.Actions>
							<CopyButton text={command} />
						</Item.Actions>
					</Item.Root>

					<Item.Root class="w-full">
						{@const command = `docker push ${registryURL}/${repositoryToken}${tagToken}`}
						<Item.Media variant="icon">
							<Icon icon="ph:cloud" />
						</Item.Media>
						<Item.Content class="flex flex-col items-start">
							<Item.Description>Push an image to a project.</Item.Description>
							<Item.Title><p class="font-mono text-xs">{command}</p></Item.Title>
						</Item.Content>
						<Item.Actions>
							<CopyButton text={command} />
						</Item.Actions>
					</Item.Root>
				</div>
				<div class="flex h-full max-w-52 flex-col justify-between p-4">
					<InputGroup.Root class="h-8">
						<InputGroup.Addon>
							<Icon icon="ph:package" />
						</InputGroup.Addon>
						<InputGroup.Addon>
							<InputGroup.Text class="text-xs">{m.image()}</InputGroup.Text>
						</InputGroup.Addon>
						<InputGroup.Input bind:value={sourceImage} />
					</InputGroup.Root>
					<InputGroup.Root class="h-8">
						<InputGroup.Addon>
							<Icon icon="ph:package" />
						</InputGroup.Addon>
						<InputGroup.Addon>
							<InputGroup.Text class="text-xs">{m.repository()}</InputGroup.Text>
						</InputGroup.Addon>
						<InputGroup.Input bind:value={repository} />
					</InputGroup.Root>
					<InputGroup.Root class="h-8">
						<InputGroup.Addon>
							<Icon icon="ph:tag" />
						</InputGroup.Addon>
						<InputGroup.Addon>
							<InputGroup.Text class="text-xs">{m.tag()}</InputGroup.Text>
						</InputGroup.Addon>
						<InputGroup.Input bind:value={tag} />
					</InputGroup.Root>
				</div>
			</div>
		</Dialog.Content>
	</Dialog.Root>
{/if}
