<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type {
		CreateVirtualMachineRestoreRequest,
		VirtualMachine
	} from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Form from '$lib/components/custom/form';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	// Props
	let {
		virtualMachine,
		scope,
		reloadManager
	}: { virtualMachine: VirtualMachine; scope: string; reloadManager: ReloadManager } = $props();

	// Context dependencies
	const transport: Transport = getContext('transport');
	const virtualMachineClient = createClient(InstanceService, transport);

	// ==================== State Variables ====================
	// UI state
	let open = $state(false);

	// Form validation state
	let invalid: boolean | undefined = $state();

	// ==================== Form State ====================
	let request = $state({} as CreateVirtualMachineRestoreRequest);

	// ==================== Local Dropdown Options ====================
	const snapshotOptions: Writable<SingleSelect.OptionType[]> = writable([]);

	$effect(() => {
		const options = (virtualMachine.snapshots || [])
			.filter((snapshot) => snapshot.phase == 'Succeeded')
			.map((snapshot) => ({
				value: snapshot.name,
				label: snapshot.name,
				icon: 'ph:camera'
			}));
		snapshotOptions.set(options);
	});

	$effect(() => {
		if (request.snapshotName) {
			const now = new Date();
			const year = now.getFullYear();
			const month = (now.getMonth() + 1).toString().padStart(2, '0');
			const day = now.getDate().toString().padStart(2, '0');
			const hours = now.getHours().toString().padStart(2, '0');
			const minutes = now.getMinutes().toString().padStart(2, '0');
			request.name = `${request.snapshotName}-${year}${month}${day}${hours}${minutes}`;
		}
	});

	// ==================== Utility Functions ====================
	function reset() {
		console.log('reset');
		request = {
			scope: scope,
			namespace: virtualMachine.namespace,
			name: '',
			virtualMachineName: virtualMachine.name,
			snapshotName: ''
		} as CreateVirtualMachineRestoreRequest;
	}
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open onOpenChange={(isOpen) => isOpen && reset()}>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_restore()}</Modal.Header>
		<Form.Root>
			<!-- ==================== Basic Configuration ==================== -->
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleSelect.Root
						required
						options={snapshotOptions}
						bind:value={request.snapshotName}
						bind:invalid
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $snapshotOptions as option (option.value)}
											<SingleSelect.Item {option}>
												<Icon
													icon={option.icon ? option.icon : 'ph:empty'}
													class={cn('size-5', option.icon ? 'visible' : 'invisible')}
												/>
												{option.label}
												<SingleSelect.Check {option} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
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
						toast.promise(() => virtualMachineClient.createVirtualMachineRestore(request), {
							loading: `Creating restore ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully created restore ${request.name}`;
							},
							error: (error) => {
								let message = `Failed to create restore ${request.name}`;
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
