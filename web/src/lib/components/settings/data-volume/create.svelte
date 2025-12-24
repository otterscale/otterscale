<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type {
		CreateDataVolumeRequest,
		DataVolume_Source
	} from '$lib/api/instance/v1/instance_pb';
	import { DataVolume_Source_Type, InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';

	let { scope, reloadManager }: { scope: string; reloadManager: ReloadManager } = $props();

	// Context dependencies
	const transport: Transport = getContext('transport');
	const virtualMachineClient = createClient(InstanceService, transport);

	// ==================== State Variables ====================
	let request: CreateDataVolumeRequest = $state({} as CreateDataVolumeRequest);

	let invalidity = $state({} as Booleanified<CreateDataVolumeRequest>);
	const invalid = $derived(invalidity.name || invalidity.sizeBytes);

	let open = $state(false);
	function init() {
		request = {
			scope: scope,
			name: '',
			namespace: 'kubevirt',
			source: { type: DataVolume_Source_Type.BLANK_IMAGE, data: '' } as DataVolume_Source,
			bootImage: false,
			sizeBytes: BigInt(10 * 1024 ** 3)
		} as CreateDataVolumeRequest;
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
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_data_volume()}</Modal.Header>
		<Form.Root>
			<!-- ==================== Basic Configuration ==================== -->
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.GeneralRule
						required
						type="text"
						bind:value={request.name}
						bind:invalid={invalidity.name}
						validateRule="rfc1123"
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.namespace()}</Form.Label>
					<SingleInput.General disabled type="text" bind:value={request.namespace} />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.size()}</Form.Label>
					<SingleInput.Measurement
						required
						bind:value={request.sizeBytes}
						bind:invalid={invalidity.sizeBytes}
						transformer={(value) => (typeof value === 'number' ? BigInt(value) : undefined)}
						units={[
							{ value: 1024 ** 3, label: 'GB' } as SingleInput.UnitType,
							{ value: 1024 ** 4, label: 'TB' } as SingleInput.UnitType
						]}
					/>
				</Form.Field>
				<Form.Field>
					<div class="flex items-center justify-between">
						<Form.Label class="whitespace-nowrap">{m.boot_image()}</Form.Label>
						<SingleInput.Boolean
							bind:value={request.bootImage}
							descriptor={() => ''}
							onCheckedChange={(checked) => {
								if (request.source) {
									request.source.type = checked
										? DataVolume_Source_Type.HTTP_URL
										: DataVolume_Source_Type.BLANK_IMAGE;
									if (!checked) {
										request.source.data = '';
									}
								}
							}}
						/>
					</div>
				</Form.Field>
				{#if request.bootImage && request.source}
					<Form.Label>{m.source()}</Form.Label>
					<SingleInput.General
						type="text"
						bind:value={request.source.data}
						placeholder="https://cloud-images.ubuntu.com/xxx/xxx/xxx.img"
					/>
					<div class="flex justify-end gap-2">
						<Button
							variant="outline"
							size="sm"
							href="https://cloud-images.ubuntu.com/"
							target="_blank"
							class="flex items-center gap-1"
						>
							<Icon icon="ph:arrow-square-out" />
							{m.cloud_image()}
						</Button>
					</div>
				{/if}
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
						toast.promise(() => virtualMachineClient.createDataVolume(request), {
							loading: `Creating ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully created ${request.name}`;
							},
							error: (error) => {
								let message = `Failed to create ${request.name}`;
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
