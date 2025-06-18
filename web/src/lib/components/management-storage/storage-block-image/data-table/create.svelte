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
	export const pools: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'rbd',
			label: 'RADOS Block Device',
			icon: 'ph:cube'
		}
	]);
	export const features: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'deep_flatten',
			label: 'Deep Flatten',
			icon: 'ph:command',
			information: `Feature can be disabled but can't be re-enabled later`
		},
		{
			value: 'layering',
			label: 'Layering',
			icon: 'ph:command',
			information: `Feature flag can't be manipulated after the image is created. Disabling this option will also disable the Protect and Clone actions on Snapshot`
		},
		{
			value: 'exclusive_lock',
			label: 'Exclusive lock',
			icon: 'ph:command'
		},
		{
			value: 'object_map',
			label: 'Object map',
			icon: 'ph:command',
			information: `requires exclusive-lock`
		},
		{
			value: 'fast_diff',
			label: 'Fast Diff',
			icon: 'ph:command',
			information: `interlocked with object-map`
		}
	]);
	export type Request = {
		name: string;
		pool: string;
		size: number;
		features: string[];
		objectSize: number;
		stripeUnit: number;
		stripeCount: number;
		bpsLimit: number;
		iopsLimit: number;
		readBpsLimit: number;
		readIopsLimit: number;
		writeBpsLimit: number;
		writeIopsLimit: number;
		bpsBurst: number;
		iopsBurst: number;
		readBpsBurst: number;
		readIopsBurst: number;
		writeBpsBurst: number;
		writeIopsBurst: number;
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
		<AlertDialog.Trigger class={cn(buttonVariants({ variant: 'default', size: 'sm' }))}>
			<div class="flex items-center gap-1">
				<Icon icon="ph:plus" />
				Create
			</div>
		</AlertDialog.Trigger>
	</div>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Create RADOS Block Device
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label for="rados-block-device-name">Name</Form.Label>
					<SingleInput.General
						required
						type="text"
						id="rados-block-device-name"
						bind:value={request.name}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-type">Pool</Form.Label>
					<SingleSelect.Root required options={pools} bind:value={request.pool}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>No results found.</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $pools as option}
											<SingleSelect.Item {option}>
												<Icon
													icon={option.icon ? option.icon : 'ph:empty'}
													class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
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
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-size">Size</Form.Label>
					<SingleInput.General
						required
						type="number"
						id="rados-block-device-size"
						bind:value={request.size}
					/>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Field>
					<Form.Label for="rados-block-device-features">Features</Form.Label>
					<MultipleSelect.Root bind:value={request.features} options={features}>
						<MultipleSelect.Viewer />
						<MultipleSelect.Controller>
							<MultipleSelect.Trigger />
							<MultipleSelect.Content>
								<MultipleSelect.Options>
									<MultipleSelect.Input />
									<MultipleSelect.List>
										<MultipleSelect.Empty>No results found.</MultipleSelect.Empty>
										<MultipleSelect.Group>
											{#each $features as option}
												<MultipleSelect.Item {option}>
													<Icon
														icon={option.icon ? option.icon : 'ph:empty'}
														class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
													/>
													{option.label}
													<MultipleSelect.Check {option} />
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

			<Form.Fieldset>
				<Form.Legend>Striping</Form.Legend>

				<Form.Field>
					<Form.Label for="rados-block-device-object-size">Object Size</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-object-size"
						bind:value={request.size}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-stripe-unit">Stripe Unit</Form.Label>
					<SingleInput.General
						required
						type="number"
						id="rados-block-device-stripe-unit"
						bind:value={request.size}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-stripe-count">Stripe Count</Form.Label>
					<SingleInput.General
						required
						type="number"
						id="rados-block-device-stripe-count"
						bind:value={request.size}
					/>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Quality of Service</Form.Legend>

				<Form.Field>
					<Form.Label for="rados-block-device-bps-limit">BPS Limit</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-bps-limit"
						bind:value={request.bpsLimit}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-iops-limit">IOPS Limit</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-iops-limit"
						bind:value={request.iopsLimit}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-read-bps-limit">Read BPS Limit</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-read-bps-limit"
						bind:value={request.readBpsLimit}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-read-iops-limit">Read IOPS Limit</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-read-iops-limit"
						bind:value={request.readIopsLimit}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-write-bps-limit">Write BPS Limit</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-write-bps-limit"
						bind:value={request.writeBpsLimit}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-write-iops-limit">Write IOPS Limit</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-write-iops-limit"
						bind:value={request.writeIopsLimit}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-bps-burst">BPS Burst</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-bps-burst"
						bind:value={request.bpsBurst}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-iops-burst">IOPS Burst</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-iops-burst"
						bind:value={request.iopsBurst}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-read-bps-burst">Read BPS Burst</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-read-bps-burst"
						bind:value={request.readBpsBurst}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-read-iops-burst">Read IOPS Burst</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-read-iops-burst"
						bind:value={request.readIopsBurst}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-write-bps-burst">Write BPS Burst</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-write-bps-burst"
						bind:value={request.writeBpsBurst}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label for="rados-block-device-write-iops-burst">Write IOPS Burst</Form.Label>
					<SingleInput.General
						type="number"
						id="rados-block-device-write-iops-burst"
						bind:value={request.writeIopsBurst}
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
