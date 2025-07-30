<script lang="ts" module>
	import type { CreateTestResultRequest, ExternalObjectService, InternalObjectService, TestResult, Warp, Warp_Input } from '$gen/api/bist/v1/bist_pb';
	import { BISTService, Warp_Input_Operation } from '$gen/api/bist/v1/bist_pb';
	import ObjectServicesPicker from '$lib/components/bist/utils/object-services-picker.svelte';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as MultipleStepModal from '$lib/components/custom/mutiple-step-modal';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, type Snippet } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable, type Writable } from 'svelte/store';

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
		.filter(key => isNaN(Number(key)))
		.map(key => ({value: Warp_Input_Operation[key as keyof typeof Warp_Input_Operation], label: key}));
	const warpInputOperation: Writable<SingleSelect.OptionType[]> = writable(Options);
</script>

<script lang="ts">
	let {
		testResult,
		data = $bindable(), 
		trigger
	}: { 
		testResult?: TestResult;
		data: Writable<TestResult[]>; 
		trigger?: Snippet<[]>;
	} = $props();

	// Request
    const DEFAULT_WARP_REQUEST = testResult 
		? { target: {value: testResult.kind.value?.target.value, case: testResult.kind.value?.target.case} } as Warp 
		: { target: {value: {}, case: {}}} as Warp;
	const DEFAULT_REQUEST = { kind: {value: DEFAULT_WARP_REQUEST, case: "warp"}, createdBy: "Woody Lin" } as CreateTestResultRequest;
    const DEFAULT_INTERNAL_OBJECT_SERVICE = testResult && testResult.kind.value?.target?.case === 'internalObjectService'
        ? testResult.kind.value.target.value as InternalObjectService
        : {} as InternalObjectService;
    const DEFAULT_DEFAULT_EXTERNAL_OBJECT_SERVICE = testResult && testResult.kind.value?.target?.case === 'externalObjectService'
        ? testResult.kind.value.target.value as ExternalObjectService
        : {} as ExternalObjectService;
	const DEFAULT_WARP_INPUT = testResult && testResult.kind.value?.input
    ? testResult.kind.value.input as Warp_Input
	: { duration: "60s", objectSize: "4MiB", objectCount: "500" } as unknown as Warp_Input; 

	let request: CreateTestResultRequest = $state(DEFAULT_REQUEST);
	let requestWarp: Warp = $state(DEFAULT_WARP_REQUEST);
	let requestInternalObjectService: InternalObjectService = $state(DEFAULT_INTERNAL_OBJECT_SERVICE);
	let requestExternalObjectService: ExternalObjectService = $state(DEFAULT_DEFAULT_EXTERNAL_OBJECT_SERVICE);
	let warpOperation = $state(DEFAULT_WARP_INPUT.operation);
	let warpDuration = $state(DEFAULT_WARP_INPUT.duration);
	let warpObjectSize = $state(DEFAULT_WARP_INPUT.objectSize);
	let warpObjectCount = $state(DEFAULT_WARP_INPUT.objectCount);
	function reset() {
		request = DEFAULT_REQUEST;
		requestWarp = DEFAULT_WARP_REQUEST;
		requestInternalObjectService = DEFAULT_INTERNAL_OBJECT_SERVICE;
		requestExternalObjectService = DEFAULT_DEFAULT_EXTERNAL_OBJECT_SERVICE;
		warpOperation = DEFAULT_WARP_INPUT.operation; 
		warpDuration = DEFAULT_WARP_INPUT.duration; 
		warpObjectSize = DEFAULT_WARP_INPUT.objectSize; 
		warpObjectCount = DEFAULT_WARP_INPUT.objectCount; 
	}

	// Modal state
	const stateController = new DialogStateController(false);

	// grpc
	const transport: Transport = getContext('transport');
	const bistClient = createClient(BISTService, transport);
</script>


