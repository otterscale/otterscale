<script lang="ts" module>
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { Writable } from 'svelte/store';
	import type { Role } from './types';

	export type Request = {
		rolename: string;
		path: string;
		assumeRolePolicyDocument: string;
	};
</script>

<script lang="ts">
	let { data = $bindable() }: { data: Writable<Role[]> } = $props();

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
			Create Role
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label for="filesystem-name">Role Name</Form.Label>
					<SingleInput.General required type="text" bind:value={request.rolename} />
				</Form.Field>
				<Form.Field>
					<Form.Label for="filesystem-name">Path</Form.Label>
					<SingleInput.General required type="text" bind:value={request.path} />
				</Form.Field>
				<Form.Field>
					<Form.Label for="filesystem-name">Assume Role Policy Document</Form.Label>
					<SingleInput.General required type="text" bind:value={request.assumeRolePolicyDocument} />
				</Form.Field>
				<Form.Help>
					Paste a json assume role policy document, to find more information on how to get this
					document, click here.
				</Form.Help>
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
