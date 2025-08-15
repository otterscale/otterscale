<script lang="ts" module>
	import {
		ConfigurationService,
		type Configuration,
		type Configuration_PackageRepository,
		type UpdatePackageRepositoryRequest
	} from '$lib/api/configuration/v1/configuration_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
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
		packageRepository,
		configuration = $bindable()
	}: {
		packageRepository: Configuration_PackageRepository;
		configuration: Writable<Configuration>;
	} = $props();

	const transport: Transport = getContext('transport');

	const client = createClient(ConfigurationService, transport);
	const requestManager = new RequestManager<UpdatePackageRepositoryRequest>({
		id: packageRepository.id,
		url: packageRepository.url,
		skipJuju: false
	} as UpdatePackageRepositoryRequest);
	const stateController = new StateController(false);
</script>

<div>
	<Modal.Root bind:open={stateController.state}>
		<Modal.Trigger variant="creative">
			<Icon icon="ph:pencil" />
			Edit
		</Modal.Trigger>
		<Modal.Content>
			<Modal.Header>Edit NTP Server</Modal.Header>
			<Form.Root>
				<Form.Fieldset>
					<Form.Field>
						<Form.Label>URL</Form.Label>
						<SingleInput.General type="text" bind:value={requestManager.request.url} />
					</Form.Field>
				</Form.Fieldset>
			</Form.Root>
			<Modal.Footer>
				<Modal.Cancel
					onclick={() => {
						requestManager.reset();
					}}>Cancel</Modal.Cancel
				>
				<Modal.ActionsGroup>
					<Modal.Action
						onclick={() => {
							toast.promise(() => client.updatePackageRepository(requestManager.request), {
								loading: 'Loading...',
								success: () => {
									client.getConfiguration({}).then((response) => {
										configuration.set(response);
									});
									return `Update ${packageRepository.name} success`;
								},
								error: (error) => {
									let message = `Fail to update ${packageRepository.name}`;
									toast.error(message, {
										description: (error as ConnectError).message.toString(),
										duration: Number.POSITIVE_INFINITY
									});
									return message;
								}
							});

							requestManager.reset();
							stateController.close();
						}}
					>
						Edit
					</Modal.Action>
				</Modal.ActionsGroup>
			</Modal.Footer>
		</Modal.Content>
	</Modal.Root>
</div>
