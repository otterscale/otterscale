<script lang="ts" module>
	import {
		ConfigurationService,
		type Configuration,
		type ImportBootImagesRequest,
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Form from '$lib/components/custom/form';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	let { configuration }: { configuration: Writable<Configuration> } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ConfigurationService, transport);

	const DEFAULT_REQUEST = {} as ImportBootImagesRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	let isImportingBootImages = $state(false);
	onMount(async () => {
		while (true) {
			const response = await client.isImportingBootImages({});
			if (response.importing) {
				await new Promise((resolve) => setTimeout(resolve, 5000));
			} else {
				isImportingBootImages = false;
				break;
			}
		}
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger disabled={isImportingBootImages}>
		{#if isImportingBootImages == true}
			<Icon icon="ph:spinner" class="text-muted-foreground size-5 animate-spin" />
			{m.importing()}
		{:else}
			<Icon icon="ph:arrows-clockwise" />
			{m.import()}
		{/if}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_boot_image()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.distro_series()}</Form.Label>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					onclick={() => {
						toast.promise(() => client.importBootImages(request), {
							loading: 'Loading...',
							success: () => {
								client.getConfiguration({}).then((response) => {
									configuration.set(response);
								});
								return `Import boot images success`;
							},
							error: (error) => {
								let message = `Fail to import boot images`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY,
								});
								return message;
							},
						});

						reset();
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
