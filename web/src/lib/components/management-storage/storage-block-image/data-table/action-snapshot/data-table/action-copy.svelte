<script lang="ts" module>
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { type Request } from './create.svelte';
	import type { Snapshot } from './types';
</script>

<script lang="ts">
	let { blockImagesnapshot }: { blockImagesnapshot: Snapshot } = $props();

	const DEFAULT_REQUEST = { name: blockImagesnapshot.name } as Request;

	let request: Request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class={cn('flex h-full w-full items-center gap-2')}>
		<Icon icon="ph:copy" />
		Copy
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Create RADOS Block Device Snapshot
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label for="rados-block-device-snapshot-name">Name</Form.Label>
					<SingleInput.General
						required
						type="text"
						id="rados-block-device-snapshot-name"
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
					Create
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
