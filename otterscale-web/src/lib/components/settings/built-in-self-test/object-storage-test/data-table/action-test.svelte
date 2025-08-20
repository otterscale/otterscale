<script lang="ts">
	import type {
		CreateTestResultRequest,
		ExternalObjectService,
		InternalObjectService,
		TestResult,
		Warp,
		Warp_Input
	} from '$lib/api/bist/v1/bist_pb';
	import { BISTService, Warp_Input_Operation } from '$lib/api/bist/v1/bist_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as MultipleStepModal from '$lib/components/custom/modal/mutiple-step';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { buttonVariants } from '$lib/components/ui/button';
	import { formatCapacity, formatSecond } from '$lib/formatter';
	import { cn } from '$lib/utils';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable, type Writable } from 'svelte/store';
	import ObjectServicesPicker from '../../utils/object-services-picker.svelte';

	// WARP Target
	const warpTarget: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'internalObjectService',
			label: 'Internal Object Service'
		},
		{
			value: 'externalObjectService',
			label: 'External Object Service'
		}
	]);

	// WARP AccessMode
	const Options: SingleSelect.OptionType[] = Object.keys(Warp_Input_Operation)
		.filter((key) => isNaN(Number(key)))
		.map((key) => ({
			value: Warp_Input_Operation[key as keyof typeof Warp_Input_Operation],
			label: key
		}));
	const warpInputOperation: Writable<SingleSelect.OptionType[]> = writable(Options);

	let {
		testResult
	}: {
		testResult?: TestResult;
	} = $props();

	// Request
	const DEFAULT_WARP_REQUEST = testResult
		? ({
				target: {
					value: testResult.kind.value?.target.value,
					case: testResult.kind.value?.target.case
				}
			} as Warp)
		: ({ target: { value: {}, case: {} } } as Warp);
	const DEFAULT_REQUEST = {
		kind: { value: DEFAULT_WARP_REQUEST, case: 'warp' },
		createdBy: 'Woody Lin'
	} as CreateTestResultRequest;
	const DEFAULT_INTERNAL_OBJECT_SERVICE =
		testResult && testResult.kind.value?.target?.case === 'internalObjectService'
			? (testResult.kind.value.target.value as InternalObjectService)
			: ({} as InternalObjectService);
	const DEFAULT_DEFAULT_EXTERNAL_OBJECT_SERVICE =
		testResult && testResult.kind.value?.target?.case === 'externalObjectService'
			? (testResult.kind.value.target.value as ExternalObjectService)
			: ({} as ExternalObjectService);
	const DEFAULT_WARP_INPUT =
		testResult && testResult.kind.value?.input
			? (testResult.kind.value.input as Warp_Input)
			: ({
					durationSeconds: 60,
					objectSizeBytes: 4 * 1024 * 1024,
					objectCount: 500
				} as unknown as Warp_Input);

	let request: CreateTestResultRequest = $state(DEFAULT_REQUEST);
	let requestWarp: Warp = $state(DEFAULT_WARP_REQUEST);
	let requestInternalObjectService: InternalObjectService = $state(DEFAULT_INTERNAL_OBJECT_SERVICE);
	let requestExternalObjectService: ExternalObjectService = $state(
		DEFAULT_DEFAULT_EXTERNAL_OBJECT_SERVICE
	);
	let warpOperation = $state(DEFAULT_WARP_INPUT.operation);
	let warpDuration = $state(DEFAULT_WARP_INPUT.durationSeconds);
	let warpObjectSize = $state(DEFAULT_WARP_INPUT.objectSizeBytes);
	let warpObjectCount = $state(DEFAULT_WARP_INPUT.objectCount);
	function reset() {
		request = DEFAULT_REQUEST;
		requestWarp = DEFAULT_WARP_REQUEST;
		requestInternalObjectService = DEFAULT_INTERNAL_OBJECT_SERVICE;
		requestExternalObjectService = DEFAULT_DEFAULT_EXTERNAL_OBJECT_SERVICE;
		warpOperation = DEFAULT_WARP_INPUT.operation;
		warpDuration = DEFAULT_WARP_INPUT.durationSeconds;
		warpObjectSize = DEFAULT_WARP_INPUT.objectSizeBytes;
		warpObjectCount = DEFAULT_WARP_INPUT.objectCount;
	}

	// Modal state
	let open = $state(false);
	function close() {
		open = false;
	}

	// grpc
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	const client = createClient(BISTService, transport);

	let invalidName = $state(false);
	let invalidTarget = $state(false);
	let invalidEndPoint = $state(false);
	let invalidAccessKey = $state(false);
	let invalidSecretKey = $state(false);
	let invalidOperation = $state(false);

	const invalidBasic = $derived(
		invalidName ||
			invalidTarget ||
			(requestWarp.target.case == 'externalObjectService' &&
				(invalidEndPoint || invalidAccessKey || invalidSecretKey))
	);
	const invalidAdvanced = $derived(invalidOperation);
	const invalid = $derived(invalidBasic || invalidAdvanced);
