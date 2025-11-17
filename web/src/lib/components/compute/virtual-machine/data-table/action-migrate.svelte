<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { MigrateInstanceRequest, VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		virtualMachine,
		scope,
		reloadManager
	}: { virtualMachine: VirtualMachine; scope: string; reloadManager: ReloadManager } = $props();

	const transport: Transport = getContext('transport');

	const virtualMachineClient = createClient(InstanceService, transport);
	let invalid = $state(false);

	const defaults = {
		scope: scope,
		namespace: virtualMachine.namespace,
		name: virtualMachine.name,
		hostname: ''
	} as MigrateInstanceRequest;
	let request = $state({ ...defaults });
	function reset() {
		request = { ...defaults };
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:arrows-clockwise" />
		{m.migrate()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.migrate_virtual_machine()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.hostname()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.hostname} bind:invalid />
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
					disabled={invalid}
					onclick={() => {
						toast.promise(() => virtualMachineClient.migrateInstance(request), {
							loading: `Migrating ${request.name} to ${request.hostname}...`,
							success: () => {
								reloadManager.force();
								return `Successfully migrated ${request.name} to ${request.hostname}`;
							},
							error: (error) => {
								let message = `Failed to migrate ${request.name}`;
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
