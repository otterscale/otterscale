<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

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
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { configuration }: { configuration: Writable<Configuration> } = $props();

	const transport: Transport = getContext('transport');
	let distroSeriesOptions = $state(writable<SingleSelect.OptionType[]>([]));
	let distroSeriesArchitecturesMap: Record<string, Writable<SingleSelect.OptionType[]>> = {};
	const client = createClient(ConfigurationService, transport);
	const defaults = {} as CreateBootImageRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	const architecturesOptions = $derived(distroSeriesArchitecturesMap[request.distroSeries]);
	$effect(() => {
		if (request.distroSeries) {
			request.architectures = [];
		}
	});

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
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_boot_image()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.distro_series()}</Form.Label>
					<SingleSelect.Root options={distroSeriesOptions} bind:value={request.distroSeries}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
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
						<Form.Label>{m.architecture()}</Form.Label>
						<MultipleSelect.Root bind:value={request.architectures} options={architecturesOptions}>
							<MultipleSelect.Viewer />
							<MultipleSelect.Controller>
								<MultipleSelect.Trigger />
								<MultipleSelect.Content>
									<MultipleSelect.Options>
										<MultipleSelect.Input />
										<MultipleSelect.List>
											<MultipleSelect.Empty>{m.no_result()}</MultipleSelect.Empty>
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
											<MultipleSelect.ActionAll>{m.all()}</MultipleSelect.ActionAll>
											<MultipleSelect.ActionClear>{m.clear()}</MultipleSelect.ActionClear>
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
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					onclick={() => {
						const distroSeries = request.distroSeries;
						const architectures = `${request.architectures.join(', ')}`;
						toast.promise(() => client.createBootImage(request), {
							loading: 'Loading...',
							success: () => {
								client.getConfiguration({}).then((response) => {
									configuration.set(response);
								});
								return `Create boot images ${distroSeries}: ${architectures} success`;
							},
							error: (error) => {
								let message = `Fail to create boot images ${distroSeries}: ${architectures}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						reset();
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
