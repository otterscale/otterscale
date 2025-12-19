<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		type Configuration,
		ConfigurationService,
		type CreateBootImageRequest,
		type UpdateBootImageRequest
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Form from '$lib/components/custom/form';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import {
		Multiple as MultipleSelect,
		Single as SingleSelect
	} from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let { configuration }: { configuration: Writable<Configuration> } = $props();

	const transport: Transport = getContext('transport');
	const distroSeriesOptions = writable<SingleSelect.OptionType[]>([]);
	const client = createClient(ConfigurationService, transport);

	let request = $state({} as CreateBootImageRequest);
	let distroSeriesArchitecturesMap: Record<string, Writable<SingleSelect.OptionType[]>> = {};
	let open = $state(false);

	const architecturesOptions = $derived(distroSeriesArchitecturesMap[request.distroSeries]);
	const existingBootImage = $derived(
		$configuration.bootImages.find((img) => img.distroSeries === request.distroSeries)
	);

	$effect(() => {
		if (request.distroSeries) {
			if (existingBootImage) {
				request.architectures = [...existingBootImage.architectures];
			} else {
				request.architectures = [];
			}
		}
	});

	onMount(async () => {
		try {
			const response = await client.listBootImageSelections({});
			distroSeriesOptions.set(
				response.bootImageSelections
					.map((bootImageSelection) => ({
						value: bootImageSelection.distroSeries,
						label: bootImageSelection.name,
						icon: 'ph:disc'
					}))
					.sort((a, b) => b.label.localeCompare(a.label))
			);
			distroSeriesArchitecturesMap = Object.fromEntries(
				response.bootImageSelections.map((bootImageSelection) => [
					bootImageSelection.distroSeries,
					writable(
						bootImageSelection.architectures.map((architecture) => ({
							value: architecture,
							label: architecture,
							icon: 'ph:cpu'
						}))
					)
				])
			);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});

	function init() {
		request = {} as CreateBootImageRequest;
	}
	function close() {
		open = false;
	}
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
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
										{#each $distroSeriesOptions as option (option.value)}
											<SingleSelect.Item {option}>
												<Icon
													icon={option.icon ? option.icon : 'ph:empty'}
													class={cn('size-5', option.icon ? 'visible' : 'invisible')}
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
												{#each $architecturesOptions as option (option.value)}
													<MultipleSelect.Item {option}>
														<Icon
															icon={option.icon ? option.icon : 'ph:empty'}
															class={cn('size-5', option.icon ? 'visible' : 'invisible')}
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
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					onclick={() => {
						const architectures = `${request.architectures.join(', ')}`;

						if (existingBootImage) {
							const updateRequest = {
								id: existingBootImage.id,
								distroSeries: request.distroSeries,
								architectures: request.architectures
							} as UpdateBootImageRequest;

							toast.promise(() => client.updateBootImage(updateRequest), {
								loading: 'Updating...',
								success: () => {
									client.getConfiguration({}).then((response) => {
										configuration.set(response);
									});
									return `Update boot images ${request.distroSeries}: ${architectures} success`;
								},
								error: (error) => {
									let message = `Fail to update boot images ${request.distroSeries}: ${architectures}`;
									toast.error(message, {
										description: (error as ConnectError).message.toString(),
										duration: Number.POSITIVE_INFINITY
									});
									return message;
								}
							});
						} else {
							toast.promise(() => client.createBootImage(request), {
								loading: 'Loading...',
								success: () => {
									client.getConfiguration({}).then((response) => {
										configuration.set(response);
									});
									return `Create boot images ${request.distroSeries}: ${architectures} success`;
								},
								error: (error) => {
									let message = `Fail to create boot images ${request.distroSeries}: ${architectures}`;
									toast.error(message, {
										description: (error as ConnectError).message.toString(),
										duration: Number.POSITIVE_INFINITY
									});
									return message;
								}
							});
						}
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
