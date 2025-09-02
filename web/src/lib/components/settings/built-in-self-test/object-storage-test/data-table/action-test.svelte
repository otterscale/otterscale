<script lang="ts" module>
	import type {
		CreateTestResultRequest,
		ExternalObjectService,
		InternalObjectService,
		TestResult,
		Warp,
		Warp_Input,
	} from '$lib/api/bist/v1/bist_pb';
	import { BISTService, Warp_Input_Operation } from '$lib/api/bist/v1/bist_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as Modal from '$lib/components/custom/modal/mutiple-step';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { formatCapacity, formatSecond } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
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
			label: 'Internal Object Service',
		},
		{
			value: 'externalObjectService',
			label: 'External Object Service',
		},
	]);

	// WARP AccessMode
	const Options: SingleSelect.OptionType[] = Object.keys(Warp_Input_Operation)
		.filter((key) => isNaN(Number(key)))
		.map((key) => ({
			value: Warp_Input_Operation[key as keyof typeof Warp_Input_Operation],
			label: key,
		}));
	const warpInputOperation: Writable<SingleSelect.OptionType[]> = writable(Options);
</script>

<script lang="ts">
	let {
		testResult,
	}: {
		testResult?: TestResult;
	} = $props();

	// Request
	const DEFAULT_WARP_REQUEST = testResult
		? ({
				target: {
					value: testResult.kind.value?.target.value,
					case: testResult.kind.value?.target.case,
				},
			} as Warp)
		: ({ target: { value: {}, case: {} } } as Warp);
	const DEFAULT_REQUEST = {
		kind: { value: DEFAULT_WARP_REQUEST, case: 'warp' },
		createdBy: 'Woody Lin',
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
					objectCount: 500,
				} as unknown as Warp_Input);

	let request: CreateTestResultRequest = $state(DEFAULT_REQUEST);
	let requestWarp: Warp = $state(DEFAULT_WARP_REQUEST);
	let requestInternalObjectService: InternalObjectService = $state(DEFAULT_INTERNAL_OBJECT_SERVICE);
	let requestExternalObjectService: ExternalObjectService = $state(DEFAULT_DEFAULT_EXTERNAL_OBJECT_SERVICE);
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
				(invalidEndPoint || invalidAccessKey || invalidSecretKey)),
	);
	const invalidAdvanced = $derived(invalidOperation);
	const invalid = $derived(invalidBasic || invalidAdvanced);
</script>

