<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Single as SingleInput, Multiple as MultipleInput } from '$lib/components/custom/input';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Form from '$lib/components/custom/form';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { writable, type Writable } from 'svelte/store';
	import type { Snapshot } from './types';

	export type Request = {
		name: string;
	};
	export const subvolumnes: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'not_allow_access',
			label: 'Not Allow Access',
			icon: 'ph:cube'
		}
	]);
	export const subvolumneGroups: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: '_nogroup',
			label: '_nogroup',
			icon: 'ph:cube'
		}
	]);
</script>

<script lang="ts">
	let { data = $bindable() }: { data: Writable<Snapshot[]> } = $props();

	const DEFAULT_REQUEST = {} as Request;
	let request: Request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<div class="flex justify-end">
		<AlertDialog.Trigger class={cn(buttonVariants({ variant: 'default', size: 'sm' }))}>
			<div class="flex items-center gap-1">
				<Icon icon="ph:plus" />
				Create
			</div>
		</AlertDialog.Trigger>
	</div>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Create Snapshot
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label for="filesystem-name">Name</Form.Label>
					<SingleInput.General
						required
						type="text"
						id="filesystem-name"
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
						data.set([]);
						console.log(request);
						stateController.close();
					}}
				>
					Create
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
