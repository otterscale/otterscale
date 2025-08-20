<script lang="ts" module>
	import {
		ConfigurationService,
		type Configuration,
		type CreateBootImageRequest
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Form from '$lib/components/custom/form';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import {
		Multiple as MultipleSelect,
		Single as SingleSelect
	} from '$lib/components/custom/select';
	import { cn } from '$lib/utils';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
	import { writable } from 'svelte/store';
</script>

<script lang="ts">
	let { configuration = $bindable() }: { configuration: Writable<Configuration> } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ConfigurationService, transport);

	const DEFAULT_REQUEST = {} as CreateBootImageRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	let distroSeriesOptions = $state(writable<SingleSelect.OptionType[]>([]));
	let distroSeriesArchitecturesMap: Record<string, Writable<SingleSelect.OptionType[]>> = {};
	const architecturesOptions = $derived(distroSeriesArchitecturesMap[request.distroSeries]);
	$effect(() => {
		request.distroSeries;
		request.architectures = [];
	});

	let open = $state(false);
	function close() {
		open = false;
	}

	let isMounted = false;
	onMount(async () => {
		try {
			await client.listBootImageSelections({}).then((response) => {
				distroSeriesOptions.set(
					response.bootImageSelections.map((bootImageSelection) => ({
						value: bootImageSelection.distroSeries,
						label: bootImageSelection.name,
						icon: 'ph:empty'
					}))
				);
				distroSeriesArchitecturesMap = Object.fromEntries(
					response.bootImageSelections.map((bootImageSelection) => [
						bootImageSelection.distroSeries,
						writable(
							bootImageSelection.architectures.map((architecture) => ({
								value: architecture,
								label: architecture,
								icon: 'ph:empty'
							}))
						)
					])
				);
			});

			isMounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		Boot Image
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Create Boot Image</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Distro Series</Form.Label>
					<SingleSelect.Root options={distroSeriesOptions} bind:value={request.distroSeries}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>No results found.</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $distroSeriesOptions as option}
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

				{#if request.distroSeries}
					<Form.Field>
						<Form.Label>Architectures</Form.Label>
						<MultipleSelect.Root bind:value={request.architectures} options={architecturesOptions}>
							<MultipleSelect.Viewer />
							<MultipleSelect.Controller>
								<MultipleSelect.Trigger />
								<MultipleSelect.Content>
									<MultipleSelect.Options>
										<MultipleSelect.Input />
										<MultipleSelect.List>
											<MultipleSelect.Empty>No results found.</MultipleSelect.Empty>
											<MultipleSelect.Group>
												{#each $architecturesOptions as option}
													<MultipleSelect.Item {option}>
														<Icon
															icon={option.icon ? option.icon : 'ph:empty'}
															class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
														/>
														{option.label}
														<MultipleSelect.Check {option} />
													</MultipleSelect.Item>
												{/each}
											</MultipleSelect.Group>
										</MultipleSelect.List>
										<MultipleSelect.Actions>
											<MultipleSelect.ActionAll>All</MultipleSelect.ActionAll>
											<MultipleSelect.ActionClear>Clear</MultipleSelect.ActionClear>
										</MultipleSelect.Actions>
									</MultipleSelect.Options>
								</MultipleSelect.Content>
							</MultipleSelect.Controller>
						</MultipleSelect.Root>
					</Form.Field>
				{/if}
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel onclick={reset}>Cancel</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					onclick={() => {
						toast.promise(() => client.createBootImage(request), {
							loading: 'Loading...',
							success: () => {
								client.getConfiguration({}).then((response) => {
									configuration.set(response);
								});
								return `Create boot images ${request.distroSeries}: ${request.architectures.join(', ')} success`;
							},
							error: (error) => {
								let message = `Fail to create boot images ${request.distroSeries}: ${request.architectures.join(', ')}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});

						reset();
						close();
					}}>Create</Modal.Action
				>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
