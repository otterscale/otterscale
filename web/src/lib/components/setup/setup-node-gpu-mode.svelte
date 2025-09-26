<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { EssentialService } from '$lib/api/essential/v1/essential_pb';
	import type { Facility_Unit } from '$lib/api/facility/v1/facility_pb';
	import { MachineService } from '$lib/api/machine/v1/machine_pb';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';

	const gpuModeOptions = [
		{
			value: 'vgpu',
			label: 'Virtual',
		},
		{
			value: 'vm-passthrough',
			label: 'Passthrough',
		},
	];
</script>

<script lang="ts">
	let { unit, class: className }: { unit: Facility_Unit; class: string } = $props();

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const essentialClient = createClient(EssentialService, transport);

	let hasGPUs: undefined | boolean = $state(undefined);
	let selectedGPUMode: undefined | string = $state(undefined);

	async function set() {}
	async function fetch() {
		essentialClient
			.listKubernetesNodeLabels({
				scopeUuid: $currentKubernetes?.scopeUuid,
				facilityName: $currentKubernetes?.name,
				hostname: 'proxmox-4090x2-197-114',
				all: true,
			})
			.then((response) => {
				selectedGPUMode = response.labels['nvidia.com/gpu.workload.config'];
			});
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	onMount(async () => {
		try {
			await machineClient.getMachine({ id: unit.machineId }).then((response) => {
				hasGPUs = response.gpuDevices.length > 0;
			});

			if (hasGPUs) {
				await fetch();
			}
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
						{#each gpuModeOptions as option}
							<Command.Item
								value={option.value}
								onSelect={() => {
									toast.promise(() => set(), {
										loading: 'Loading...',
										success: () => {
											fetch();
											return `Set ${unit.hostname} as ${selectedGPUMode}`;
										},
										error: (e) => {
											let msg = `Fail to set  ${unit.hostname} as ${selectedGPUMode}`;
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
									class={option.value == selectedGPUMode ? 'visible' : 'invisible'}
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
