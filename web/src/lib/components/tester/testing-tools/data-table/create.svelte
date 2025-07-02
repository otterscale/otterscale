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

	
	export type Request = {
		name: string;
		rwMode: string;
		fileSize: string;
		numberJobs: number;
		blockSize: string;
		runtime: number;

		// optional
		timeBased: boolean;
		exitallOnError: boolean;
		createSerialize: boolean;
		norandommap: boolean;
		direct: boolean;
		groupReporting: boolean;

		filenameFormat: string;
		directory:  string;
		clocksource: string;
		ioengine: string;
		diskUtil: string;

		startDelay: number;
		rwmixread: number;
		ioDepth: number;
		bufferCompressPercentage: number;
	};

	const readWriteMode: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'read',
			label: 'Read',
			icon: 'ph:upload'
		},
		{
			value: 'randread',
			label: 'Random  Read',
			icon: 'ph:upload'
		},
		{
			value: 'write',
			label: 'Write',
			icon: 'ph:download'
		},
		{
			value: 'randwrite',
			label: 'Random Write',
			icon: 'ph:download'
		}
	]);

	const DEFAULT_REQUEST = { 
		timeBased: true,
		exitallOnError: true,
		createSerialize: true,
		norandommap: true,
		direct: true,
		groupReporting: true,
	} as Request;
	let request: Request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<div class="flex justify-end">
		<AlertDialog.Trigger class={cn(buttonVariants({ variant: 'default', size: 'sm' }))}>
			<div class="flex items-center gap-1">
				<Icon icon="ph:plus" />
				Create
			</div>
		</AlertDialog.Trigger>
	</div>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Create FIO config
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label for="fio-name">Name</Form.Label>
					<SingleInput.General
						required
						type="text"
						id="fio-name"
						bind:value={request.name}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label for="fio-read-write-mode">Read Write Mode</Form.Label>
					<SingleSelect.Root options={readWriteMode} bind:value={request.rwMode}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>No results found.</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $readWriteMode as rw}
											<SingleSelect.Item option={rw}>
												<Icon
													icon={rw.icon ? rw.icon : 'ph:empty'}
													class={cn('size-5', rw.icon ? 'visibale' : 'invisible')}
												/>
												{rw.label}
												<SingleSelect.Check option={rw} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>
				<Form.Field>
					<Form.Label for="fio-file-size">File Size</Form.Label>
					<SingleInput.General
						required
						type="text"
						id="fio-file-size"
						bind:value={request.fileSize}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label for="fio-number-jobs">Number Jobs</Form.Label>
					<SingleInput.General
						required
						type="number"
						id="fio-number-jobs"
						bind:value={request.numberJobs}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label for="fio-block-size">Block Size</Form.Label>
					<SingleInput.General
						required
						type="text"
						id="fio-block-size"
						bind:value={request.blockSize}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label for="fio-runtime">Runtime</Form.Label>
					<SingleInput.General
						required
						type="number"
						id="fio-runtime"
						bind:value={request.runtime}
					/>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Optional</Form.Legend>
				<!-- Boolean -->
				<Form.Field>
					<Form.Label>Time Based</Form.Label>
					<SingleInput.Boolean bind:value={request.timeBased} />
				</Form.Field>
				<Form.Field>
					<Form.Label>Exitall On Error</Form.Label>
					<SingleInput.Boolean bind:value={request.exitallOnError} />
				</Form.Field>
				<Form.Field>
					<Form.Label>Create Serialize</Form.Label>
					<SingleInput.Boolean bind:value={request.createSerialize} />
				</Form.Field>
				<Form.Field>
					<Form.Label>norandommap</Form.Label>
					<SingleInput.Boolean bind:value={request.norandommap} />
				</Form.Field>
				<Form.Field>
					<Form.Label>Direct</Form.Label>
					<SingleInput.Boolean bind:value={request.direct} />
				</Form.Field>
				<Form.Field>
					<Form.Label>Group Reporting</Form.Label>
					<SingleInput.Boolean bind:value={request.groupReporting} />
				</Form.Field>

				<!-- String -->
				<Form.Field>
					<Form.Label for="fio-filename-format">Filename Format</Form.Label>
					<SingleInput.General
						type="text"
						id="fio-filename-format"
						bind:value={request.filenameFormat}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label for="fio-directory">Directory</Form.Label>
					<SingleInput.General
						type="text"
						id="fio-directory"
						bind:value={request.directory}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label for="fio-clocksource">Clocksource</Form.Label>
					<SingleInput.General
						type="text"
						id="fio-clocksource"
						bind:value={request.clocksource}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label for="fio-ioengine">I/O Engine</Form.Label>
					<SingleInput.General
						type="text"
						id="fio-ioengine"
						bind:value={request.ioengine}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label for="fio-disk-util">Disk Util</Form.Label>
					<SingleInput.General
						type="text"
						id="fio-disk-util"
						bind:value={request.diskUtil}
					/>
				</Form.Field>

				<!-- Number -->
				<Form.Field>
					<Form.Label for="fio-start-delay">Start Delay</Form.Label>
					<SingleInput.General
						type="number"
						id="fio-start-delay"
						bind:value={request.startDelay}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label for="fio-rwmixread">Mixed Read Write</Form.Label>
					<SingleInput.General
						type="number"
						id="fio-rwmixread"
						bind:value={request.rwmixread}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label for="fio-io-depth">IO Depth</Form.Label>
					<SingleInput.General
						type="number"
						id="fio-io-depth"
						bind:value={request.ioDepth}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label for="fio-buffer-compress-percentage">Buffer Compress Percentage</Form.Label>
					<SingleInput.General
						type="number"
						id="fio-buffer-compress-percentage"
						bind:value={request.bufferCompressPercentage}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						console.log(request);
					}}
				>
					Create
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
