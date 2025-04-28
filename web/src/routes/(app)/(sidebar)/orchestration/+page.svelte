<script lang="ts">
	import '@xyflow/svelte/dist/style.css';
	import { type Edge, type Node } from '@xyflow/svelte';
	import Icon from '@iconify/svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import {
		StackService,
		type Application as Cluster,
		type Machine,
		type Model
	} from '$gen/api/stack/v1/stack_pb';
	import { KubeService, type Application } from '$gen/api/kube/v1/kube_pb';
	import { OrchestrationFlow } from '$lib/components/otterscale';

	const transport: Transport = getContext('transport');

	const nodes: Node[] = [];
	const edges: Edge[] = [];

	const stackClient = createClient(StackService, transport);
	const machinesStore = writable<Machine[]>([]);
	const machinesIsLoading = writable(true);
	async function fetchMachines() {
		try {
			const response = await stackClient.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching machines:', error);
		} finally {
			machinesIsLoading.set(false);
		}
	}
	const modelsStore = writable<Model[]>([]);
	const modelsIsLoading = writable(true);
	async function fetchModels() {
		try {
			const response = await stackClient.listModels({});
			modelsStore.set(response.models);
		} catch (error) {
			console.error('Error fetching models:', error);
		} finally {
			modelsIsLoading.set(false);
		}
	}
	const clustersStore = writable<Cluster[]>([]);
	const clustersIsLoading = writable(true);
	async function fetchClusters(modelUuid: string) {
		try {
			const response = await stackClient.listApplications({
				modelUuid: modelUuid
			});
			clustersStore.set(response.applications);
		} catch (error) {
			console.error('Error fetching applications:', error);
		} finally {
			clustersIsLoading.set(false);
		}
	}

	const kubeClient = createClient(KubeService, transport);
	const applicationsStore = writable<Application[]>([]);
	const applicationsIsLoading = writable(true);
	async function fetchApplications(modelUuid: string, clusterName: string) {
		try {
			const response = await kubeClient.listApplications({
				modelUuid: modelUuid,
				clusterName: clusterName
			});
			applicationsStore.set(response.applications);
		} catch (error) {
			console.error('Error fetching models:', error);
		} finally {
			applicationsIsLoading.set(false);
		}
	}

	let mounted = false;
	onMount(async () => {
		try {
			const modelClusterPairs: { modelUuid: string; clusterName: string; charmName: string }[] = [];

			await fetchMachines();
			await fetchModels();

			$machinesStore.flat().forEach((machine) => {
				nodes.push({
					type: 'MAAS',
					id: machine.systemId,
					data: { ...machine },
					position: { x: 0, y: 0 }
				});
			});
			$modelsStore.flat().forEach((model) => {
				nodes.push({
					type: 'JUJU',
					id: model.uuid,
					data: { ...model },
					position: { x: 0, y: 0 }
				});
			});

			for (const model of $modelsStore) {
				await fetchClusters(model.uuid);
				$clustersStore.flat().forEach((cluster) => {
					modelClusterPairs.push({
						modelUuid: model.uuid,
						clusterName: cluster.name,
						charmName: cluster.charmName
					});
				});
			}

			for (const modelClusterPair of modelClusterPairs) {
				if (modelClusterPair.charmName.toLowerCase().includes('kubernetes-worker')) {
					console.log(modelClusterPair.charmName);
					await fetchApplications(modelClusterPair.modelUuid, modelClusterPair.clusterName);
					$applicationsStore.flat().forEach((application) => {
						nodes.push({
							type: 'Kubernetes',
							id: `${application.namespace}/${application.name}`,
							data: { ...application },
							position: { x: 0, y: 0 }
						});
					});
				}
			}

			$machinesStore.forEach((machine) => {
				const jujuModelUuid = machine.workloadAnnotations['juju-model-uuid'];
				if (jujuModelUuid) {
					edges.push({
						id: `${machine.systemId}-${jujuModelUuid}`,
						source: jujuModelUuid,
						target: machine.systemId
					});
				}
			});

			for (const modelClusterPair of modelClusterPairs) {
				if (modelClusterPair.charmName.toLowerCase().includes('kubernetes-worker')) {
					$applicationsStore.forEach((application) => {
						edges.push({
							id: `${application.namespace}/${application.name}-${modelClusterPair.modelUuid}`,
							source: `${application.namespace}/${application.name}`,
							target: modelClusterPair.modelUuid
						});
					});
				}
			}
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<main class="h-[calc(100vh_-_theme(spacing.16))]">
	{#if mounted}
		<OrchestrationFlow {nodes} {edges} horizontal />
	{:else}
		<div class="flex h-full w-full items-center justify-center gap-2 text-sm text-muted-foreground">
			<Icon icon="ph:spinner" class="size-8 animate-spin" />
			Loading...
		</div>
	{/if}
</main>
