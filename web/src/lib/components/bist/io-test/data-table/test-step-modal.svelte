<script lang="ts" module>
	import type { CephBlockDevice, CreateTestResultRequest, FIO, FIO_Input, NetworkFileSystem, TestResult } from '$gen/api/bist/v1/bist_pb';
	import { BISTService, FIO_Input_AccessMode } from '$gen/api/bist/v1/bist_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as MultipleStepModal from '$lib/components/custom/mutiple-step-modal';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import CephPicker from '$lib/components/management-storage/utils/ceph-picker.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, type Snippet } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable, type Writable } from 'svelte/store';
	import { formatByte, formatSecond } from '$lib/formatter';

	// FIO Target
	const fioTarget: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'cephBlockDevice',
			label: 'Ceph Block Device',
		},
		{
			value: 'networkFileSystem',
			label: 'Network File System',
		},
	]);

	// FIO AccessMode
	const Options: SingleSelect.OptionType[] = Object.keys(FIO_Input_AccessMode)
		.filter(key => isNaN(Number(key)))
		.map(key => ({value: FIO_Input_AccessMode[key as keyof typeof FIO_Input_AccessMode], label: key}));
	const fioInputeAccessMode: Writable<SingleSelect.OptionType[]> = writable(Options);
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
    const DEFAULT_FIO_REQUEST = testResult 
		? { target: {value: testResult.kind.value?.target.value, case: testResult.kind.value?.target.case} } as FIO 
		: { target: {value: {}, case: {}}} as FIO;
	const DEFAULT_REQUEST = { kind: {value: DEFAULT_FIO_REQUEST, case: "fio"}, createdBy: "Woody Lin" } as CreateTestResultRequest;
    const DEFAULT_CEPH_BLOCK_DEVICE = testResult && testResult.kind.value?.target?.case === 'cephBlockDevice'
        ? testResult.kind.value.target.value as CephBlockDevice
        : {} as CephBlockDevice;
    const DEFAULT_NETWORK_FILE_SYSTEM = testResult && testResult.kind.value?.target?.case === 'networkFileSystem'
        ? testResult.kind.value.target.value as NetworkFileSystem
        : {} as NetworkFileSystem;
	const DEFAULT_FIO_INPUT = testResult && testResult.kind.value?.input
    ? testResult.kind.value.input as FIO_Input
	: { jobCount: 32, runTimeSeconds: 60, blockSizeBytes: 4096, fileSizeBytes: 1024 * 1024 * 1024, ioDepth: 1 } as unknown as FIO_Input; 
    let selectedScope = $state(
        testResult && testResult.kind.value?.target?.case === 'cephBlockDevice'
            ? testResult.kind.value.target.value?.scopeUuid ?? ''
            : ''
    );
    let selectedFacility = $state(
        testResult && testResult.kind.value?.target?.case === 'cephBlockDevice'
            ? testResult.kind.value.target.value?.facilityName ?? ''
            : ''
    );
	let request: CreateTestResultRequest = $state(DEFAULT_REQUEST);
	let requestFio: FIO = $state(DEFAULT_FIO_REQUEST);
	let requestCephBlockDevice: CephBlockDevice = $state(DEFAULT_CEPH_BLOCK_DEVICE);
	let requestNetworkFileSystem: NetworkFileSystem = $state(DEFAULT_NETWORK_FILE_SYSTEM);
	let fioAccessMode = $state(DEFAULT_FIO_INPUT.accessMode);
	let fioJobCount = $state(DEFAULT_FIO_INPUT.jobCount);
	let fioRunTime = $state(DEFAULT_FIO_INPUT.runTimeSeconds);
	let fioBlockSize = $state(DEFAULT_FIO_INPUT.blockSizeBytes);
	let fioFileSize = $state(DEFAULT_FIO_INPUT.fileSizeBytes);
	let fioIoDepth = $state(DEFAULT_FIO_INPUT.ioDepth);
	function reset() {
		request = DEFAULT_REQUEST;
		requestFio = DEFAULT_FIO_REQUEST;
		requestCephBlockDevice = DEFAULT_CEPH_BLOCK_DEVICE;
		requestNetworkFileSystem = DEFAULT_NETWORK_FILE_SYSTEM;
		fioAccessMode = DEFAULT_FIO_INPUT.accessMode; 
		fioJobCount = DEFAULT_FIO_INPUT.jobCount; 
		fioRunTime = DEFAULT_FIO_INPUT.runTimeSeconds; 
		fioBlockSize = DEFAULT_FIO_INPUT.blockSizeBytes; 
		fioFileSize = DEFAULT_FIO_INPUT.fileSizeBytes; 
		fioIoDepth = DEFAULT_FIO_INPUT.ioDepth; 
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
								<SingleSelect.Root options={fioTarget} required bind:value={requestFio.target.case}>
									<SingleSelect.Trigger />
									<SingleSelect.Content>
										<SingleSelect.Options>
											<SingleSelect.Input />
											<SingleSelect.List>
												<SingleSelect.Empty>No results found.</SingleSelect.Empty>
												<SingleSelect.Group>
													{#each $fioTarget as item}
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
						{#if requestFio.target.case == 'cephBlockDevice'}
							<Form.Fieldset>
								<Form.Legend>Target</Form.Legend>
								<Form.Field>
									<Form.Label>Ceph</Form.Label>
									<CephPicker bind:selectedScope bind:selectedFacility />
								</Form.Field>
							</Form.Fieldset>
						{:else if requestFio.target.case == 'networkFileSystem'}
							<Form.Fieldset>
								<Form.Legend>Target</Form.Legend>
								<Form.Field>
									<Form.Label>Endpoint</Form.Label>
									<SingleInput.General type="text" required bind:value={requestNetworkFileSystem.endpoint}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Path</Form.Label>
									<SingleInput.General type="text" required bind:value={requestNetworkFileSystem.path}/>
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
							<!-- fioInputeAccessMode -->
							<Form.Field>
								<Form.Label for="fio-access-mode">Access Mode</Form.Label>
								<SingleSelect.Root options={fioInputeAccessMode} bind:value={fioAccessMode}>
									<SingleSelect.Trigger />
									<SingleSelect.Content>
										<SingleSelect.Options>
											<SingleSelect.Input />
											<SingleSelect.List>
												<SingleSelect.Empty>No results found.</SingleSelect.Empty>
												<SingleSelect.Group>
													{#each $fioInputeAccessMode as item}
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
							<!-- jobCount -->
							<Form.Field>
								<Form.Label>Job Count</Form.Label>
								<SingleInput.General type="number" placeholder="32" bind:value={fioJobCount}/>
							</Form.Field>
							<!-- runTime -->
							<Form.Field>
								<Form.Label>Run Time</Form.Label>
								<SingleInput.Measurement
									bind:value={fioRunTime}
									units={[
										{ value: 1, label: 's' } as SingleInput.UnitType,
										{ value: 60, label: 'm' } as SingleInput.UnitType,
										{ value: 3600, label: 'h' } as SingleInput.UnitType,
										{ value: 86400, label: 'd' } as SingleInput.UnitType
									]}
								/>
							</Form.Field>
							<!-- blockSize -->
							<Form.Field>
								<Form.Label>Block Size</Form.Label>
								<SingleInput.Measurement
									bind:value={fioBlockSize}
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
							<!-- fileSize -->
							<Form.Field>
								<Form.Label>File Size</Form.Label>
								<SingleInput.Measurement
									bind:value={fioFileSize}
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
							<!-- ioDepth -->
							<Form.Field>
								<Form.Label>I/O Depth</Form.Label>
								<SingleInput.General type="number" placeholder="1" bind:value={fioIoDepth}/>
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
							<Form.Description>Target: {requestFio.target.case}</Form.Description>
							{#if requestFio.target.case == 'cephBlockDevice'}
								<Form.Description>Scope UUID: {selectedScope}</Form.Description>
								<Form.Description>Facility Name: {selectedFacility}</Form.Description>
							{:else if requestFio.target.case == 'networkFileSystem'}
								<Form.Description>type: {requestNetworkFileSystem.endpoint}</Form.Description>
								<Form.Description>name: {requestNetworkFileSystem.path}</Form.Description>
							{/if}
						</Form.Fieldset>
						<!-- Step 2 -->
						<Form.Fieldset>
							{@const runTime = formatSecond(Number(fioRunTime))}
							{@const blockSize = formatByte(Number(fioBlockSize))}
							{@const fileSize = formatByte(Number(fioFileSize))}
							<Form.Legend>Step 2</Form.Legend>
							<Form.Description>Access Mode: {FIO_Input_AccessMode[fioAccessMode]}</Form.Description>
							<Form.Description>Job Count: {fioJobCount}</Form.Description>
							<Form.Description>Run Time: {runTime.value} {runTime.unit}</Form.Description>
							<Form.Description>Block Size: {blockSize.value} {blockSize.unit}</Form.Description>
							<Form.Description>File Size: {fileSize.value} {fileSize.unit}</Form.Description>
							<Form.Description>I/O Depth: {fioIoDepth}</Form.Description>
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
                    if (requestFio.target.case == 'cephBlockDevice') {
						requestCephBlockDevice.scopeUuid = selectedScope;
						requestCephBlockDevice.facilityName = selectedFacility;
                        requestFio.target.value = requestCephBlockDevice;
                    } else if (requestFio.target.case == 'networkFileSystem') {
                        requestFio.target.value = requestNetworkFileSystem;
                    }
					requestFio.input = {
						accessMode: fioAccessMode,
						jobCount: BigInt(fioJobCount),
						runTimeSeconds: BigInt(fioRunTime),
						blockSizeBytes: BigInt(fioBlockSize),
						fileSizeBytes: BigInt(fioFileSize),
						ioDepth: BigInt(fioIoDepth)
					} as FIO_Input;
					request.kind.value = requestFio;
					// request
					bistClient
						.createTestResult(request)
						.then((r) => {
							toast.success(`Create ${r.name}`);
							bistClient
								.listTestResults({})
								.then((r) => {
									data.set(r.testResults.filter((result) => result.kind.case === 'fio' ));
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