</script>

<MultipleStepModal.Root bind:open steps={3}>
	{#if testResult}
		<MultipleStepModal.Trigger class="flex h-full w-full items-center gap-2">
			<Icon icon="ph:play" />
			Retest
		</MultipleStepModal.Trigger>
	{:else}
		<div class="flex justify-end">
			<MultipleStepModal.Trigger class={cn(buttonVariants({ variant: 'default', size: 'sm' }))}>
				<div class="flex items-center gap-1">
					<Icon icon="ph:plus" />
					Create
				</div>
			</MultipleStepModal.Trigger>
		</div>
	{/if}
	<MultipleStepModal.Content>
		<MultipleStepModal.Header>Built-in Self Test</MultipleStepModal.Header>
		<MultipleStepModal.Stepper>
			<MultipleStepModal.Steps>
				<MultipleStepModal.Step icon="ph:number-one">
					{#snippet text()}
						Basic
					{/snippet}
				</MultipleStepModal.Step>
				<MultipleStepModal.Step icon="ph:number-two">
					{#snippet text()}
						Advance
					{/snippet}
				</MultipleStepModal.Step>
				<MultipleStepModal.Step icon="ph:number-three">
					{#snippet text()}
						Check
					{/snippet}
				</MultipleStepModal.Step>
			</MultipleStepModal.Steps>
			<MultipleStepModal.Models>
				<!-- Step One -->
				<MultipleStepModal.Model>
					<Form.Root class="max-h-[65vh]">
						<Form.Fieldset>
							<!-- Name -->
							<Form.Field>
								<Form.Label for="bist-name">Name</Form.Label>
								<SingleInput.General
									type="text"
									id="name"
									required
									bind:value={request.name}
									bind:invalid={invalidName}
								/>
							</Form.Field>
							<!-- Choose Target -->
							<Form.Field>
								<Form.Label for="bist-input">Target</Form.Label>
								<SingleSelect.Root
									options={warpTarget}
									required
									bind:value={requestWarp.target.case}
									bind:invalid={invalidTarget}
								>
									<SingleSelect.Trigger />
									<SingleSelect.Content>
										<SingleSelect.Options>
											<SingleSelect.Input />
											<SingleSelect.List>
												<SingleSelect.Empty>No results found.</SingleSelect.Empty>
												<SingleSelect.Group>
													{#each $warpTarget as item}
														<SingleSelect.Item option={item}>
															<Icon
																icon={item.icon ? item.icon : 'ph:empty'}
																class={cn('size-5', item.icon ? 'visible' : 'invisible')}
															/>
															{item.label}
															<SingleSelect.Check option={item} />
														</SingleSelect.Item>
													{/each}
												</SingleSelect.Group>
											</SingleSelect.List>
										</SingleSelect.Options>
									</SingleSelect.Content>
								</SingleSelect.Root>
							</Form.Field>
						</Form.Fieldset>
						<!-- Target -->
						{#if requestWarp.target.case == 'internalObjectService'}
							<Form.Fieldset>
								<Form.Legend>Target</Form.Legend>
								<Form.Field>
									<Form.Label>Interanl Object Service</Form.Label>
									<ObjectServicesPicker
										bind:selectedInternalObjectService={requestInternalObjectService}
									/>
								</Form.Field>
							</Form.Fieldset>
						{:else if requestWarp.target.case == 'externalObjectService'}
							<Form.Fieldset>
								<Form.Legend>Target</Form.Legend>
								<Form.Field>
									<Form.Label>Endpoint</Form.Label>
									<SingleInput.General
										type="text"
										required
										bind:value={requestExternalObjectService.endpoint}
										bind:invalid={invalidEndPoint}
									/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Access Key</Form.Label>
									<SingleInput.General
										type="text"
										required
										bind:value={requestExternalObjectService.accessKey}
										bind:invalid={invalidAccessKey}
									/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Secret Key</Form.Label>
									<SingleInput.General
										type="text"
										required
										bind:value={requestExternalObjectService.secretKey}
										bind:invalid={invalidSecretKey}
									/>
								</Form.Field>
							</Form.Fieldset>
						{/if}
					</Form.Root>
				</MultipleStepModal.Model>

				<!-- Step two -->
				<MultipleStepModal.Model>
					<Form.Root class="max-h-[65vh]">
						<Form.Fieldset>
							<Form.Legend>Parameter</Form.Legend>
							<!-- warpInputOperation -->
							<Form.Field>
								<Form.Label for="warp-operation">Operation</Form.Label>
								<SingleSelect.Root
									options={warpInputOperation}
									required
									bind:value={warpOperation}
									bind:invalid={invalidOperation}
								>
									<SingleSelect.Trigger />
									<SingleSelect.Content>
										<SingleSelect.Options>
											<SingleSelect.Input />
											<SingleSelect.List>
												<SingleSelect.Empty>No results found.</SingleSelect.Empty>
												<SingleSelect.Group>
													{#each $warpInputOperation as item}
														<SingleSelect.Item option={item}>
															<Icon
																icon={item.icon ? item.icon : 'ph:empty'}
																class={cn('size-5', item.icon ? 'visible' : 'invisible')}
															/>
															{item.label}
															<SingleSelect.Check option={item} />
														</SingleSelect.Item>
													{/each}
												</SingleSelect.Group>
											</SingleSelect.List>
										</SingleSelect.Options>
									</SingleSelect.Content>
								</SingleSelect.Root>
							</Form.Field>
							<!-- Duration -->
							<Form.Field>
								<Form.Label>Duration</Form.Label>
								<SingleInput.Measurement
									bind:value={warpDuration}
									units={[
										{ value: 1, label: 's' } as SingleInput.UnitType,
										{ value: 60, label: 'm' } as SingleInput.UnitType,
										{ value: 3600, label: 'h' } as SingleInput.UnitType,
										{ value: 86400, label: 'd' } as SingleInput.UnitType
									]}
								/>
							</Form.Field>
							<!-- ObjectSize -->
							<Form.Field>
								<Form.Label>Object Size</Form.Label>
								<SingleInput.Measurement
									bind:value={warpObjectSize}
									units={[
										{ value: Math.pow(2, 10 * 0), label: 'B' } as SingleInput.UnitType,
										{ value: Math.pow(2, 10 * 1), label: 'KB' } as SingleInput.UnitType,
										{ value: Math.pow(2, 10 * 2), label: 'MB' } as SingleInput.UnitType,
										{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
										{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType,
										{ value: Math.pow(2, 10 * 5), label: 'PB' } as SingleInput.UnitType
									]}
								/>
							</Form.Field>
							<!-- ObjectCount -->
							<Form.Field>
								<Form.Label>Object Count</Form.Label>
								<SingleInput.General
									type="number"
									placeholder="500"
									bind:value={warpObjectCount}
									disabled={warpOperation == Warp_Input_Operation.PUT}
								/>
							</Form.Field>
						</Form.Fieldset>
					</Form.Root>
				</MultipleStepModal.Model>

				<!-- Step three Overview -->
				<MultipleStepModal.Model>
					<Form.Root>
						<!-- Step 1 -->
						<Form.Fieldset>
							<Form.Legend>Step 1</Form.Legend>
							<Form.Description>Name: {request.name}</Form.Description>
							<Form.Description>Target: {requestWarp.target.case}</Form.Description>
							{#if requestWarp.target.case == 'internalObjectService'}
								<Form.Description>Type: {requestInternalObjectService.type}</Form.Description>
								<Form.Description>Name: {requestInternalObjectService.name}</Form.Description>
								<Form.Description
									>Endpoint: {requestInternalObjectService.endpoint}</Form.Description
								>
							{:else if requestWarp.target.case == 'externalObjectService'}
								<Form.Description
									>Endpoint: {requestExternalObjectService.endpoint}</Form.Description
								>
								<Form.Description
									>Access Key: {requestExternalObjectService.accessKey}</Form.Description
								>
								<Form.Description
									>Secret Key: {requestExternalObjectService.secretKey}</Form.Description
								>
							{/if}
						</Form.Fieldset>
						<!-- Step 2 -->
						<Form.Fieldset>
							{@const duration = formatSecond(Number(warpDuration))}
							{@const objectSize = formatCapacity(Number(warpObjectSize))}
							<Form.Legend>Step 2</Form.Legend>
							<Form.Description>Operation: {Warp_Input_Operation[warpOperation]}</Form.Description>
							<Form.Description>Duration: {duration.value} {duration.unit}</Form.Description>
							<Form.Description>Object Size: {objectSize.value} {objectSize.unit}</Form.Description>
							{#if warpOperation !== Warp_Input_Operation.PUT}
								<Form.Description>Object Count: {warpObjectCount}</Form.Description>
							{/if}
						</Form.Fieldset>
					</Form.Root>
				</MultipleStepModal.Model>
			</MultipleStepModal.Models>
		</MultipleStepModal.Stepper>

		<MultipleStepModal.Footer>
			<MultipleStepModal.Cancel
				onclick={() => {
					reset();
				}}>Cancel</MultipleStepModal.Cancel
			>
			<MultipleStepModal.Controllers>
				<MultipleStepModal.Back>Back</MultipleStepModal.Back>
				<MultipleStepModal.Next>Next</MultipleStepModal.Next>
				<MultipleStepModal.Confirm
					disabled={invalid}
					onclick={() => {
						// prepare request
						if (requestWarp.target.case == 'internalObjectService') {
							requestWarp.target.value = requestInternalObjectService;
						} else if (requestWarp.target.case == 'externalObjectService') {
							requestWarp.target.value = requestExternalObjectService;
						}
						requestWarp.input = {
							operation: warpOperation,
							durationSeconds: BigInt(warpDuration),
							objectSizeBytes: BigInt(warpObjectSize),
							objectCount: warpOperation === Warp_Input_Operation.PUT ? 0n : BigInt(warpObjectCount)
						} as Warp_Input;
						request.kind.value = requestWarp;
						toast.promise(() => client.createTestResult(request), {
							loading: `Testing ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Conduct ${request.name}`;
							},
							error: (error) => {
								let message = `Fail to test ${request.name}`;
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
					Confirm
				</MultipleStepModal.Confirm>
			</MultipleStepModal.Controllers>
		</MultipleStepModal.Footer>
	</MultipleStepModal.Content>
</MultipleStepModal.Root>
