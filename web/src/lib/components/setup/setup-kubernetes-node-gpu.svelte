<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { MachineService } from '$lib/api/machine/v1/machine_pb';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { cn } from '$lib/utils';

	const gpuOptions = [
		{
			value: 'virtual',
			label: 'Virtual',
		},
		{
			value: 'passthrough',
			label: 'Passthrough',
		},
	];
</script>

<script lang="ts">
	let {
		machine_id,
		unit_name,
		class: className,
	}: { machine_id: string; unit_name: string; class: string } = $props();

	let VALUE = 'virtual';

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	let selectedGPU: undefined | string = $state(undefined);

	async function set(value: string) {
		VALUE = value;
	}
	async function fetch() {
		selectedGPU = VALUE;
	}

	let hasGPUs: undefined | boolean = $state(undefined);

	let open = $state(false);
	function close() {
		open = false;
	}

	onMount(async () => {
		try {
			await machineClient.getMachine({ id: machine_id }).then((response) => {
				if (response.gpuDevices.length > 0) {
					hasGPUs = true;
					fetch();
				} else {
					hasGPUs = false;
				}
			});
		} catch (error) {
			console.error('Failed to get machine informationo:', error);
		}
	});
</script>

{#if hasGPUs}
	<Popover.Root bind:open>
		<Popover.Trigger
			onclick={(e) => {
				e.stopPropagation();
			}}
			class={cn('flex items-center justify-center', className)}
		>
			<Icon icon="ph:graphics-card" />
		</Popover.Trigger>
		<Popover.Content class="w-[200px] p-0">
			<Command.Root>
				<Command.List>
					<Command.Group>
						{#each gpuOptions as option}
							<Command.Item
								value={option.value}
								onSelect={() => {
									toast.promise(() => set(option.value), {
										loading: 'Loading...',
										success: () => {
											fetch();
											return `Set ${unit_name} as ${selectedGPU}`;
										},
										error: (e) => {
											let msg = `Fail to set  ${unit_name} as ${selectedGPU}`;
											toast.error(msg, {
												description: (e as ConnectError).message.toString(),
												duration: Number.POSITIVE_INFINITY,
											});
											return msg;
										},
									});
									close();
								}}
							>
								<Icon
									icon="ph:graphics-card"
									class={option.value == selectedGPU ? 'visible' : 'invisible'}
								/>
								{option.label}
							</Command.Item>
						{/each}
					</Command.Group>
				</Command.List>
			</Command.Root>
		</Popover.Content>
	</Popover.Root>
{/if}