<Modal.Root bind:open steps={3}>
	{#if testResult}
		<Modal.Trigger variant="creative">
			<Icon icon="ph:play" />
			{m.retest()}
		</Modal.Trigger>
	{:else}
		<div class="flex justify-end">
			<Modal.Trigger variant="default">
				<Icon icon="ph:plus" />
				{m.create()}
			</Modal.Trigger>
		</div>
	{/if}
	<Modal.Content>
		<Modal.Header>{m.built_in_self_test()}</Modal.Header>
		<Modal.Stepper>
			<Modal.Steps>
				<Modal.Step icon="ph:number-one">
					{#snippet text()}
						{m.basic()}
					{/snippet}
				</Modal.Step>
				<Modal.Step icon="ph:number-two">
					{#snippet text()}
						{m.advance()}
					{/snippet}
				</Modal.Step>
				<Modal.Step icon="ph:number-three">
					{#snippet text()}
						{m.check()}
					{/snippet}
				</Modal.Step>
			</Modal.Steps>
			<Modal.Models>
				<!-- Step One -->
				<Modal.Model>
					<Form.Root class="max-h-[65vh]">
						<Form.Fieldset>
							<!-- Name -->
							<Form.Field>
								<Form.Label for="bist-name">{m.name()}</Form.Label>
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
								<Form.Label for="bist-input">{m.target()}</Form.Label>
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
												<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
												<SingleSelect.Group>
													{#each $warpTarget as item}
														<SingleSelect.Item option={item}>
															<Icon
																icon={item.icon ? item.icon : 'ph:empty'}
																class={cn(
																	'size-5',
																	item.icon ? 'visible' : 'invisible',
																)}
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
								<Form.Legend>{m.target()}</Form.Legend>
								<Form.Field>
									<Form.Label>{m.internal_object_service()}</Form.Label>
									<ObjectServicesPicker
										bind:selectedInternalObjectService={requestInternalObjectService}
									/>
								</Form.Field>
							</Form.Fieldset>
						{:else if requestWarp.target.case == 'externalObjectService'}
							<Form.Fieldset>
								<Form.Legend>{m.target()}</Form.Legend>
								<Form.Field>
									<Form.Label>{m.endpoint()}</Form.Label>
									<SingleInput.General
										type="text"
										required
										bind:value={requestExternalObjectService.endpoint}
										bind:invalid={invalidEndPoint}
									/>
								</Form.Field>
								<Form.Field>
									<Form.Label>{m.access_key()}</Form.Label>
									<SingleInput.General
										type="text"
										required
										bind:value={requestExternalObjectService.accessKey}
										bind:invalid={invalidAccessKey}
									/>
								</Form.Field>
								<Form.Field>
									<Form.Label>{m.secret_key()}</Form.Label>
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
				</Modal.Model>

				<!-- Step two -->
				<Modal.Model>
					<Form.Root class="max-h-[65vh]">
						<Form.Fieldset>
							<Form.Legend>{m.parameter()}</Form.Legend>
							<!-- warpInputOperation -->
							<Form.Field>
								<Form.Label for="warp-operation">{m.operation()}</Form.Label>
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
												<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
												<SingleSelect.Group>
													{#each $warpInputOperation as item}
														<SingleSelect.Item option={item}>
															<Icon
																icon={item.icon ? item.icon : 'ph:empty'}
																class={cn(
																	'size-5',
																	item.icon ? 'visible' : 'invisible',
																)}
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
								<Form.Label>{m.duration()}</Form.Label>
								<SingleInput.Measurement
									bind:value={warpDuration}
									units={[
										{ value: 1, label: 's' } as SingleInput.UnitType,
										{ value: 60, label: 'm' } as SingleInput.UnitType,
										{ value: 3600, label: 'h' } as SingleInput.UnitType,
										{ value: 86400, label: 'd' } as SingleInput.UnitType,
									]}
								/>
							</Form.Field>
							<!-- ObjectSize -->
							<Form.Field>
								<Form.Label>{m.object_size()}</Form.Label>
								<SingleInput.Measurement
									bind:value={warpObjectSize}
									units={[
										{ value: Math.pow(2, 10 * 0), label: 'B' } as SingleInput.UnitType,
										{ value: Math.pow(2, 10 * 1), label: 'KB' } as SingleInput.UnitType,
										{ value: Math.pow(2, 10 * 2), label: 'MB' } as SingleInput.UnitType,
										{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
										{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType,
										{ value: Math.pow(2, 10 * 5), label: 'PB' } as SingleInput.UnitType,
									]}
								/>
							</Form.Field>
							<!-- ObjectCount -->
							<Form.Field>
								<Form.Label>{m.object_count()}</Form.Label>
								<SingleInput.General
									type="number"
									placeholder="500"
									bind:value={warpObjectCount}
									disabled={warpOperation == Warp_Input_Operation.PUT}
								/>
							</Form.Field>
						</Form.Fieldset>
					</Form.Root>
				</Modal.Model>

				<!-- Step three Overview -->
				<Modal.Model>
					<Form.Root>
						<!-- Step 1 -->
						<Form.Fieldset>
							<Form.Legend>{m.basic()}</Form.Legend>
							<Form.Description>{m.name()}: {request.name}</Form.Description>
							<Form.Description>{m.target()}: {requestWarp.target.case}</Form.Description>
							{#if requestWarp.target.case == 'internalObjectService'}
								<Form.Description>{m.type()}: {requestInternalObjectService.type}</Form.Description>
								<Form.Description>{m.name()}: {requestInternalObjectService.name}</Form.Description>
								<Form.Description
									>{m.endpoint()}: {requestInternalObjectService.endpoint}</Form.Description
								>
							{:else if requestWarp.target.case == 'externalObjectService'}
								<Form.Description
									>{m.endpoint()}: {requestExternalObjectService.endpoint}</Form.Description
								>
								<Form.Description
									>{m.access_key()}: {requestExternalObjectService.accessKey}</Form.Description
								>
								<Form.Description
									>{m.secret_key()}: {requestExternalObjectService.secretKey}</Form.Description
								>
							{/if}
						</Form.Fieldset>
						<!-- Step 2 -->
						<Form.Fieldset>
							{@const duration = formatSecond(Number(warpDuration))}
							{@const objectSize = formatCapacity(Number(warpObjectSize))}
							<Form.Legend>{m.advance()}</Form.Legend>
							<Form.Description>{m.operation()}: {Warp_Input_Operation[warpOperation]}</Form.Description>
							<Form.Description>{m.duration()}: {duration.value} {duration.unit}</Form.Description>
							<Form.Description>{m.object_size()}: {objectSize.value} {objectSize.unit}</Form.Description>
							{#if warpOperation !== Warp_Input_Operation.PUT}
								<Form.Description>{m.object_count()}: {warpObjectCount}</Form.Description>
							{/if}
						</Form.Fieldset>
					</Form.Root>
				</Modal.Model>
			</Modal.Models>
		</Modal.Stepper>

		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}>{m.cancel()}</Modal.Cancel
			>
			<Modal.Controllers>
				<Modal.Back>{m.back()}</Modal.Back>
				<Modal.Next>{m.next()}</Modal.Next>
				<Modal.Confirm
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
							objectCount: warpOperation === Warp_Input_Operation.PUT ? 0n : BigInt(warpObjectCount),
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
									duration: Number.POSITIVE_INFINITY,
								});
								return message;
							},
						});
						reset();
						close();
					}}
				>
					{m.confirm()}
				</Modal.Confirm>
			</Modal.Controllers>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
