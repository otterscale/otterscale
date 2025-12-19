<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		type Configuration,
		type Configuration_BootImage,
		ConfigurationService,
		type SetDefaultBootImageRequest
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		bootImage,
		configuration,
		closeActions
	}: {
		bootImage: Configuration_BootImage;
		configuration: Writable<Configuration>;
		closeActions: () => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ConfigurationService, transport);

	let request = $state({} as SetDefaultBootImageRequest);
	let open = $state(false);

	function init() {
		request = {
			distroSeries: bootImage.distroSeries
		} as SetDefaultBootImageRequest;
	}

	function close() {
		open = false;
	}
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
	onOpenChangeComplete={(isOpen) => {
		if (!isOpen) {
			closeActions();
		}
	}}
>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:star" />
		{m.set_as_default()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.set_default_boot_image()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.Confirm target={bootImage.name} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					onclick={() => {
						toast.promise(() => client.setDefaultBootImage(request), {
							loading: 'Loading...',
							success: () => {
								client.getConfiguration({}).then((response) => {
									configuration.set(response);
								});
								return `Set ${request.distroSeries} as default`;
							},
							error: (error) => {
								let message = `Fail to set ${request.distroSeries} as default`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
