<script lang="ts">
	import type { Pool, UpdatePoolRequest } from '$gen/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { QUOTAS_BYTES_HELP_TEXT, QUOTAS_OBJECTS_HELP_TEXT } from './helper';

	let { pool }: { pool: Pool } = $props();

	const DEFAULT_REQUEST = { poolName: pool.name } as UpdatePoolRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class={cn('flex h-full w-full items-center gap-2')}>
		<Icon icon="ph:pencil" />
		Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Create Pool
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General required type="text" bind:value={request.poolName} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Quotas</Form.Legend>

				<Form.Field>
					<Form.Label>Bytes</Form.Label>
					<SingleInput.General type="number" bind:value={request.quotaBytes} />
				</Form.Field>
				<Form.Help>
					{QUOTAS_BYTES_HELP_TEXT}
				</Form.Help>

				<Form.Field>
					<Form.Label>Objects</Form.Label>
					<SingleInput.General type="number" bind:value={request.quotaObjects} />
				</Form.Field>
				<Form.Help>
					{QUOTAS_OBJECTS_HELP_TEXT}
				</Form.Help>
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
					Edit
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
