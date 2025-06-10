<script lang="ts">
	import Icon from '@iconify/svelte';
	import { Single as SingleInput, Multiple as MultipleInput } from '$lib/components/custom/input';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Form from '$lib/components/custom/form';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { writable, type Writable } from 'svelte/store';

	type Request = {
		name: string;
		placement: any;
		hosts: string[];
		labels: string[];
	};

	const placements: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'host',
			label: 'Host',
			icon: 'ph:computer-tower'
		},
		{
			value: 'label',
			label: 'Label',
			icon: 'ph:tag'
		}
	]);

	const DEFAULT_REQUEST = {} as Request;
	let request: Request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<div class="flex justify-end">
		<AlertDialog.Trigger class={cn(buttonVariants({ variant: 'default' }))}>
			<div class="flex items-center gap-2">
				<Icon icon="ph:plus" />
				<p class="text-base">Create</p>
			</div>
		</AlertDialog.Trigger>
	</div>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Create File System
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
					<Form.Label for="filesystem-placement">Placement</Form.Label>
					<SingleSelect.Root options={placements} bind:value={request.placement}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>No results found.</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $placements as placement}
											<SingleSelect.Item option={placement}>
												<Icon
													icon={placement.icon ? placement.icon : 'ph:empty'}
													class={cn('size-5', placement.icon ? 'visibale' : 'invisible')}
												/>
												{placement.label}
												<SingleSelect.Check option={placement} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>
				{#if request.placement === 'host'}
					<Form.Field>
						<Form.Label for="filesystem-hosts">Hosts</Form.Label>
						<MultipleInput.Root
							type="text"
							bind:values={request.hosts}
							id="filesystem-hosts"
							contextData={{ icon: 'ph:computer-tower' }}
						>
							<MultipleInput.Viewer />
							<MultipleInput.Controller>
								<MultipleInput.Input />
								<MultipleInput.Add />
								<MultipleInput.Clear />
							</MultipleInput.Controller>
						</MultipleInput.Root>
					</Form.Field>
				{:else if request.placement === 'label'}
					<Form.Field>
						<Form.Label for="filesystem-labels">Labels</Form.Label>
						<MultipleInput.Root
							type="text"
							bind:values={request.labels}
							id="filesystem-labels"
							contextData={{ icon: 'ph:tag' }}
						>
							<MultipleInput.Viewer />
							<MultipleInput.Controller>
								<MultipleInput.Input />
								<MultipleInput.Add />
								<MultipleInput.Clear />
							</MultipleInput.Controller>
						</MultipleInput.Root>
					</Form.Field>
				{/if}
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
