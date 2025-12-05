<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import {
		type CreateModelRequest,
		type Model_Decode,
		Model_Mode,
		type Model_Prefill,
		ModelService
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

	import SelectModel from './util-select-model.svelte';
</script>

<script lang="ts">
	let {
		scope,
		namespace,
		reloadManager
	}: { scope: string; namespace: string; reloadManager: ReloadManager } = $props();

	const transport: Transport = getContext('transport');

	const modelClient = createClient(ModelService, transport);
	const applicationClient = createClient(ApplicationService, transport);

	const defaults = {
		scope: scope,
		namespace: namespace,
		sizeBytes: BigInt(100 * 1024 ** 3),
		maxModelLength: 8192,
		mode: Model_Mode.INTELLIGENT_INFERENCE_SCHEDULING
	} as CreateModelRequest;
	const defaultPrefillResource = {} as Model_Prefill;
	const defaultDecodeResource = { replica: 1, tensor: 1 } as Model_Decode;

	let request = $state(defaults);
	let requestPrefillResource = $state({ ...defaultPrefillResource });
	let requestDecodeResource = $state({ ...defaultDecodeResource });
	let isDisaggregationMode = $state(false);
	let fromLocal: boolean = $state(false);
	let selectedModel: string = $state('');

	function resetPrefillResources() {
		requestPrefillResource = { ...defaultPrefillResource };
	}
	function resetDecodeResources() {
		requestDecodeResource = { ...defaultDecodeResource };
	}
	function resetFromLocal() {
		fromLocal = false;
	}
	function resetSelectedModel() {
		selectedModel = '';
	}

	function reset() {
		request = { ...defaults };
		resetPrefillResources();
		resetDecodeResources();
		resetFromLocal();
		resetSelectedModel();
	}

	function integrate() {
		request.prefill = requestPrefillResource;
		request.decode = requestDecodeResource;
		request.modelName = fromLocal ? `pvc-${selectedModel}` : selectedModel;
		request.mode = isDisaggregationMode
			? Model_Mode.PREFILL_DECODE_DISAGGREGATION
			: Model_Mode.INTELLIGENT_INFERENCE_SCHEDULING;
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

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			reset();
		}
	}}
>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_model()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General
						type="text"
						bind:value={request.name}
						required
						bind:invalid={invalidity.name}
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

				<Form.Field>
					<Form.Label>{m.model_name()}</Form.Label>
					<SelectModel
						bind:value={selectedModel}
						bind:fromLocal
						required
						bind:invalid={invalidity.modelName}
						{scope}
						{namespace}
					/>

					{#if invalidity.modelName}
						<p class="text-sm text-destructive/50">{m.select_model_description()}</p>
					{/if}
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.size()}</Form.Label>
					<SingleInput.Measurement
						type="number"
						bind:value={request.sizeBytes}
						required
						bind:invalid={invalidity.sizeBytes}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 5), label: 'PB' } as SingleInput.UnitType
						]}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.max_model_length()}</Form.Label>
					<SingleInput.General type="number" bind:value={request.maxModelLength} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Field>
				<div class="flex items-center justify-between gap-4">
					<div class="flex items-center gap-2 font-bold">
						<Icon icon="ph:aperture" class="size-5" />
						<p>{m.disaggregation_mode()}</p>
					</div>
					<Switch bind:checked={isDisaggregationMode} />
				</div>
			</Form.Field>
			{#if !isDisaggregationMode}
				<Form.Fieldset>
					<Form.Legend>{m.inference()}</Form.Legend>

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
								min={1}
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
									min={1}
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
							/>
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
									min={1}
									max={100}
									step={1}
									class="w-full py-4 **:data-[slot=slider-track]:h-3"
								/>
							</div>
						</Form.Field>
					</Form.Fieldset>
				</div>
			{/if}
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.Action
				disabled={invalid}
				onclick={() => {
					integrate();
					console.log(request);
					toast.promise(() => modelClient.createModel(request), {
						loading: `Creating ${request.modelName}...`,
						success: () => {
							reloadManager.force();
							return `Create ${request.modelName} success`;
						},
						error: (error) => {
							let message = `Fail to create ${request.modelName}`;
							toast.error(message, {
								description: (error as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY
							});
							return message;
						}
					});
					close();
				}}
			>
				{m.confirm()}
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
