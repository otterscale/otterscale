<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { getContext } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import {
		Nexus,
		type CreateStorageClassRequest,
		type Facility_Info
	} from '$gen/api/nexus/v1/nexus_pb';
	import { toast } from 'svelte-sonner';
	import { Button } from '$lib/components/ui/button';
	import { writable } from 'svelte/store';
	import * as Select from '$lib/components/ui/select/index.js';
	import { onMount } from 'svelte';

	let {
		scopeUuid
	}: {
		scopeUuid: string;
	} = $props();

	const kubernetesesStore = writable<Facility_Info[]>([]);
	const kubernetesesLoading = writable(true);
	async function fetchKuberneteses(scopeUuid: string) {
		try {
			const response = await client.listKuberneteses({ scopeUuid: scopeUuid });
			kubernetesesStore.set(response.kuberneteses);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			kubernetesesLoading.set(false);
		}
	}

	const cephesStore = writable<Facility_Info[]>([]);
	const cephesLoading = writable(true);
	async function fetchCephes(scopeUuid: string) {
		try {
			const response = await client.listCephes({ scopeUuid: scopeUuid });
			cephesStore.set(response.cephes);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			cephesLoading.set(false);
		}
	}

	function getKey(facility_information: Facility_Info) {
		if (facility_information.scopeName && facility_information.facilityName) {
			return facility_information.scopeUuid + facility_information.facilityName;
		} else {
			return 'Select';
		}
	}
	function getIdentifier(facility_information: Facility_Info) {
		if (facility_information.scopeName && facility_information.facilityName) {
			return [facility_information.scopeName, facility_information.facilityName].join('/');
		} else {
			return 'Select';
		}
	}

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_FACILITY_INFORMATION = {} as {} as Facility_Info;
	const DEFAULT_REQUEST = {} as CreateStorageClassRequest;

	let createStorageClassRequest = $state(DEFAULT_REQUEST);
	let defaultKubernetes = $state(DEFAULT_FACILITY_INFORMATION);
	let defaultCeph = $state(DEFAULT_FACILITY_INFORMATION);

	function reset() {
		defaultKubernetes = DEFAULT_FACILITY_INFORMATION;
		defaultCeph = DEFAULT_FACILITY_INFORMATION;
		createStorageClassRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchKuberneteses(scopeUuid);
			await fetchCephes(scopeUuid);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Button class="flex items-center justify-between gap-2" variant="ghost">
			<Icon icon="ph:plus" /> Storage Classes
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Create Storage Classes</AlertDialog.Title>
			<AlertDialog.Description class="grid gap-2 rounded-lg border p-4">
				<span class="grid gap-2">
					<Label>Kubernetes</Label>
					<Select.Root type="single">
						<Select.Trigger>
							{getIdentifier(defaultKubernetes)}
						</Select.Trigger>
						<Select.Content>
							{#if $kubernetesesStore.length > 0}
								{#each $kubernetesesStore as kubernetes}
									{@const selection = getKey(kubernetes)}
									<Select.Item
										value={selection}
										onclick={() => {
											defaultKubernetes = kubernetes;
											createStorageClassRequest.kubernetes = kubernetes;
										}}
									>
										{getIdentifier(kubernetes)}
									</Select.Item>
								{/each}
							{:else}
								<p class="p-2 text-xs text-muted-foreground">No Kubernetes</p>
							{/if}
						</Select.Content>
					</Select.Root>
				</span>
				<span class="grid gap-2">
					<Label>Ceph</Label>
					<Select.Root type="single">
						<Select.Trigger>
							{getIdentifier(defaultCeph)}
						</Select.Trigger>
						<Select.Content>
							{#if $cephesStore.length > 0}
								{#each $cephesStore as ceph}
									{@const selection = getKey(ceph)}
									<Select.Item
										value={selection}
										onclick={() => {
											defaultCeph = ceph;
											createStorageClassRequest.ceph = ceph;
										}}
									>
										{getIdentifier(ceph)}
									</Select.Item>
								{/each}
							{:else}
								<p class="p-2 text-xs text-muted-foreground">No Ceph</p>
							{/if}
						</Select.Content>
					</Select.Root>
				</span>

				<span class="grid w-full items-center gap-2">
					<Label>Prefix</Label>
					<Input class="w-full" bind:value={createStorageClassRequest.prefix} />
				</span>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client
						.createStorageClass(createStorageClassRequest)
						.then((r) => {
							toast.info(`Create storage classes`);
						})
						.catch((e) => {
							toast.error(`Create storage classes fail`);
						});
					// toast.info(`Create Storage Classes`);
					console.log(createStorageClassRequest);
					reset();
					close();
				}}
			>
				Create
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
