<script lang="ts" module>
	import { createClient } from '@connectrpc/connect';
	import { createConnectTransport } from '@connectrpc/connect-web';
	import Icon from '@iconify/svelte';
	import { type Edge, type Node } from '@xyflow/svelte';
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import type { LargeLangeageModel } from '../type';

	import { EssentialService, type GpuRelation } from '$lib/api/essential/v1/essential_pb';
	import { Complex } from '$lib/components/custom/flow/index';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { currentKubernetes } from '$lib/stores';
	import '@xyflow/svelte/dist/style.css';

	const position = { x: 0, y: 0 };
</script>

<script lang="ts">
	let { model }: { model: LargeLangeageModel } = $props();

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
				.getGpuRelationByModel({
					scopeUuid: $currentKubernetes?.scopeUuid,
					facilityName: $currentKubernetes?.name,
					modelName: model.name,
				})
				.then((response) => {
					if (response.gpuRelation) {
						relation.set(response.gpuRelation);
						isMounted = true;
					}
				})
				.catch((error) => {
					console.log(essentialClient);
					console.error(`Failed to fetch relation of model ${model.name}:`, error);
				});
		} catch (error) {
			console.error(error);
		}
	});
</script>

{#if !isMounted}
	Loading...
{:else}
	<Dialog.Root bind:open>
		<Dialog.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
			<Icon icon="ph:graph" />
		</Dialog.Trigger>
		{#if open}
			<Dialog.Content class="min-h-[77vh] min-w-[77vw]">
				<Complex.Flow initialNodes={nodes} initialEdges={edges} />
			</Dialog.Content>
		{/if}
	</Dialog.Root>
{/if}
