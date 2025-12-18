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
	import * as Tooltip from '$lib/components/ui/tooltip';
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

	let request = $state({} as MigrateInstanceRequest);
	let invalid = $state(false);
	let open = $state(false);

	function init() {
		request = {
			scope: scope,
			namespace: virtualMachine.namespace,
			name: virtualMachine.name,
			hostname: ''
		} as MigrateInstanceRequest;
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
>
	<!-- TODO: disabled until feature is implemented -->
	<!-- <Modal.Trigger variant="creative">
		<Icon icon="ph:arrows-clockwise" />
		{m.migrate()}
	</Modal.Trigger> -->
	<Tooltip.Provider>
		<Tooltip.Root>
			<Tooltip.Trigger class="w-full">
				<Modal.Trigger variant="creative" disabled>
					<Icon icon="ph:arrows-clockwise" />
					{m.migrate()}
				</Modal.Trigger>
			</Tooltip.Trigger>
			<Tooltip.Content>{m.under_development()}</Tooltip.Content>
		</Tooltip.Root>
	</Tooltip.Provider>
	<Modal.Content>
		<Modal.Header>{m.migrate_virtual_machine()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.hostname()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.hostname} required bind:invalid />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
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
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
