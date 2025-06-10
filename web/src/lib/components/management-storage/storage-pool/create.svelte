<script lang="ts" module>
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import {
		Multiple as MultipleSelect,
		Single as SingleSelect
	} from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { writable, type Writable } from 'svelte/store';

	export const poolTypes: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'erasure',
			label: 'Erasure',
			icon: 'ph:scales'
		},
		{
			value: 'replicated',
			label: 'Replicated',
			icon: 'ph:copy-simple'
		}
	]);
	export const PGAutoscales: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'on',
			label: 'ON',
			icon: 'ph:toggle-right'
		},
		{
			value: 'off',
			label: 'OFF',
			icon: 'ph:toggle-left'
		},
		{
			value: 'warn',
			label: 'WARN',
			icon: 'ph:warning-circle'
		}
	]);
	export const profiles: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'default',
			label: 'Default',
			icon: 'ph:squares-four'
		}
	]);
	export const applications: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'cephfs',
			label: 'Ceph File System',
			icon: 'ph:squares-four'
		},
		{
			value: 'rbd',
			label: 'RADOS Block Device',
			icon: 'ph:squares-four'
		},
		{
			value: 'rgw',
			label: 'RADOS Gateway',
			icon: 'ph:squares-four'
		}
	]);
	export const modes: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'none',
			label: 'None',
			icon: 'ph:list'
		},
		{
			value: 'passive',
			label: 'Passive',
			icon: 'ph:list'
		},
		{
			value: 'aggressive',
			label: 'Aggressive',
			icon: 'ph:list'
		},
		{
			value: 'force',
			label: 'Force',
			icon: 'ph:list'
		}
	]);
	export const algorithms: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'snappy',
			label: 'Snappy',
			icon: 'ph:math-operations'
		},
		{
			value: 'zlib',
			label: 'ZLIB',
			icon: 'ph:math-operations'
		},
		{
			value: 'zstd',
			label: 'ZSTD',
			icon: 'ph:math-operations'
		},
		{
			value: 'lz4',
			label: 'LZ4',
			icon: 'ph:math-operations'
		}
	]);
	export type Request = {
		name: string;
		poolType: string;
		pgAutoScale: string;
		flags: boolean;
		replicatedSize: number;
		applications: string[];
		placementGroups: number;

		erasureCodeProfile: string;
		crushRuleset: string[];

		mode: string;
		algorithm: string;
		minimumBlobSize: string;
		maximumBlobSize: string;
		ratio: number;

		maxBytes: string;
		maxObjects: number;
	};
</script>

