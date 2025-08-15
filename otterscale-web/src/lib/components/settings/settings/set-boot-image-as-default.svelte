<script lang="ts" module>
	import {
		ConfigurationService,
		type Configuration,
		type Configuration_BootImage,
		type SetDefaultBootImageRequest
	} from '$lib/api/configuration/v1/configuration_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		bootImage,
		configuration = $bindable()
	}: {
		bootImage: Configuration_BootImage;
		configuration: Writable<Configuration>;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ConfigurationService, transport);

	const DEFAULT_REQUEST = {
		distroSeries: bootImage.distroSeries
	} as SetDefaultBootImageRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	let name = $state('');

	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:star" />
		Default
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Set Default Boot Image</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.Confirm target={bootImage.name} bind:value={name} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel onclick={reset}>Cancel</Modal.Cancel>
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

						reset();
						stateController.close();
					}}
				>
					Create
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
