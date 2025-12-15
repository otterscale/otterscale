<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';

	import { RegistryService } from '$lib/api/registry/v1/registry_pb';
	import * as Code from '$lib/components/custom/code';
	import CopyButton from '$lib/components/custom/copy-button/copy-button.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Item from '$lib/components/ui/item';
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
		<Dialog.Content class="max-h-[77vh] w-full min-w-[50vw] overflow-y-auto">
			<Dialog.Header>
				<Dialog.Title class="text-center">{m.commands()}</Dialog.Title>
			</Dialog.Header>

			<Item.Root class="w-full">
				{@const command = `docker push ${registryURL}/<repository>[:<tag>]`}
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
				{@const command = `helm push <chart_package> oci://${registryURL}/<namespace> --plain-http`}
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

			<div class="m-4 space-y-2 rounded-lg bg-muted p-4">
				<h3 class="text-sm font-semibold">{m.trouble_pushing_image()}</h3>
				<Collapsible.Root>
					<Collapsible.Trigger class="text-sm hover:underline">
						{@html m.trouble_pushing_image_step_one({
							daemon: '<span class="font-mono">daemon.json</span>',
							insecure_registries: '<span class="font-mono">insecure-registries</span>'
						})}
					</Collapsible.Trigger>
					<Collapsible.Content>
						<div class="space-y-2 p-4 text-sm text-muted-foreground">
							{@html m.trouble_pushing_image_step_one_details({
								insecure_registries: '<span class="font-mono">insecure-registries</span>',
								daemon: '<span class="font-mono">daemon.json</span>',
								path: '<span class="font-mono">/etc/docker/daemon.json</span>'
							})}
							<Code.Root
								class="w-fit border-none bg-transparent"
								lang="json"
								code={JSON.stringify(
									{
										'insecure-registries': [registryURL]
									},
									null,
									2
								)}
								hideLines
							/>
						</div>
					</Collapsible.Content>
				</Collapsible.Root>

				<Collapsible.Root>
					<Collapsible.Trigger class="text-sm hover:underline">
						{m.trouble_pushing_image_step_two()}
					</Collapsible.Trigger>
					<Collapsible.Content>
						<div class="space-y-2 p-4 text-sm text-muted-foreground">
							{m.trouble_pushing_image_step_two_details()}
							<Code.Root
								class="w-fit border-none bg-transparent"
								lang="bash"
								code={['sudo systemctl daemon-reload', 'sudo systemctl restart docker'].join('\n')}
								hideLines
							/>
						</div>
					</Collapsible.Content>
				</Collapsible.Root>
			</div>
		</Dialog.Content>
	</Dialog.Root>
{/if}
