<script lang="ts">
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import {
		Multiple as MultipleSelect,
		Single as SingleSelect
	} from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { writable, type Writable } from 'svelte/store';
	import type { FileSystem } from './types';

	type Request = {
		userId: string;
		directory: string;
		permissions: string[];
	};

	const permissions: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'read',
			label: 'Read',
			icon: 'ph:lock',
			information: 'Read permission is the minimum givable access',
			enabled: true
		},
		{
			value: 'write',
			label: 'Write',
			icon: 'ph:lock',
			enabled: false
		},
		{
			value: 'quota',
			label: 'Quota',
			icon: 'ph:lock',
			information: 'Permission to set layouts or quotas, write access needed',
			enabled: true
		},
		{
			value: 'snapshot',
			label: 'Snapshot',
			icon: 'ph:lock',
			information: 'Permission to create or delete snapshots, write access needed',
			enabled: true
		},
		{
			value: 'root_squash',
			label: 'Root Squash',
			icon: 'ph:lock',
			information: 'Safety measure to prevent scenarios such as accidental sudo rm -rf /path',
			enabled: false
		}
	]);

	const DEFAULT_REQUEST = { permissions: ['write', 'root_squash'] } as Request;
	let request: Request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	let { fileSystem }: { fileSystem: FileSystem } = $props();

	const stateController = new DialogStateController(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class={cn('flex h-full w-full items-center gap-2')}>
		<Icon icon="ph:user" />
		Authorize
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Update Access
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label for="filesystem-name">Name</Form.Label>
					<SingleInput.General disabled type="text" id="filesystem-name" value={fileSystem.name} />
				</Form.Field>
				<Form.Field>
					<Form.Label for="filesystem-user">User ID</Form.Label>
					<SingleInput.General required type="text" id="filesystem-user" value={request.userId} />
					<Form.Help>You can manage users from Ceph Users page</Form.Help>
				</Form.Field>
				<Form.Field>
					<Form.Label for="filesystem-directory">Directory</Form.Label>
					<SingleInput.General
						required
						type="text"
						id="filesystem-directory"
						value={request.directory}
					/>
					<Form.Help>Path to restrict access to</Form.Help>
				</Form.Field>
				<Form.Field>
					<Form.Label for="filesystem-directory">Permissions</Form.Label>
					<MultipleSelect.Root bind:value={request.permissions} options={permissions}>
						<MultipleSelect.Viewer />
						<MultipleSelect.Controller>
							<MultipleSelect.Trigger />
							<MultipleSelect.Content>
								<MultipleSelect.Options>
									<MultipleSelect.Input />
									<MultipleSelect.List>
										<MultipleSelect.Empty>No results found.</MultipleSelect.Empty>
										<MultipleSelect.Group>
											{#each $permissions as permission}
												<MultipleSelect.Item option={permission} disabled={!permission.enabled}>
													<Icon
														icon={permission.icon ? permission.icon : 'ph:empty'}
														class={cn('size-5', permission.icon ? 'visibale' : 'invisible')}
													/>
													{permission.label}
													<MultipleSelect.ItemInformation>
														{permission.information}
													</MultipleSelect.ItemInformation>
													<MultipleSelect.Check option={permission} />
												</MultipleSelect.Item>
											{/each}
										</MultipleSelect.Group>
									</MultipleSelect.List>
								</MultipleSelect.Options>
							</MultipleSelect.Content>
						</MultipleSelect.Controller>
					</MultipleSelect.Root>
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
