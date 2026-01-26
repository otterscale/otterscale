<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		type Configuration,
		type Configuration_PackageRepository,
		ConfigurationService,
		type UpdatePackageRepositoryRequest
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';

	let {
		packageRepository,
		configuration = $bindable()
	}: {
		packageRepository: Configuration_PackageRepository;
		configuration: Configuration;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ConfigurationService, transport);

	let open = $state(false);

	let request = $state({} as UpdatePackageRepositoryRequest);
	function init() {
		request = {
			id: packageRepository.id,
			url: packageRepository.url
		} as UpdatePackageRepositoryRequest;
	}

	function close() {
		open = false;
	}
</script>

<span>
	<Modal.Root
		bind:open
		onOpenChange={(isOpen) => {
			if (isOpen) init();
		}}
	>
		<Modal.Trigger variant="creative">
			<Icon icon="ph:pencil" />
			{m.edit()}
		</Modal.Trigger>
		<Modal.Content>
			<Modal.Header>{m.edit_package_repository()}</Modal.Header>
			<Form.Root>
				<Form.Fieldset>
					<Form.Field>
						<Form.Label>{m.url()}</Form.Label>
						<SingleInput.General type="text" bind:value={request.url} />
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
							toast.promise(() => client.updatePackageRepository(request), {
								loading: 'Loading...',
								success: () => {
									client.getConfiguration({}).then((response) => {
										configuration = response;
									});
									return `Update ${packageRepository.name} success`;
								},
								error: (error) => {
									let message = `Fail to update ${packageRepository.name}`;
									toast.error(message, {
										description: (error as ConnectError).message.toString(),
										duration: Number.POSITIVE_INFINITY,
										closeButton: true
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
</span>
