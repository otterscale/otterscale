<script lang="ts" module>
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { writable, type Writable } from 'svelte/store';
	import type { SubvolumeGroup } from './types';

	export type Request = {
		name: string;
		volumeName: string;
		size: number;
		pool: string;
		uid: string;
		gid: string;
		mode: string[];
	};

	export const pools: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'ceph-fs_data',
			label: 'Ceph File System Data',
			icon: 'ph:cube'
		}
	]);
</script>

<script lang="ts">
	let { data = $bindable() }: { data: Writable<SubvolumeGroup[]> } = $props();

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
			Create Subvolume
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

				<Form.Field>
					<Form.Label>Volume Name</Form.Label>
					<SingleInput.General type="text" bind:value={request.volumeName} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Subvolume Group</Form.Legend>

				<Form.Field>
					<Form.Label>Size</Form.Label>
					<SingleInput.General type="number" bind:value={request.size} />
				</Form.Field>
				<Form.Help>
					The size of the subvolume is specified by setting a quota on it. If left blank or put 0,
					then quota will be infinite
				</Form.Help>

				<Form.Field>
					<Form.Label>Pool</Form.Label>
					<SingleSelect.Root options={pools} bind:value={request.pool}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>No results found.</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $pools as option}
											<SingleSelect.Item {option}>
												<Icon
													icon={option.icon ? option.icon : 'ph:empty'}
													class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
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
				<Form.Help>By default, the data_pool_layout of the parent directory is selected.</Form.Help>

				<Form.Field>
					<Form.Label>UID</Form.Label>
					<SingleInput.General type="text" bind:value={request.uid} />
				</Form.Field>

				<Form.Field>
					<Form.Label>GID</Form.Label>
					<SingleInput.General type="text" bind:value={request.gid} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						data.set([]);
						console.log($data);
						stateController.close();
					}}
				>
					Create
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
