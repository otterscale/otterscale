<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/state';
	import type {
		CephBlockDevice,
		CreateTestResultRequest,
		FIO,
		FIO_Input,
		NetworkFileSystem,
		TestResult
	} from '$lib/api/configuration/v1/configuration_pb';
	import {
		ConfigurationService,
		FIO_Input_AccessMode
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { MultipleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { formatCapacity, formatSecond } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	// FIO Target
	const fioTarget: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'cephBlockDevice',
			label: 'Ceph Block Device',
			icon: 'ph:hard-drives'
		},
		{
			value: 'networkFileSystem',
			label: 'Network File System',
			icon: 'ph:network'
		}
	]);

	// FIO AccessMode
	const Options: SingleSelect.OptionType[] = Object.keys(FIO_Input_AccessMode)
		.filter((key) => isNaN(Number(key)))
		.map((key) => ({
			value: FIO_Input_AccessMode[key as keyof typeof FIO_Input_AccessMode],
			label: key,
			icon: 'ph:cube'
		}));
	const fioInputAccessMode: Writable<SingleSelect.OptionType[]> = writable(Options);
</script>

<script lang="ts">
	let {
		testResult,
		scope,
		reloadManager,
		closeActions
	}: {
		testResult?: TestResult;
		scope: string;
		reloadManager: ReloadManager;
		closeActions?: () => void;
	} = $props();

	let request: CreateTestResultRequest = $state({} as CreateTestResultRequest);
	let requestFio: FIO = $state({} as FIO);
	let requestCephBlockDevice: CephBlockDevice = $state({} as CephBlockDevice);
	let requestNetworkFileSystem: NetworkFileSystem = $state({} as NetworkFileSystem);
	let fioInput = $state({} as FIO_Input);

	function init() {
		request = {
			kind: { value: {} as FIO, case: 'fio' },
			createdBy: page.data.user?.username ?? ''
		} as CreateTestResultRequest;
		requestFio = testResult
			? ({
					target: {
						value: testResult.kind.value?.target.value,
						case: testResult.kind.value?.target.case
					}
				} as FIO)
			: ({ target: { value: {}, case: {} } } as FIO);
		requestCephBlockDevice =
			testResult && testResult.kind.value?.target?.case === 'cephBlockDevice'
				? (testResult.kind.value.target.value as CephBlockDevice)
				: ({} as CephBlockDevice);
		requestNetworkFileSystem =
			testResult && testResult.kind.value?.target?.case === 'networkFileSystem'
				? (testResult.kind.value.target.value as NetworkFileSystem)
				: ({} as NetworkFileSystem);
		fioInput =
			testResult && testResult.kind.value?.input
				? (testResult.kind.value.input as FIO_Input)
				: ({
						jobCount: 32,
						runTimeSeconds: 60,
						blockSizeBytes: 4096,
						fileSizeBytes: 1024 * 1024 * 1024,
						ioDepth: 1
					} as unknown as FIO_Input);
	}

	// Modal state
	let open = $state(false);
	function close() {
		open = false;
	}

	// grpc
	const transport: Transport = getContext('transport');
	const configClient = createClient(ConfigurationService, transport);
</script>

<Modal.Root
	bind:open
	steps={3}
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
	onOpenChangeComplete={(isOpen) => {
		if (closeActions && !isOpen) {
			closeActions();
		}
	}}
