<script lang="ts" module>
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import Icon from '@iconify/svelte';
	import { type Request } from '../create.svelte';
	import type { BlockImageSnapshot } from './types';
</script>

<script lang="ts">
	let { blockImagesnapshot }: { blockImagesnapshot: BlockImageSnapshot } = $props();

	const DEFAULT_REQUEST = {} as Request;

	let request: Request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class="text-destructive flex h-full w-full items-center gap-2">
		<Icon icon="ph:trash" />
		Delete
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Delete RADOS Block Device Snapshot
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Help>
					Please type the RADOS block device snapshot name exactly to confirm deletion. This action
					cannot be undone.
				</Form.Help>
				<Form.Field>
					<SingleInput.DeletionConfirm
						required
						id="rados-block-device-snapshot-delete"
						target={blockImagesnapshot.name}
						bind:value={request.name}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						console.log(request);
					}}
				>
					Delete
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
