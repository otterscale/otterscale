<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { CreateSubvolumeRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages.js';
</script>

<script lang="ts">
	let {
		scope,
		volume,
		group,
		reloadManager
	}: {
		scope: string;
		volume: string;
		group: string;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');

	let open = $state(false);
	function close() {
		open = false;
	}
	let invalid = $state(false);
	const storageClient = createClient(StorageService, transport);

	let request = $state({} as CreateSubvolumeRequest);
	function init() {
		request = {
			scope: scope,
			volumeName: volume,
			groupName: group,
			export: true
		} as CreateSubvolumeRequest;
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
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header class="flex items-center justify-center text-xl font-bold">
			{m.create_nfs()}
		</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.GeneralRule
						required
						type="text"
						bind:value={request.subvolumeName}
						bind:invalid
						validateRule="lower-alphanum-dash-start-alpha"
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.group()}</Form.Label>
					<SingleInput.General disabled type="text" value={group === '' ? 'default' : group} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.quota_size()}</Form.Label>
					<Form.Help>
						{m.nfs_quota_size_direction()}
					</Form.Help>
					<SingleInput.Measurement
						bind:value={request.quotaBytes}
						transformer={(value) => (value !== undefined ? BigInt(value) : undefined)}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
						]}
					/>
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
						toast.promise(() => storageClient.createSubvolume(request), {
							loading: `Creating ${request.subvolumeName}...`,
							success: () => {
								reloadManager.force();
								return `Create ${request.subvolumeName}`;
							},
							error: (error) => {
								let message = `Fail to create ${request.subvolumeName}`;
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
