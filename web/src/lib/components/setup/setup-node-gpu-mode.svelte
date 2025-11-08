<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { Facility_Unit } from '$lib/api/facility/v1/facility_pb';
	import { MachineService } from '$lib/api/machine/v1/machine_pb';
	import { OrchestratorService } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';

	const gpuModeOptions = [
		{
			value: 'vgpu',
			label: 'Virtual'
		},
		{
			value: 'vm-passthrough',
			label: 'Passthrough',
			disabled: true
		}
	];
</script>

<script lang="ts">
	let { unit, class: className }: { unit: Facility_Unit; class: string } = $props();

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const orchestratorClient = createClient(OrchestratorService, transport);

	let hasGPUs: undefined | boolean = $state(undefined);
	let selectedGPUMode: string = $state('');

	async function fetch() {
		orchestratorClient
			.listKubernetesNodeLabels({
				scope: $currentKubernetes?.scope,
				facility: $currentKubernetes?.name,
				hostname: unit.hostname,
				all: true
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
			console.error('Failed to get machine information:', error);
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
									toast.promise(
										() => {
											selectedGPUMode = option.value;
											return orchestratorClient.updateKubernetesNodeLabels({
												scope: $currentKubernetes?.scope,
												facility: $currentKubernetes?.name,
												hostname: unit.hostname,
												labels: {
													'nvidia.com/gpu.workload.config': selectedGPUMode
												}
											});
										},

										{
											loading: 'Loading...',
											success: () => {
												fetch();
												return `Set ${unit.hostname} as ${selectedGPUMode}`;
											},
											error: (error) => {
												let message = `Failed to set ${unit.hostname} as ${selectedGPUMode}`;
												toast.error(message, {
													description: (error as ConnectError).message.toString(),
													duration: Number.POSITIVE_INFINITY
												});
												return message;
											}
										}
									);
									close();
								}}
							>
								<Icon
									icon="ph:check"
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
