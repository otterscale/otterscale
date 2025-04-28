<script lang="ts">
	import { getContext } from 'svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { toast } from 'svelte-sonner';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Button } from '$lib/components/ui/button';
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import {
		Nexus,
		type CreateKubernetesRequest,
		type Scope,
		type Machine
	} from '$gen/api/nexus/v1/nexus_pb';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const scopesStore = writable<Scope[]>([]);
	const scopesLoading = writable(true);
	async function fetchScopes() {
		try {
			const response = await client.listScopes({});
			scopesStore.set(response.scopes);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			scopesLoading.set(false);
		}
	}

	const machinesStore = writable<Machine[]>([]);
	const machinesLoading = writable(true);
	async function fetchMachines() {
		try {
			const response = await client.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			machinesLoading.set(false);
		}
	}

	const DEFAULT_REQUEST = {} as CreateKubernetesRequest;

	let createKubernetesRequest = $state(DEFAULT_REQUEST);
	let selectedScopeName = $state('');
	let selectedMachineFQDN = $state('');

	function reset() {
		createKubernetesRequest = DEFAULT_REQUEST;
		selectedScopeName = '';
		selectedMachineFQDN = '';
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchScopes();
			await fetchMachines();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="p-4">
		<Button variant="destructive">Install</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Create Kubernetes</AlertDialog.Title>
			<AlertDialog.Description class="flex flex-col gap-4">
				<div class="grid w-full items-center gap-2">
					<Label>Scope</Label>
					<Select.Root type="single" bind:value={createKubernetesRequest.scopeUuid}>
						<Select.Trigger class="w-full">
							{selectedScopeName || 'Select'}
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								{#each $scopesStore as scope}
									<Select.Item
										value={scope.uuid}
										onclick={() => {
											selectedScopeName = scope.name;
										}}
									>
										{scope.name}
									</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
				</div>
				<div class="grid w-full items-center gap-2">
					<Label>Machine</Label>
					<Select.Root type="single" bind:value={createKubernetesRequest.machineId}>
						<Select.Trigger class="w-full">
							{selectedMachineFQDN || 'Select'}
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								{#each $machinesStore as machine}
									<Select.Item
										value={machine.id}
										onclick={() => {
											selectedMachineFQDN = machine.fqdn;
										}}
									>
										{machine.fqdn}
									</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
				</div>
				<div class="grid w-full items-center gap-2">
					<Label>Prefix</Label>
					<Input class="w-full" bind:value={createKubernetesRequest.prefixName} />
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.createKubernetes(createKubernetesRequest).then((r) => {
						toast.info(`Create Kubernetes ${createKubernetesRequest.scopeUuid}`);
					});
					// toast.info(`Create Kubernetes ${createKubernetesRequest.scopeUuid}`);
					console.log(createKubernetesRequest);
					reset();
					close();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