>
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
													{#each $fioTarget as item (item.value)}
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
						{#if requestFio.target.case == 'networkFileSystem'}
							<Form.Fieldset>
								<Form.Legend>{m.target()}</Form.Legend>
								<Form.Field>
									<Form.Label>{m.host()}</Form.Label>
									<SingleInput.General
										type="text"
										required
										bind:value={requestNetworkFileSystem.host}
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
							<!-- fioInputAccessMode -->
							<Form.Field>
								<Form.Label for="fio-access-mode">{m.access_mode()}</Form.Label>
								<SingleSelect.Root
									options={fioInputAccessMode}
									required
									bind:value={fioInput.accessMode}
								>
									<SingleSelect.Trigger />
									<SingleSelect.Content>
										<SingleSelect.Options>
											<SingleSelect.Input />
											<SingleSelect.List>
												<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
												<SingleSelect.Group>
													{#each $fioInputAccessMode as item (item.value)}
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
								<Form.Label>{m.job_count()}</Form.Label>
								<SingleInput.General
									type="number"
									placeholder="32"
									bind:value={fioInput.jobCount}
								/>
							</Form.Field>
							<!-- runTime -->
							<Form.Field>
								<Form.Label>{m.run_time()}</Form.Label>
								<SingleInput.Measurement
									bind:value={fioInput.runTimeSeconds}
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
								<Form.Label>{m.block_size()}</Form.Label>
								<SingleInput.Measurement
									bind:value={fioInput.blockSizeBytes}
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
								<Form.Label>{m.file_size()}</Form.Label>
								<SingleInput.Measurement
									bind:value={fioInput.fileSizeBytes}
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
								<Form.Label>{m.io_depth()}</Form.Label>
								<SingleInput.General type="number" placeholder="1" bind:value={fioInput.ioDepth} />
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
								<Form.Description>{m.scope()}: {scope}</Form.Description>
							{:else if requestFio.target.case == 'networkFileSystem'}
								<Form.Description>{m.type()}: {requestNetworkFileSystem.host}</Form.Description>
								<Form.Description>{m.name()}: {requestNetworkFileSystem.path}</Form.Description>
							{/if}
						</Form.Fieldset>
						<!-- Step 2 -->
						<Form.Fieldset>
							{@const runTime = formatSecond(Number(fioInput.runTimeSeconds))}
							{@const blockSize = formatCapacity(Number(fioInput.blockSizeBytes))}
							{@const fileSize = formatCapacity(Number(fioInput.fileSizeBytes))}
							<Form.Legend>{m.advance()}</Form.Legend>
							<Form.Description
								>{m.access_mode()}: {FIO_Input_AccessMode[fioInput.accessMode]}</Form.Description
							>
							<Form.Description>{m.job_count()}: {fioInput.jobCount}</Form.Description>
							<Form.Description>{m.run_time()}: {runTime.value} {runTime.unit}</Form.Description>
							<Form.Description
								>{m.block_size()}: {blockSize.value} {blockSize.unit}</Form.Description
							>
							<Form.Description>{m.file_size()}: {fileSize.value} {fileSize.unit}</Form.Description>
							<Form.Description>{m.io_depth()}: {fioInput.ioDepth}</Form.Description>
						</Form.Fieldset>
					</Form.Root>
				</Modal.Model>
			</Modal.Models>
		</Modal.Stepper>

		<Modal.Footer>
			<Modal.Cancel>{m.cancel()}</Modal.Cancel>
			<Modal.Controllers>
				<Modal.Back>{m.back()}</Modal.Back>
				<Modal.Next>{m.next()}</Modal.Next>
				<Modal.Confirm
					onclick={() => {
						// prepare request
						if (requestFio.target.case == 'cephBlockDevice') {
							requestCephBlockDevice.scope = scope;
							requestFio.target.value = requestCephBlockDevice;
						} else if (requestFio.target.case == 'networkFileSystem') {
							requestFio.target.value = requestNetworkFileSystem;
						}
						requestFio.input = {
							accessMode: fioInput.accessMode,
							jobCount: BigInt(fioInput.jobCount),
							runTimeSeconds: BigInt(fioInput.runTimeSeconds),
							blockSizeBytes: BigInt(fioInput.blockSizeBytes),
							fileSizeBytes: BigInt(fioInput.fileSizeBytes),
							ioDepth: BigInt(fioInput.ioDepth)
						} as FIO_Input;
						request.kind.value = requestFio;
						// request
						toast.promise(() => configClient.createTestResult(request), {
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
									closeButton: true
								});
								return message;
							}
						});
						close();
					}}
				>
					{m.confirm()}
				</Modal.Confirm>
			</Modal.Controllers>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