<MultipleStepModal.Root bind:open={stateController.state} steps={3}>
	<!-- {@render trigger()} -->
	{#if trigger}
		{@render trigger()}
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

		<MultipleStepModal.Stepper>
			<MultipleStepModal.Steps>
				<MultipleStepModal.Step icon="ph:number-one" />
				<MultipleStepModal.Step icon="ph:number-two" />
				<MultipleStepModal.Step icon="ph:number-three" />
			</MultipleStepModal.Steps>
			<!-- <MultipleStepModal.Header class="flex m-4 items-center justify-center text-xl font-bold"> -->
			<MultipleStepModal.Header class="flex mt-6 mb-6 justify-center text-xl font-bold">
				BIST
			</MultipleStepModal.Header>	
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
									bind:value={request.name}
								/>
							</Form.Field>
							<!-- Choose Target -->
							<Form.Field>
								<Form.Label for="bist-input">Target</Form.Label>
								<SingleSelect.Root options={warpTarget} required bind:value={requestWarp.target.case}>
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
									<ObjectServicesPicker bind:selectedInternalObjectService={requestInternalObjectService} />
								</Form.Field>
							</Form.Fieldset>
						{:else if requestWarp.target.case == 'externalObjectService'}
							<Form.Fieldset>
								<Form.Legend>Target</Form.Legend>
								<Form.Field>
									<Form.Label>Endpoint</Form.Label>
									<SingleInput.General type="text" required bind:value={requestExternalObjectService.endpoint}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Access Key</Form.Label>
									<SingleInput.General type="text" required bind:value={requestExternalObjectService.accessKey}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Secret Key</Form.Label>
									<SingleInput.General type="text" required bind:value={requestExternalObjectService.secretKey}/>
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
								<SingleSelect.Root options={warpInputOperation} bind:value={warpOperation}>
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
								<SingleInput.General type="text" placeholder="32" bind:value={warpDuration}/>
							</Form.Field>
							<!-- ObjectSize -->
							<Form.Field>
								<Form.Label>Object Size</Form.Label>
								<SingleInput.General type="text" placeholder="100" bind:value={warpObjectSize}/>
							</Form.Field>
							<!-- ObjectCount -->
							<Form.Field>
								<Form.Label>Object Count</Form.Label>
								<SingleInput.General type="text" placeholder="4k" bind:value={warpObjectCount}/>
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
								<Form.Description>type: {requestInternalObjectService.type}</Form.Description>
								<Form.Description>name: {requestInternalObjectService.name}</Form.Description>
								<Form.Description>endpoint: {requestInternalObjectService.endpoint}</Form.Description>
							{:else if requestWarp.target.case == 'externalObjectService'}
								<Form.Description>Endpoint: {requestExternalObjectService.endpoint}</Form.Description>
								<Form.Description>Access Key: {requestExternalObjectService.accessKey}</Form.Description>
								<Form.Description>Secret Key: {requestExternalObjectService.secretKey}</Form.Description>
							{/if}
						</Form.Fieldset>
						<!-- Step 2 -->
						<Form.Fieldset>
							<Form.Legend>Step 2</Form.Legend>
							<Form.Description>Operation: {Warp_Input_Operation[warpOperation]}</Form.Description>
							<Form.Description>Duration: {warpDuration}</Form.Description>
							<Form.Description>Object Size: {warpObjectSize}</Form.Description>
							<Form.Description>Object Count: {warpObjectCount}</Form.Description>
						</Form.Fieldset>
					</Form.Root>
				</MultipleStepModal.Model>
			</MultipleStepModal.Models>
		</MultipleStepModal.Stepper>
		
		<MultipleStepModal.Footer>
			<MultipleStepModal.Cancel onclick={() => { reset(); }}>Cancel</MultipleStepModal.Cancel>
			<MultipleStepModal.Confirm
				onclick={() => {
					// prepare request
                    if (requestWarp.target.case == 'internalObjectService') {
                        requestWarp.target.value = requestInternalObjectService;
                    } else if (requestWarp.target.case == 'externalObjectService') {
                        requestWarp.target.value = requestExternalObjectService;
                    }
					requestWarp.input = {
						operation: warpOperation,
						duration: warpDuration,
						objectSize: warpObjectSize,
						objectCount: warpObjectCount
					} as Warp_Input;
					request.kind.value = requestWarp;
					// request
					bistClient
						.createTestResult(request)
						.then((r) => {
							toast.success(`Create ${r.name}`);
							bistClient
								.listTestResults({})
								.then((r) => {
									data.set(r.testResults.filter((result) => result.kind.case === 'warp' ));
								});
						})
						.catch((e) => {
							console.log(e.toString());
							toast.error(`Fail to create test: ${e.toString()}`);
						})
						.finally(() => {
							reset();
							stateController.close();
						});
					stateController.close();
				}}
			>
				Confirm
			</MultipleStepModal.Confirm>
			<MultipleStepModal.Controllers>
				<MultipleStepModal.Back>Back</MultipleStepModal.Back>
				<MultipleStepModal.Next>Next</MultipleStepModal.Next>
			</MultipleStepModal.Controllers>
		</MultipleStepModal.Footer>
	</MultipleStepModal.Content>
</MultipleStepModal.Root>