<script lang="ts">
	const DEFAULT_REQUEST = {} as Request;
	let request: Request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<div class="flex justify-end">
		<AlertDialog.Trigger class={cn(buttonVariants({ variant: 'default' }))}>
			<div class="flex items-center gap-2">
				<Icon icon="ph:plus" />
				<p class="text-base">Create</p>
			</div>
		</AlertDialog.Trigger>
	</div>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Create Pool
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label for="pool-name">Name</Form.Label>
					<SingleInput.General required type="text" id="pool-name" bind:value={request.name} />
				</Form.Field>

				<Form.Field>
					<Form.Label for="pool-type">Pool Type</Form.Label>
					<SingleSelect.Root required options={poolTypes} bind:value={request.poolType}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>No results found.</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $poolTypes as type}
											<SingleSelect.Item option={type}>
												<Icon
													icon={type.icon ? type.icon : 'ph:empty'}
													class={cn('size-5', type.icon ? 'visibale' : 'invisible')}
												/>
												{type.label}
												<SingleSelect.Check option={type} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>

				<Form.Field>
					<Form.Label for="pool-pg-autoscale">PG Autoscale</Form.Label>
					<SingleSelect.Root options={PGAutoscales} bind:value={request.pgAutoScale}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>No results found.</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $PGAutoscales as autoscale}
											<SingleSelect.Item option={autoscale} class="text-xs">
												<Icon
													icon={autoscale.icon ? autoscale.icon : 'ph:empty'}
													class={cn('size-5', autoscale.icon ? 'visibale' : 'invisible')}
												/>
												{autoscale.label}
												<SingleSelect.Check option={autoscale} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>

				{#if request.poolType === 'erasure'}
					<Form.Field>
						<Form.Label for="pool-flags">Flags</Form.Label>
						<SingleInput.Boolean
							required
							id="pool-flags"
							descriptor={(value) => {
								if (value === true) {
									return 'EC Overwrites';
								} else if (value === false) {
									return 'EC Not Overwrites';
								} else {
									return 'Undetermined';
								}
							}}
							bind:value={request.flags}
						/>
					</Form.Field>
				{/if}

				{#if request.pgAutoScale !== 'on'}
					<Form.Field>
						<Form.Label for="pool-placement-groups">Placement Groups</Form.Label>
						<SingleInput.General
							required
							type="number"
							id="pool-placement-groups"
							bind:value={request.placementGroups}
						/>
					</Form.Field>
				{/if}

				{#if request.poolType === 'replicated'}
					<Form.Field>
						<Form.Label for="pool-replicated-size">Replcated Size</Form.Label>
						<SingleInput.General
							required
							type="number"
							id="pool-replicated-size"
							bind:value={request.replicatedSize}
						/>
					</Form.Field>
					<Form.Help>
						A size of 1 will not create a replication of the object. The 'Replicated size' includes
						the object itself.
					</Form.Help>
				{/if}

				<Form.Field>
					<Form.Label for="pool-applications">Applications</Form.Label>
					<MultipleSelect.Root bind:value={request.applications} options={applications}>
						<MultipleSelect.Viewer />
						<MultipleSelect.Controller>
							<MultipleSelect.Trigger />
							<MultipleSelect.Content>
								<MultipleSelect.Options>
									<MultipleSelect.Input>
										{#snippet addition({ accessor })}
											<AlertDialog.Root>
												<AlertDialog.Trigger class={cn(buttonVariants({ variant: 'outline' }))}>
													<div class="gap- flex items-center">
														<Icon icon="ph:plus" />
														<p class="text-xs">Add</p>
													</div>
												</AlertDialog.Trigger>
												<AlertDialog.Content>
													<AlertDialog.Header
														class="flex items-center justify-center text-xl font-bold"
													>
														Create Application
													</AlertDialog.Header>
													<Form.Root>
														<Form.Fieldset>
															<Form.Field>
																<Form.Label>Name</Form.Label>
																<SingleInput.General required bind:value={accessor.input} />
															</Form.Field>
														</Form.Fieldset>
													</Form.Root>
													<AlertDialog.Footer>
														<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
														<AlertDialog.Action
															onclick={() => {
																applications.update((current) => [
																	...current,
																	{
																		value: accessor.input,
																		label: accessor.input,
																		icon: 'ph:squares-four'
																	} as MultipleSelect.OptionType
																]);
																accessor.input = '';
															}}
														>
															Confirm
														</AlertDialog.Action>
													</AlertDialog.Footer>
												</AlertDialog.Content>
											</AlertDialog.Root>
										{/snippet}
									</MultipleSelect.Input>
									<MultipleSelect.List>
										<MultipleSelect.Empty>No results found.</MultipleSelect.Empty>
										<MultipleSelect.Group>
											{#each $applications as application}
												<MultipleSelect.Item option={application}>
													<Icon
														icon={application.icon ? application.icon : 'ph:empty'}
														class={cn('size-5', application.icon ? 'visibale' : 'invisible')}
													/>
													{application.label}
													<MultipleSelect.Check option={application} />
												</MultipleSelect.Item>
											{/each}
										</MultipleSelect.Group>
									</MultipleSelect.List>
									<MultipleSelect.Actions>
										<MultipleSelect.ActionAll>All</MultipleSelect.ActionAll>
										<MultipleSelect.ActionClear>Clear</MultipleSelect.ActionClear>
									</MultipleSelect.Actions>
								</MultipleSelect.Options>
							</MultipleSelect.Content>
						</MultipleSelect.Controller>
					</MultipleSelect.Root>
				</Form.Field>
			</Form.Fieldset>

			{#if request.poolType}
				<Form.Fieldset>
					<Form.Legend>Crush</Form.Legend>

					{#if request.poolType === 'erasure'}
						<Form.Field>
							<Form.Label for="pool-erasure-code-profile">Erasure Code Profile</Form.Label>
							<SingleSelect.Root options={profiles} bind:value={request.erasureCodeProfile}>
								<SingleSelect.Trigger />
								<SingleSelect.Content>
									<SingleSelect.Options>
										<SingleSelect.Input />
										<SingleSelect.List>
											<SingleSelect.Empty>No results found.</SingleSelect.Empty>
											<SingleSelect.Group>
												{#each $profiles as profile}
													<SingleSelect.Item option={profile}>
														<Icon
															icon={profile.icon ? profile.icon : 'ph:empty'}
															class={cn('size-5', profile.icon ? 'visibale' : 'invisible')}
														/>
														{profile.label}
														<SingleSelect.Check option={profile} />
													</SingleSelect.Item>
												{/each}
											</SingleSelect.Group>
										</SingleSelect.List>
									</SingleSelect.Options>
								</SingleSelect.Content>
							</SingleSelect.Root>
						</Form.Field>
					{/if}
					<Form.Field>
						<Form.Label for="pool-crush-ruleset">Crush Ruleset</Form.Label>
						{#if request.poolType === 'erasure'}
							<Form.Help>A new crush ruleset will be implicitly created.</Form.Help>
						{:else if request.poolType === 'replicated'}
							<SingleSelect.Root options={profiles} bind:value={request.erasureCodeProfile}>
								<SingleSelect.Trigger />
								<SingleSelect.Content>
									<SingleSelect.Options>
										<SingleSelect.Input />
										<SingleSelect.List>
											<SingleSelect.Empty>No results found.</SingleSelect.Empty>
											<SingleSelect.Group>
												{#each $profiles as profile}
													<SingleSelect.Item option={profile}>
														<Icon
															icon={profile.icon ? profile.icon : 'ph:empty'}
															class={cn('size-5', profile.icon ? 'visibale' : 'invisible')}
														/>
														{profile.label}
														<SingleSelect.Check option={profile} />
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
			{/if}

			<Form.Fieldset>
				<Form.Legend>Compression</Form.Legend>

				<Form.Field>
					<Form.Label for="pool-mode">Mode</Form.Label>
					<SingleSelect.Root options={modes} bind:value={request.mode}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>No results found.</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $modes as mode}
											<SingleSelect.Item option={mode}>
												<Icon
													icon={mode.icon ? mode.icon : 'ph:empty'}
													class={cn('size-5', mode.icon ? 'visibale' : 'invisible')}
												/>
												{mode.label}
												<SingleSelect.Check option={mode} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>

				{#if request.mode && request.mode !== 'none'}
					<Form.Field>
						<Form.Label for="pool-algorithm">Algorithm</Form.Label>
						<SingleSelect.Root options={algorithms} bind:value={request.algorithm}>
							<SingleSelect.Trigger />
							<SingleSelect.Content>
								<SingleSelect.Options>
									<SingleSelect.Input />
									<SingleSelect.List>
										<SingleSelect.Empty>No results found.</SingleSelect.Empty>
										<SingleSelect.Group>
											{#each $algorithms as algorithm}
												<SingleSelect.Item option={algorithm}>
													<Icon
														icon={algorithm.icon ? algorithm.icon : 'ph:empty'}
														class={cn('size-5', algorithm.icon ? 'visibale' : 'invisible')}
													/>
													{algorithm.label}
													<SingleSelect.Check option={algorithm} />
												</SingleSelect.Item>
											{/each}
										</SingleSelect.Group>
									</SingleSelect.List>
								</SingleSelect.Options>
							</SingleSelect.Content>
						</SingleSelect.Root>
					</Form.Field>

					<Form.Field>
						<Form.Label for="pool-minimum-blob-size">Minimum Blob Size</Form.Label>
						<SingleInput.General
							type="text"
							id="pool-minimum-blob-size"
							bind:value={request.minimumBlobSize}
						/>
					</Form.Field>

					<Form.Field>
						<Form.Label for="pool-maximum-blob-size">Maximum Blob Size</Form.Label>
						<SingleInput.General
							type="text"
							id="pool-maximum-blob-size"
							bind:value={request.maximumBlobSize}
						/>
					</Form.Field>

					<Form.Field>
						<Form.Label for="pool-ratio">Ratio</Form.Label>
						<SingleInput.General type="number" id="pool-ratio" bind:value={request.ratio} />
					</Form.Field>
				{/if}
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Quotas</Form.Legend>

				<Form.Field>
					<Form.Label for="pool-max-bytes">Max Bytes</Form.Label>
					<SingleInput.General type="text" id="pool-max-bytes" bind:value={request.maxBytes} />
				</Form.Field>
				<Form.Help>
					Max objects 0 Leave it blank or specify 0 to disable this quota. A valid quota should be
					greater than 0.
				</Form.Help>

				<Form.Field>
					<Form.Label for="pool-max-objects">Max Objects</Form.Label>
					<SingleInput.General
						type="number"
						id="pool-max-objects"
						bind:value={request.maxObjects}
					/>
				</Form.Field>
				<Form.Help>
					Leave it blank or specify 0 to disable this quota. A valid quota should be greater than 0.
				</Form.Help>
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
