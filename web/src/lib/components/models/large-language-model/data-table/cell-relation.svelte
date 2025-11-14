<script lang="ts" module>
	import '@xyflow/svelte/dist/style.css';
	import '@xyflow/svelte/dist/style.css';

	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { type Edge, type Node } from '@xyflow/svelte';
	import { getContext } from 'svelte';
	import { type Writable, writable } from 'svelte/store';

	import {
		type GPURelation_GPU,
		type GPURelation_Machine,
		type GPURelation_Pod,
		OrchestratorService
	} from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { Complex as Flow } from '$lib/components/flow/index';
	import * as Sheet from '$lib/components/ui/sheet';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';

	import type { LargeLanguageModel } from '../type';
</script>

<script lang="ts">
	let { scope, model }: { scope: string; model: LargeLanguageModel } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(OrchestratorService, transport);

	const position = { x: 0, y: 0 };

	const nodes: Writable<Node[]> = writable([]);
	const edges: Writable<Edge[]> = writable([]);
	let isLoading = $state(true);

	client
		.listGPURelationsByModel({
			scope: scope,
			,
			namespace: model.application.namespace,
			modelName: model.name
		})
		.then((response) => {
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
								scope,
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
					} else if (gpuRelation.entity.case === 'pod') {
						return {
							type: 'model',
							id: `pod${gpuRelation.entity.value.namespace}${gpuRelation.entity.value.name}`,
							data: gpuRelation.entity.value,
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

			isLoading = false;
		});

	let open = $state(false);
</script>

{#if isLoading}
	Loading...
{:else}
	<Sheet.Root bind:open>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content side="right" class="min-w-[38vw] p-4">
			{#if open}
				<Sheet.Header>
					<Sheet.Title class="text-center text-lg">{m.details()}</Sheet.Title>
				</Sheet.Header>
				<Flow.Flow initialNodes={$nodes} initialEdges={$edges} />
			{/if}
		</Sheet.Content>
	</Sheet.Root>
{/if}
