<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Single as SingleInput, Multiple as MultipleInput } from '$lib/components/custom/input';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Form from '$lib/components/custom/form';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { writable, type Writable } from 'svelte/store';
	import * as MultipleStepModal from './mutiple-step-modal';
	import { TestResult_Type, TestResult_FIO_AccessMode, TestResult_Warp_Operation } from '$gen/api/bist/v1/bist_pb'
	import type { CreateTestResultRequest, TestResult_FIO, TestResult_Warp } from '$gen/api/bist/v1/bist_pb'

	const testResultType: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: TestResult_Type.BLOCK,
			label: TestResult_Type[TestResult_Type.BLOCK],
			// icon: 'ph:upload'
		},
		{
			value: TestResult_Type.NFS,
			label: TestResult_Type[TestResult_Type.NFS],
			// icon: 'ph:upload'
		},
		{
			value: TestResult_Type.S3,
			label: TestResult_Type[TestResult_Type.S3],
			// icon: 'ph:download'
		},
	]);

	const fioOptions: SingleSelect.OptionType[] = Object.keys(TestResult_FIO_AccessMode)
		.filter(key => isNaN(Number(key)))
		.map(key => ({value: TestResult_FIO_AccessMode[key as keyof typeof TestResult_FIO_AccessMode], label: key}));
	const testResultFIOMode: Writable<SingleSelect.OptionType[]> = writable(fioOptions);
	const warpOptions: SingleSelect.OptionType[] = Object.keys(TestResult_Warp_Operation)
		.filter(key => isNaN(Number(key)))
		.map(key => ({value: TestResult_Warp_Operation[key as keyof typeof TestResult_Warp_Operation], label: key}));
	const testResultWarpOperation: Writable<SingleSelect.OptionType[]> = writable(warpOptions);

	const testResultInput: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'fio',
			label: 'FIO',
		},
		{
			value: 'warp',
			label: 'WARP',
		},
	]);

	// const testResultFIOMode: Writable<SingleSelect.OptionType[]> = writable([
	// 	{value: TestResult_FIO_AccessMode.READ, label: 'FIO'},
	// ]);
	
	const DEFAULT_REQUEST = { input: {value: {}, case: {}} } as CreateTestResultRequest;
	const DEFAULT_FIO_REQUEST = { } as TestResult_FIO;
	const DEFAULT_WARP_REQUEST = { } as TestResult_Warp;
	let request: CreateTestResultRequest = $state(DEFAULT_REQUEST);
	let requestFIO: TestResult_FIO = $state(DEFAULT_FIO_REQUEST);
	let requestWarp: TestResult_Warp = $state(DEFAULT_WARP_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
		requestFIO = DEFAULT_FIO_REQUEST;
		requestWarp = DEFAULT_WARP_REQUEST;
	}

	const stateController = new DialogStateController(false);

	// STEP
	let step = $state(1);
	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<MultipleStepModal.Root bind:open steps={3}>
	<div class="flex justify-end">
		<MultipleStepModal.Trigger class={cn(buttonVariants({ variant: 'default', size: 'sm' }))}>
			<div class="flex items-center gap-1">
				<Icon icon="ph:plus" />
				Create
			</div>
		</MultipleStepModal.Trigger>
	</div>
	<MultipleStepModal.Content>

		<MultipleStepModal.Stepper>
			<MultipleStepModal.Steps>
				<MultipleStepModal.Step text="Step 1" icon="ph:number-one" />
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
					<Form.Root class="max-h-[50vh]">
					<Form.Fieldset>
						<Form.Legend>Step 1</Form.Legend>
						<Form.Field>
							<Form.Label for="bist-name">Name</Form.Label>
							<SingleInput.General
								type="text"
								id="name"
								bind:value={request.name}
							/>
						</Form.Field>
						<Form.Field>
							<Form.Label for="bist-type">Type</Form.Label>
							<SingleSelect.Root options={testResultType} required bind:value={request.type}>
								<SingleSelect.Trigger />
								<SingleSelect.Content>
									<SingleSelect.Options>
										<SingleSelect.Input />
										<SingleSelect.List>
											<SingleSelect.Empty>No results found.</SingleSelect.Empty>
											<SingleSelect.Group>
												{#each $testResultType as item}
													<SingleSelect.Item option={item}>
														<Icon
															icon={item.icon ? item.icon : 'ph:empty'}
															class={cn('size-5', item.icon ? 'visibale' : 'invisible')}
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
						<Form.Field>
							<Form.Label for="bist-input">Input</Form.Label>
							<SingleSelect.Root options={testResultInput} required bind:value={request.input.case}>
								<SingleSelect.Trigger />
								<SingleSelect.Content>
									<SingleSelect.Options>
										<SingleSelect.Input />
										<SingleSelect.List>
											<SingleSelect.Empty>No results found.</SingleSelect.Empty>
											<SingleSelect.Group>
												{#each $testResultInput as item}
													<SingleSelect.Item option={item}>
														<Icon
															icon={item.icon ? item.icon : 'ph:empty'}
															class={cn('size-5', item.icon ? 'visibale' : 'invisible')}
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
					</Form.Root>
				</MultipleStepModal.Model>

				<!-- Step two -->
				<MultipleStepModal.Model>
					<Form.Root class="max-h-[50vh]">
						{#if request.input.case == 'fio'}
							<Form.Fieldset>
								<Form.Field>
									<Form.Label for="fio-access-mode">Access Mode</Form.Label>
									<SingleSelect.Root options={testResultFIOMode} bind:value={requestFIO.accessMode}>
										<SingleSelect.Trigger />
										<SingleSelect.Content>
											<SingleSelect.Options>
												<SingleSelect.Input />
												<SingleSelect.List>
													<SingleSelect.Empty>No results found.</SingleSelect.Empty>
													<SingleSelect.Group>
														{#each $testResultFIOMode as item}
															<SingleSelect.Item option={item}>
																<Icon
																	icon={item.icon ? item.icon : 'ph:empty'}
																	class={cn('size-5', item.icon ? 'visibale' : 'invisible')}
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
								<Form.Field>
									<Form.Label>Scope UUID</Form.Label>
									<SingleInput.General type="text" bind:value={requestFIO.scopeUuid}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Facility Name</Form.Label>
									<SingleInput.General type="text" bind:value={requestFIO.facilityName}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>NFS Endpoint</Form.Label>
									<SingleInput.General type="text" bind:value={requestFIO.nfsEndpoint}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>NFS Path</Form.Label>
									<SingleInput.General type="text" bind:value={requestFIO.nfsPath}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Job Count</Form.Label>
									<SingleInput.General type="number" bind:value={requestFIO.jobCount}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Run Time</Form.Label>
									<SingleInput.General type="text" bind:value={requestFIO.runTime}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Block Size</Form.Label>
									<SingleInput.General type="text" bind:value={requestFIO.blockSize}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>File Size</Form.Label>
									<SingleInput.General type="text" bind:value={requestFIO.fileSize}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>IO Depth</Form.Label>
									<SingleInput.General type="number" bind:value={requestFIO.ioDepth}/>
								</Form.Field>
							</Form.Fieldset>
						{:else if request.input.case == 'warp'}
							<Form.Fieldset>
								<Form.Field>
									<Form.Label for="fio-access-mode">Operation</Form.Label>
									<SingleSelect.Root options={testResultWarpOperation} bind:value={requestWarp.operation}>
										<SingleSelect.Trigger />
										<SingleSelect.Content>
											<SingleSelect.Options>
												<SingleSelect.Input />
												<SingleSelect.List>
													<SingleSelect.Empty>No results found.</SingleSelect.Empty>
													<SingleSelect.Group>
														{#each $testResultWarpOperation as item}
															<SingleSelect.Item option={item}>
																<Icon
																	icon={item.icon ? item.icon : 'ph:empty'}
																	class={cn('size-5', item.icon ? 'visibale' : 'invisible')}
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
								<Form.Field>
									<Form.Label>Endpoint</Form.Label>
									<SingleInput.General type="number" bind:value={requestWarp.endpoint}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Access Key</Form.Label>
									<SingleInput.General type="number" bind:value={requestWarp.accessKey}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Secret Key</Form.Label>
									<SingleInput.General type="number" bind:value={requestWarp.secretKey}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Duration</Form.Label>
									<SingleInput.General type="number" bind:value={requestWarp.duration}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Object Size</Form.Label>
									<SingleInput.General type="number" bind:value={requestWarp.objectSize}/>
								</Form.Field>
								<Form.Field>
									<Form.Label>Object Num</Form.Label>
									<SingleInput.General type="number" bind:value={requestWarp.objectNum}/>
								</Form.Field>
							</Form.Fieldset>
						{/if}
					</Form.Root>
				</MultipleStepModal.Model>

				<!-- Step three -->
				<MultipleStepModal.Model>
					TODO: Content for Tab 3
				</MultipleStepModal.Model>
			</MultipleStepModal.Models>
		</MultipleStepModal.Stepper>
		
		<MultipleStepModal.Footer>
			<MultipleStepModal.Cancel onclick={reset}>Cancel</MultipleStepModal.Cancel>
			<MultipleStepModal.Confirm>Confirm</MultipleStepModal.Confirm>
			<MultipleStepModal.Controllers>
				<MultipleStepModal.Back>Back</MultipleStepModal.Back>
				<MultipleStepModal.Next>Next</MultipleStepModal.Next>
			</MultipleStepModal.Controllers>
		</MultipleStepModal.Footer>
	</MultipleStepModal.Content>
</MultipleStepModal.Root>

