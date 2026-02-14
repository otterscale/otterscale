<script lang="ts">
	import '@xyflow/svelte/dist/style.css';

	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { type Edge, type Node } from '@xyflow/svelte';
	import { getContext } from 'svelte';
	import { type Writable, writable } from 'svelte/store';

	import type { Model } from '$lib/api/model/v1/model_pb';
	import {
		type GPURelation_GPU,
		type GPURelation_Machine,
		type GPURelation_Pod,
		OrchestratorService
	} from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { Complex as ComplexFlow } from '$lib/components/flow/index';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Drawer from '$lib/components/ui/drawer/index.js';
	import { m } from '$lib/paraglide/messages';

	let { scope, model }: { scope: string; model: Model } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(OrchestratorService, transport);

	const position = { x: 0, y: 0 };

	const nodes: Writable<Node[]> = writable([]);
	const edges: Writable<Edge[]> = writable([]);

	async function fetch() {
		try {
			const response = await client.listGPURelationsByModel({
				scope: scope,
				namespace: model.namespace,
				modelName: model.name
			});
			nodes.set(
				response.gpuRelations.map((gpuRelation) => {
					if (gpuRelation.entity.case === 'machine') {
						const gpus = response.gpuRelations
							.filter((gpuRelation) => gpuRelation.entity.case === 'gpu')
							.map((gpuRelation) => gpuRelation.entity.value as GPURelation_GPU);
						return {
							type: 'machine',
							id: `machine${gpuRelation.entity.value.id}`,
							data: {
								machine: gpuRelation.entity.value,
								gpus: gpus.filter(
									(gpu) => gpu.machineId === (gpuRelation.entity.value as GPURelation_Machine).id
								)
							},
							position
						};
					} else if (gpuRelation.entity.case === 'gpu') {
						const pods = response.gpuRelations
							.filter((gpuRelation) => gpuRelation.entity.case === 'pod')
							.map((gpuRelation) => gpuRelation.entity.value as GPURelation_Pod);
						if (
							pods
								.flatMap((pod) => pod.devices.map((device) => device.gpuId))
								.includes(gpuRelation.entity.value.id)
						) {
							return {
								type: 'gpu',
								id: `gpu${gpuRelation.entity.value.id}`,
								data: {
									scope,
									gpu: gpuRelation.entity.value,
									devices: pods.flatMap((pod) =>
										pod.devices.filter((device) => {
											return device.gpuId === (gpuRelation.entity.value as GPURelation_GPU).id;
										})
									)
								},
								position
							};
						} else {
							return {} as Node;
						}
					} else if (gpuRelation.entity.case === 'pod') {
						return {
							type: 'model',
							id: `pod${gpuRelation.entity.value.namespace}${gpuRelation.entity.value.name}`,
							data: { scope, pod: gpuRelation.entity.value },
							position
						};
					} else {
						return {} as Node;
					}
				})
			);

			edges.set(
				response.gpuRelations.flatMap((gpuRelation) => {
					if (gpuRelation.entity.case === 'gpu') {
						const gpu = gpuRelation.entity.value as GPURelation_GPU;
						return [
							{
								type: 'edge',
								id: `gpu${gpu.id}machine${gpu.machineId}`,
								source: `gpu${gpu.id}`,
								target: `machine${gpu.machineId}`,
								animated: true,
								selectable: false
							}
						];
					} else if (gpuRelation.entity.case === 'pod') {
						return gpuRelation.entity.value.devices.map((device) => {
							const pod = gpuRelation.entity.value as GPURelation_Pod;
							return {
								type: 'edge',
								id: `pod${pod.namespace}${pod.name}gpu${device.gpuId}`,
								source: `pod${pod.namespace}${pod.name}`,
								target: `gpu${device.gpuId}`,
								animated: true,
								selectable: false
							};
						});
					} else {
						return [];
					}
				})
			);
		} catch (error) {
			console.error('Failed to fetch GPU relations:', error);
		}
	}

	let open = $state(false);
	let isLoaded = $state(false);
</script>

<Drawer.Root
	bind:open
	onOpenChange={async (isOpen) => {
		if (isOpen) {
			if (!isLoaded) {
				await fetch();
				isLoaded = true;
			}
		} else {
			isLoaded = false;
		}
	}}
>
	<Drawer.Trigger class={buttonVariants({ variant: 'ghost' })}>
		<Icon icon="ph:graph" />
	</Drawer.Trigger>
	<Drawer.Content class="h-[77vh]">
		{#if isLoaded}
			<Drawer.Header>
				<Drawer.Title>{m.details()}</Drawer.Title>
				<Drawer.Description>
					<p>{m.gpu_relation_description()}</p>
					<p>{m.gpu_relation_guide_description()}</p>
				</Drawer.Description>
			</Drawer.Header>
			{#if open}
				<ComplexFlow.Flow initialNodes={$nodes} initialEdges={$edges} />
			{/if}
		{/if}
	</Drawer.Content>
</Drawer.Root>
