<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { UpdateVirtualMachineRequest, VirtualMachine } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { KubeVirtService } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const KubeVirtClient = createClient(KubeVirtService, transport);

	// Label management state
	let labelKey = $state('');
	let labelValue = $state('');

	// ==================== Default Values & Constants ====================
	const defaults = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		name: virtualMachine.metadata?.name,
		namespace: virtualMachine.metadata?.namespace,
		networkName: virtualMachine.networkName,
		labels: virtualMachine.metadata?.labels || {},
		disks: virtualMachine.disks || [],
	} as UpdateVirtualMachineRequest;

	let request = $state(defaults);

	function reset() {
		request = defaults;
		labelKey = '';
		labelValue = '';
	}

	// ==================== Label Management ====================
	function addLabel() {
		if (labelKey.trim() && labelValue.trim()) {
			request.labels = { ...request.labels, [labelKey.trim()]: labelValue.trim() };
			labelKey = '';
			labelValue = '';
		}
	}
	function removeLabel(key: string) {
		const { [key]: _, ...rest } = request.labels;
		request.labels = rest;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="outline">
		<Icon icon="ph:pencil" />
		{m.edit()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.edit()} {m.virtual_machine()}</Modal.Header>
		<Form.Root>
			<!-- ==================== Basic Configuration ==================== -->
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.virtual_machine_name()}</Form.Label>
					<SingleInput.General required readonly bind:value={request.name} />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.network_name()}</Form.Label>
					<SingleInput.General bind:value={request.networkName} />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.labels()}</Form.Label>
					<div class="space-y-2">
						<div class="flex gap-2">
							<SingleInput.General
								type="text"
								placeholder={m.label_key()}
								bind:value={labelKey}
								class="flex-1"
							/>
							<SingleInput.General
								type="text"
								placeholder={m.label_value()}
								bind:value={labelValue}
								class="flex-1"
							/>
							<Button
								type="button"
								variant="outline"
								size="sm"
								disabled={!labelKey.trim() || !labelValue.trim()}
								onclick={addLabel}
							>
								<Icon icon="ph:plus" class="size-4" />
								Add
							</Button>
						</div>
						{#if Object.keys(request.labels).length > 0}
							<div class="space-y-1">
								{#each Object.entries(request.labels) as [key, value]}
									<div class="bg-muted flex items-center justify-between rounded-md px-3 py-2">
										<span class="text-sm">
											<span class="font-medium">{key}</span>: {value}
										</span>
										<Button
											type="button"
											variant="ghost"
											size="sm"
											onclick={() => removeLabel(key)}
										>
											<Icon icon="ph:x" class="size-4" />
										</Button>
									</div>
								{/each}
							</div>
						{/if}
					</div>
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
						toast.promise(() => KubeVirtClient.updateVirtualMachine(request), {
							loading: `Updating ${virtualMachine.metadata?.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully updated ${virtualMachine.metadata?.name}`;
							},
							error: (error) => {
								let message = `Failed to update ${virtualMachine.metadata?.name}`;
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
					{m.save()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
