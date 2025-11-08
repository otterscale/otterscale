<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		ConfigurationService,
		type Configuration,
		type Configuration_PackageRepository,
		type UpdatePackageRepositoryRequest
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		packageRepository,
		configuration
	}: {
		packageRepository: Configuration_PackageRepository;
		configuration: Writable<Configuration>;
	} = $props();

	const transport: Transport = getContext('transport');

	const client = createClient(ConfigurationService, transport);
	const defaults = {
		id: packageRepository.id,
		url: packageRepository.url,
		skipJuju: false
	} as UpdatePackageRepositoryRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<span>
	<Modal.Root bind:open>
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
							toast.promise(() => client.updatePackageRepository(request), {
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
</span>
