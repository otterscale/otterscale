<script lang="ts" module>
	import type {
		CephBlockDevice,
		CreateTestResultRequest,
		FIO,
		FIO_Input,
		NetworkFileSystem,
		TestResult,
	} from '$lib/api/bist/v1/bist_pb';
	import { BISTService, FIO_Input_AccessMode } from '$lib/api/bist/v1/bist_pb';
	import { Essential_Type, EssentialService } from '$lib/api/essential/v1/essential_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as Loading from '$lib/components/custom/loading';
	import { MultipleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { formatCapacity, formatSecond } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable, type Writable } from 'svelte/store';

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
		.filter((key) => isNaN(Number(key)))
		.map((key) => ({
			value: FIO_Input_AccessMode[key as keyof typeof FIO_Input_AccessMode],
			label: key,
		}));
	const fioInputeAccessMode: Writable<SingleSelect.OptionType[]> = writable(Options);
</script>

<script lang="ts">
	let {
		testResult,
	}: {
		testResult?: TestResult;
	} = $props();

	// Request
	const DEFAULT_FIO_REQUEST = testResult
		? ({
				target: {
					value: testResult.kind.value?.target.value,
					case: testResult.kind.value?.target.case,
				},
			} as FIO)
		: ({ target: { value: {}, case: {} } } as FIO);
	const DEFAULT_REQUEST = {
		kind: { value: DEFAULT_FIO_REQUEST, case: 'fio' },
		createdBy: 'Woody Lin',
	} as CreateTestResultRequest;
	const DEFAULT_CEPH_BLOCK_DEVICE =
		testResult && testResult.kind.value?.target?.case === 'cephBlockDevice'
			? (testResult.kind.value.target.value as CephBlockDevice)
			: ({} as CephBlockDevice);
	const DEFAULT_NETWORK_FILE_SYSTEM =
		testResult && testResult.kind.value?.target?.case === 'networkFileSystem'
			? (testResult.kind.value.target.value as NetworkFileSystem)
			: ({} as NetworkFileSystem);
	const DEFAULT_FIO_INPUT =
		testResult && testResult.kind.value?.input
			? (testResult.kind.value.input as FIO_Input)
			: ({
					jobCount: 32,
					runTimeSeconds: 60,
					blockSizeBytes: 4096,
					fileSizeBytes: 1024 * 1024 * 1024,
					ioDepth: 1,
				} as unknown as FIO_Input);
	let selectedScope = $state(
		testResult && testResult.kind.value?.target?.case === 'cephBlockDevice'
			? (testResult.kind.value.target.value?.scopeUuid ?? '')
			: '',
	);
	let selectedFacility = $state(
		testResult && testResult.kind.value?.target?.case === 'cephBlockDevice'
			? (testResult.kind.value.target.value?.facilityName ?? '')
			: '',
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
	let open = $state(false);
	function close() {
		open = false;
	}

	// grpc
	const transport: Transport = getContext('transport');
	const bistClient = createClient(BISTService, transport);
	const essentialClient = createClient(EssentialService, transport);

	const reloadManager: ReloadManager = getContext('reloadManager');

	const cephOptions = writable<SingleSelect.OptionType[]>([]);
	let isCephsLoading = $state(true);
	async function fetchCephOptions() {
		try {
			const response = await essentialClient.listEssentials({ type: Essential_Type.CEPH });
			cephOptions.set(
				response.essentials.map(
					(essential) =>
						({
							value: { scopeUuid: essential.scopeUuid, facilityName: essential.name },
							label: `${essential.scopeName}-${essential.name}`,
							icon: 'ph:cube',
						}) as SingleSelect.OptionType,
				),
			);
			isCephsLoading = false;
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchCephOptions();
			isMounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<Modal.Root bind:open steps={3}>
	<!-- {@render trigger()} -->
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
		<Modal.Header class="mt-6 mb-6 flex justify-center text-xl font-bold">
			{m.built_in_self_test()}
		</Modal.Header>
		<Modal.Stepper>
			<Modal.Steps>
				<Modal.Step icon="ph:number-one" />
				<Modal.Step icon="ph:number-two" />
				<Modal.Step icon="ph:number-three" />
			</Modal.Steps>
			<Modal.Models>
				<!-- Step One -->
				<Modal.Model>
					<Form.Root class="max-h-[65vh]">
						<Form.Fieldset>
							<!-- Name -->
							<Form.Field>
								<Form.Label for="bist-name">{m.name()}</Form.Label>
								<SingleInput.General type="text" id="name" bind:value={request.name} />
							</Form.Field>
							<!-- Choose Target -->
							<Form.Field>
								<Form.Label for="bist-input">{m.target()}</Form.Label>
								<SingleSelect.Root options={fioTarget} required bind:value={requestFio.target.case}>
									<SingleSelect.Trigger />
									<SingleSelect.Content>
										<SingleSelect.Options>
											<SingleSelect.Input />
											<SingleSelect.List>
												<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
												<SingleSelect.Group>
													{#each $fioTarget as item}
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
						{#if requestFio.target.case == 'cephBlockDevice'}
							<Form.Fieldset>
								<Form.Legend>{m.target()}</Form.Legend>
								<Form.Field>
									<Form.Label>{m.ceph()}</Form.Label>
									{#if isCephsLoading}
										<Loading.Selection />
									{:else}
										<SingleSelect.Root options={cephOptions}>
											<SingleSelect.Trigger />
											<SingleSelect.Content>
												<SingleSelect.Options>
													<SingleSelect.Input />
													<SingleSelect.List>
														<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
														<SingleSelect.Group>
															{#each $cephOptions as option}
																<SingleSelect.Item
																	{option}
																	onclick={() => {
																		selectedScope = option.value.scopeUuid;
																		selectedFacility = option.value.facilityName;
																	}}
																>
																	<Icon
																		icon={option.icon ? option.icon : 'ph:empty'}
																		class={cn(
																			'size-5',
																			option.icon ? 'visibale' : 'invisible',
																		)}
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
						{:else if requestFio.target.case == 'networkFileSystem'}
							<Form.Fieldset>
								<Form.Legend>{m.target()}</Form.Legend>
								<Form.Field>
									<Form.Label>{m.endpoint()}</Form.Label>
									<SingleInput.General
										type="text"
										required
										bind:value={requestNetworkFileSystem.endpoint}
									/>
								</Form.Field>
								<Form.Field>
									<Form.Label>{m.path()}</Form.Label>
									<SingleInput.General
										type="text"
										required
										bind:value={requestNetworkFileSystem.path}
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
							<!-- fioInputeAccessMode -->
							<Form.Field>
								<Form.Label for="fio-access-mode">{m.access_mode()}</Form.Label>
								<SingleSelect.Root options={fioInputeAccessMode} required bind:value={fioAccessMode}>
									<SingleSelect.Trigger />
									<SingleSelect.Content>
										<SingleSelect.Options>
											<SingleSelect.Input />
											<SingleSelect.List>
												<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
												<SingleSelect.Group>
													{#each $fioInputeAccessMode as item}
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
							<!-- jobCount -->
							<Form.Field>
								<Form.Label>{m.job_count()}</Form.Label>
								<SingleInput.General type="number" placeholder="32" bind:value={fioJobCount} />
							</Form.Field>
							<!-- runTime -->
							<Form.Field>
								<Form.Label>{m.run_time()}</Form.Label>
								<SingleInput.Measurement
									bind:value={fioRunTime}
									units={[
										{ value: 1, label: 's' } as SingleInput.UnitType,
										{ value: 60, label: 'm' } as SingleInput.UnitType,
										{ value: 3600, label: 'h' } as SingleInput.UnitType,
										{ value: 86400, label: 'd' } as SingleInput.UnitType,
									]}
								/>
							</Form.Field>
							<!-- blockSize -->
							<Form.Field>
								<Form.Label>{m.block_size()}</Form.Label>
								<SingleInput.Measurement
									bind:value={fioBlockSize}
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
							<!-- fileSize -->
							<Form.Field>
								<Form.Label>{m.file_size()}</Form.Label>
								<SingleInput.Measurement
									bind:value={fioFileSize}
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
							<!-- ioDepth -->
							<Form.Field>
								<Form.Label>{m.io_depth()}</Form.Label>
								<SingleInput.General type="number" placeholder="1" bind:value={fioIoDepth} />
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
							<Form.Description>{m.target()}: {requestFio.target.case}</Form.Description>
							{#if requestFio.target.case == 'cephBlockDevice'}
								<Form.Description>{m.scope()}: {selectedScope}</Form.Description>
								<Form.Description>{m.facility()}: {selectedFacility}</Form.Description>
							{:else if requestFio.target.case == 'networkFileSystem'}
								<Form.Description>{m.type()}: {requestNetworkFileSystem.endpoint}</Form.Description>
								<Form.Description>{m.name()}: {requestNetworkFileSystem.path}</Form.Description>
							{/if}
						</Form.Fieldset>
						<!-- Step 2 -->
						<Form.Fieldset>
							{@const runTime = formatSecond(Number(fioRunTime))}
							{@const blockSize = formatCapacity(Number(fioBlockSize))}
							{@const fileSize = formatCapacity(Number(fioFileSize))}
							<Form.Legend>{m.advance()}</Form.Legend>
							<Form.Description>{m.access_mode()}: {FIO_Input_AccessMode[fioAccessMode]}</Form.Description
							>
							<Form.Description>{m.job_count()}: {fioJobCount}</Form.Description>
							<Form.Description>{m.run_time()}: {runTime.value} {runTime.unit}</Form.Description>
							<Form.Description>{m.block_size()}: {blockSize.value} {blockSize.unit}</Form.Description>
							<Form.Description>{m.file_size()}: {fileSize.value} {fileSize.unit}</Form.Description>
							<Form.Description>{m.io_depth()}: {fioIoDepth}</Form.Description>
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
							ioDepth: BigInt(fioIoDepth),
						} as FIO_Input;
						request.kind.value = requestFio;
						// request
						toast.promise(() => bistClient.createTestResult(request), {
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
