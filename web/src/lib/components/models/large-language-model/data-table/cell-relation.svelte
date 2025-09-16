<script lang="ts">
	import { createClient } from '@connectrpc/connect';
	import { createConnectTransport } from '@connectrpc/connect-web';
	import Icon from '@iconify/svelte';
	import { type Edge, type Node } from '@xyflow/svelte';
	import { writable } from 'svelte/store';

	import type { LargeLangeageModel } from '../protobuf.svelte';

	import { EssentialService, type GpuRelation } from '$lib/api/essential/v1/essential_pb';
	import { Complex } from '$lib/components/custom/flow/index';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { currentKubernetes } from '$lib/stores';
	import '@xyflow/svelte/dist/style.css';

	let { model }: { model: LargeLangeageModel } = $props();
	console.log(model.name);

	const transport = createConnectTransport({
		baseUrl: 'http://10.102.197.18:10888',
	});
	const essentialClient = createClient(EssentialService, transport);

	const relation = writable({} as GpuRelation);
	essentialClient
		.getGpuRelationByModel({
			scopeUuid: $currentKubernetes?.scopeUuid,
			facilityName: $currentKubernetes?.name,
			modelName: 'llama-1b-hf-32768-fpf',
		})
		.then((response) => {
			if (response.gpuRelation) {
				relation.set(response.gpuRelation);
			}
		});

	const position = { x: 0, y: 0 };

	const machines: Node[] = $derived(
		$relation.podInfos.map((podInformation) => ({
			id: podInformation.machineName,
			type: 'machine',
			data: {
				name: podInformation.machineName,
				icon: 'simple-icons:maas',
			},
			position,
		})),
	);
	const gpus: Node[] = $derived(
		$relation.podInfos
			.flatMap((podInformation) => podInformation.vgpus)
			.map((gpu) => ({
				id: gpu.physicalGpuUuid,
				type: 'gpu',
				data: {
					name: gpu.physicalGpuUuid,
					model: gpu.physicalGpuUuid,
					icon: 'simple-icons:nvidia',
				},
				position,
			})),
	);

	const nodes = $derived([...machines, ...gpus]);

	const edges: Edge[] = $derived(
		$relation.podInfos.flatMap((podInformation) =>
			podInformation.vgpus.map((gpu) => ({
				id: `${podInformation.machineName}${gpu.physicalGpuUuid}`,
				type: 'edge',
				source: podInformation.machineName,
				target: gpu.physicalGpuUuid,
				animated: true,
				selectable: false,
			})),
		),
	);

	let open = $state(false);
</script>

<!-- {machines.length}
{gpus.length} -->
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
