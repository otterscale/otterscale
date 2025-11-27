<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import {
		type CreateModelRequest,
		type Model,
		type Model_Decode,
		Model_Mode,
		type Model_Prefill,
		ModelService,
		type UpdateModelRequest
	} from '$lib/api/model/v1/model_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { Slider } from '$lib/components/ui/slider/index.js';
	import Switch from '$lib/components/ui/switch/switch.svelte';
	import { m } from '$lib/paraglide/messages.js';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		model,
		scope,
		reloadManager
	}: { model: Model; scope: string; reloadManager: ReloadManager } = $props();

	const transport: Transport = getContext('transport');

	const modelClient = createClient(ModelService, transport);
	const applicationClient = createClient(ApplicationService, transport);

	let isIntelligentMode = $state(model.mode === Model_Mode.INTELLIGENT_INFERENCE_SCHEDULING);

	const defaults = {
		scope: scope,
		name: model.name,
		namespace: model.namespace,
		mode: model.mode,
		prefill: { ...model.prefill },
		decode: { ...model.decode }
	} as UpdateModelRequest;
	const defaultPrefillResource = { ...model.prefill } as Model_Prefill;
	const defaultDecodeResource = { ...model.decode } as Model_Decode;

	let request = $state(defaults);
	let requestPrefillResource = $state({ ...defaultPrefillResource });
	let requestDecodeResource = $state({ ...defaultDecodeResource });

	function resetPrefillResources() {
		requestPrefillResource = { ...defaultPrefillResource };
	}
	function resetDecodeResources() {
		requestDecodeResource = { ...defaultDecodeResource };
	}

	function reset() {
		request = { ...defaults };
		resetPrefillResources();
		resetDecodeResources();
	}

	function integrate() {
		request.prefill = requestPrefillResource;
		request.decode = requestDecodeResource;
		request.mode = isIntelligentMode
			? Model_Mode.INTELLIGENT_INFERENCE_SCHEDULING
			: Model_Mode.PREFILL_DECODE_DISAGGREGATION;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	let isNamespaceOptionsLoaded = $state(false);
	const namespaceOptions: Writable<SingleSelect.OptionType[]> = writable([]);
	async function fetchNamespaceOptions() {
		const response = await applicationClient.listNamespaces({ scope });
		namespaceOptions.set(
			response.namespaces.map((namespace) => ({
				value: namespace.name,
				label: namespace.name,
				icon: 'ph:cube'
			}))
		);
		isNamespaceOptionsLoaded = true;
	}

	let invalidity = $state({} as Booleanified<CreateModelRequest>);
	const invalid = $derived(
		invalidity.name || invalidity.namespace || invalidity.modelName || invalidity.sizeBytes
	);

	onMount(async () => {
		try {
			await fetchNamespaceOptions();
		} catch (error) {
			console.debug('Failed to init data:', error);
		}
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		{m.update()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.update_model()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General
						type="text"
						bind:value={request.name}
						required
						bind:invalid={invalidity.name}
						readonly
						class="focus-none"
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.namespace()}</Form.Label>
					{#if isNamespaceOptionsLoaded}
						<SingleSelect.Root
							options={namespaceOptions}
							bind:value={request.namespace}
							required
							bind:invalid={invalidity.namespace}
						>
							<!-- Temporary Disabled -->
							<SingleSelect.Trigger disabled />
							<SingleSelect.Content>
								<SingleSelect.Options>
									<SingleSelect.Input />
									<SingleSelect.List>
										<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
										<SingleSelect.Group>
											{#each $namespaceOptions as option (option.value)}
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
					{/if}
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>{m.prefill()}</Form.Legend>
				<Form.Field>
					<Form.Label>{m.replica()}</Form.Label>
					<SingleInput.General type="number" bind:value={requestPrefillResource.replica} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.memory()}</Form.Label>
					<div class="flex items-center gap-4">
						<p class="w-6 whitespace-nowrap">{requestPrefillResource.vgpumemPercentage} %</p>
						<Slider
							type="single"
							bind:value={requestPrefillResource.vgpumemPercentage}
							max={100}
							step={1}
							class="w-full py-4 **:data-[slot=slider-track]:h-3"
						/>
					</div>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>{m.resources()}</Form.Legend>
				<Form.Field>
					<div class="flex items-center justify-between gap-4">
						<div class="flex items-center gap-2 font-bold">
							<Icon icon={isIntelligentMode ? 'ph:sparkle' : 'ph:aperture'} class="size-5" />
							<p>
								{isIntelligentMode
									? 'Intelligent Optimization Mode Enabled'
									: 'Disaggregation Mode Enabled'}
							</p>
						</div>
						<Switch bind:checked={isIntelligentMode} />
					</div>
				</Form.Field>
				{#if isIntelligentMode}
					<Form.Fieldset>
						<Form.Legend>{m.decode()}</Form.Legend>

						<Form.Field>
							<Form.Label>{m.replica()}</Form.Label>
							<SingleInput.General type="number" bind:value={requestDecodeResource.replica} />
						</Form.Field>

						<Form.Field>
							<Form.Label>{m.tensor()}</Form.Label>
							<SingleInput.General type="number" bind:value={requestDecodeResource.tensor} />
						</Form.Field>

						<Form.Field>
							<Form.Label>{m.memory()}</Form.Label>
							<div class="flex items-center gap-4">
								<p class="w-6 whitespace-nowrap">{requestDecodeResource.vgpumemPercentage} %</p>
								<Slider
									type="single"
									bind:value={requestDecodeResource.vgpumemPercentage}
									max={100}
									step={1}
									class="w-full py-4 **:data-[slot=slider-track]:h-3"
								/>
							</div>
						</Form.Field>
					</Form.Fieldset>
				{:else}
					<div class="flex gap-4">
						<Form.Fieldset>
							<Form.Legend>{m.prefill()}</Form.Legend>
							<Form.Field>
								<Form.Label>{m.replica()}</Form.Label>
								<SingleInput.General type="number" bind:value={requestPrefillResource.replica} />
							</Form.Field>

							<Form.Field>
								<Form.Label>{m.memory()}</Form.Label>
								<div class="flex items-center gap-4">
									<p class="w-6 whitespace-nowrap">{requestPrefillResource.vgpumemPercentage} %</p>
									<Slider
										type="single"
										bind:value={requestPrefillResource.vgpumemPercentage}
										max={100}
										step={1}
										class="w-full py-4 **:data-[slot=slider-track]:h-3"
									/>
								</div>
							</Form.Field>
						</Form.Fieldset>

						<Form.Fieldset>
							<Form.Legend>{m.decode()}</Form.Legend>

							<Form.Field>
								<Form.Label>{m.replica()}</Form.Label>
								<SingleInput.General
									type="number"
									bind:value={requestDecodeResource.replica}
									readonly
									class="focus-none"
								/>
							</Form.Field>

							<Form.Field>
								<Form.Label>{m.tensor()}</Form.Label>
								<SingleInput.General
									type="number"
									bind:value={requestDecodeResource.tensor}
									readonly
									class="focus-none"
								/>
							</Form.Field>

							<Form.Field>
								<Form.Label>{m.memory()}</Form.Label>
								<div class="flex items-center gap-4">
									<p class="w-6 whitespace-nowrap">{requestDecodeResource.vgpumemPercentage} %</p>
									<Slider
										type="single"
										bind:value={requestDecodeResource.vgpumemPercentage}
										max={100}
										step={1}
										class="w-full py-4 **:data-[slot=slider-track]:h-3"
									/>
								</div>
							</Form.Field>
						</Form.Fieldset>
					</div>
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
			<Modal.Action
				disabled={invalid}
				onclick={() => {
					integrate();
					toast.promise(() => modelClient.updateModel(request), {
						loading: `Updating ${request.name}...`,
						success: () => {
							reloadManager.force();
							return `Update ${request.name} success`;
						},
						error: (error) => {
							let message = `Fail to update ${request.name}`;
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
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
