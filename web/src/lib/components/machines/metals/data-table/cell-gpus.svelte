<script lang="ts" module>
	import { createClient } from '@connectrpc/connect';
	import { createConnectTransport } from '@connectrpc/connect-web';
	import Icon from '@iconify/svelte';
	import { type Edge, type Node } from '@xyflow/svelte';
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { EssentialService, type GpuRelation } from '$lib/api/essential/v1/essential_pb';
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Complex } from '$lib/components/custom/flow/index';
	import * as Table from '$lib/components/custom/table';
	import * as Sheet from '$lib/components/ui/sheet';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
	import '@xyflow/svelte/dist/style.css';

	const position = { x: 0, y: 0 };
</script>

<script lang="ts">
	let {
		machine,
	}: {
		machine: Machine;
	} = $props();

	const transport = createConnectTransport({
		baseUrl: 'http://10.102.197.18:10888',
	});
	const essentialClient = createClient(EssentialService, transport);

	const relation = writable({} as GpuRelation);
	const machines: Node[] = $derived(
		$relation.podInfos.map((podInformation) => ({
			id: podInformation.machineName,
			type: 'machine',
			data: podInformation,
			position,
		})),
	);
	const gpus: Node[] = $derived(
		$relation.podInfos
			.flatMap((podInformation) => podInformation.vgpus)
			.map((gpu) => ({
				id: gpu.physicalGpuUuid,
				type: 'gpu',
				data: gpu,
				position,
			})),
	);
	const models: Node[] = $derived(
		$relation.podInfos
			.filter((podInformation) => podInformation.modelName)
			.map((podInformation) => ({
				id: podInformation.modelName,
				type: 'model',
				data: podInformation,
				position,
			})),
	);
	const machineGPUs: Edge[] = $derived(
		$relation.podInfos.flatMap((podInformation) =>
			podInformation.vgpus.map((gpu) => ({
				id: `${podInformation.machineName}${gpu.physicalGpuUuid}`,
				type: 'edge',
				source: gpu.physicalGpuUuid,
				target: podInformation.machineName,
				animated: true,
				selectable: false,
			})),
		),
	);
	const gpuModels: Edge[] = $derived(
		$relation.podInfos
			.filter((podInformation) => podInformation.modelName)
			.flatMap((podInformation) =>
				podInformation.vgpus.map((gpu) => ({
					id: `${gpu.physicalGpuUuid}${podInformation.modelName}`,
					type: 'edge',
					source: podInformation.modelName,
					target: gpu.physicalGpuUuid,
					animated: true,
					selectable: false,
				})),
			),
	);
	const nodes = $derived([...machines, ...gpus, ...models]);
	const edges: Edge[] = $derived([...machineGPUs, ...gpuModels]);

	let open = $state(false);
	let isMounted = $state(false);
	onMount(async () => {
		try {
			essentialClient
				.getGpuRelationByMachine({
					scopeUuid: $currentKubernetes?.scopeUuid,
					facilityName: $currentKubernetes?.name,
					machineName: machine.hostname,
				})
				.then((response) => {
					if (response.gpuRelation) {
						relation.set(response.gpuRelation);
						isMounted = true;
					}
				})
				.catch((error) => {
					console.log(essentialClient);
					console.error(`Failed to fetch relation of machine ${machine.hostname}:`, error);
				});
		} catch (error) {
			console.error(error);
		}
	});
</script>

{#if isMounted}
	<Sheet.Root bind:open>
		<Sheet.Trigger>
			<span class="flex items-center gap-1">
				{machine.gpuDevices.length}
				<Icon icon="ph:arrow-square-out" />
			</span>
		</Sheet.Trigger>
		<Sheet.Content side="right" class="min-w-[38vw] p-4">
			{#if open}
				<Sheet.Header>
					<Sheet.Title class="text-center text-lg">{m.details()}</Sheet.Title>
				</Sheet.Header>
				<div class="rounded-lg border shadow">
					<Table.Root>
						<Table.Header>
							<Table.Row
								class="bg-muted [&_th]:bg-muted [&_th]:first:rounded-tl-lg [&_th]:last:rounded-tr-lg"
							>
								<Table.Head>{m.product()}</Table.Head>
								<Table.Head>{m.vendor()}</Table.Head>
								<Table.Head>{m.bus()}</Table.Head>
								<Table.Head>{m.pci_address()}</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each machine.gpuDevices as gpuDevice}
								<Table.Row>
									<Table.Cell
										>{gpuDevice.productName !== ''
											? gpuDevice.productName
											: gpuDevice.productId}</Table.Cell
									>
									<Table.Cell
										>{gpuDevice.vendorName !== ''
											? gpuDevice.vendorName
											: gpuDevice.vendorId}</Table.Cell
									>
									<Table.Cell>{gpuDevice.busName}</Table.Cell>
									<Table.Cell>{gpuDevice.pciAddress}</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</div>
				<div class="flex h-full items-center justify-center">
					<Complex.Flow initialNodes={nodes} initialEdges={edges} class="max-h-[50vh]" />
				</div>
			{/if}
		</Sheet.Content>
	</Sheet.Root>
{/if}
