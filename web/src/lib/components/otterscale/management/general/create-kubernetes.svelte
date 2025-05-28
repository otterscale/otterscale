NEW BOOTSTRAP FLOW

<!-- <script lang="ts">
	import { getContext } from 'svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import Icon from '@iconify/svelte';
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
		type Machine,
		type CreateCephRequest,
		type Facility_Info,
		type SetCephCSIRequest
	} from '$gen/api/nexus/v1/nexus_pb';

	const transport: Transport = getContext('transport');
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

	async function fetchMachines(scopeUuid: string) {
		try {
			const response = await client.listMachines({ scopeUuid: scopeUuid });
			return response.machines;
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	async function fetchOSDDevices(machineId: string) {
		try {
			const response = await client.getMachine({ id: machineId });
			return response.blockDevices
				.filter((device) => !device.bootDisk)
				.map((device) => `/dev/${device.name}`);
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	const DEFAULT_CREATE_CEPH_REQUEST = {
		osdDevices: [] as string[],
		development: true
	} as CreateCephRequest;
	const DEFAULT_CREATE_KUBERNETES_REQUEST = {
		virtualIps: [] as string[]
	} as CreateKubernetesRequest;

	const DEFAULT_SCOPE = {} as Scope;
	const DEFAULT_NACHINE = {} as Machine;
	const DEFAULT_PREFIX = '';

	let createCephRequest = $state(DEFAULT_CREATE_CEPH_REQUEST);
	let createKubernetesRequest = $state(DEFAULT_CREATE_KUBERNETES_REQUEST);
	let selectedScope = $state(DEFAULT_SCOPE);
	let selectedMachine = $state(DEFAULT_NACHINE);
	let inputedPrefix = $state(DEFAULT_PREFIX);

	function integrate() {
		createCephRequest.scopeUuid = selectedScope.uuid;
		createCephRequest.machineId = selectedMachine.id;
		createCephRequest.prefixName = inputedPrefix;

		createKubernetesRequest.scopeUuid = selectedScope.uuid;
		createKubernetesRequest.machineId = selectedMachine.id;
		createKubernetesRequest.prefixName = inputedPrefix;
	}

	function getSetCephCSIRequest(
		ceph: Facility_Info,
		kubernetes: Facility_Info,
		development: boolean
	) {
		return {
			ceph: ceph,
			kubernetes: kubernetes,
			prefix: inputedPrefix,
			development: development
		} as SetCephCSIRequest;
	}

	function reset() {
		inputedPrefix = '';
		createKubernetesRequest = DEFAULT_CREATE_KUBERNETES_REQUEST;
		createCephRequest = DEFAULT_CREATE_CEPH_REQUEST;
		resetSeletedScope();
		resetSeletedMachine();
		resetSeletedOSDDevices();
	}
	function resetSeletedScope() {
		selectedScope = DEFAULT_SCOPE;
	}
	function resetSeletedMachine() {
		selectedMachine = DEFAULT_NACHINE;
	}
	function resetSeletedOSDDevices() {
		createCephRequest.osdDevices = [] as string[];
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchScopes();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>
		<Button>Deploy</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Create Kubernetes</AlertDialog.Title>
			<AlertDialog.Description class="flex flex-col gap-4">
				<div class="grid w-full items-center gap-2">
					<Label>Prefix</Label>
					<Input class="w-full" bind:value={inputedPrefix} />
				</div>

				<div class="grid w-full items-center gap-2">
					<Label>Scope</Label>
					<Select.Root type="single">
						<Select.Trigger class="w-full">
							{selectedScope.name || 'Select'}
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								{#each $scopesStore as scope}
									<Select.Item
										value={scope.uuid}
										onclick={() => {
											selectedScope = scope;
											resetSeletedMachine();
											resetSeletedOSDDevices();
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
					{#if !selectedScope.uuid}
						<p class="w-full rounded-lg bg-muted/30 p-2 text-center text-muted-foreground">
							Select a scope
						</p>
					{:else}
						{#await fetchMachines(selectedScope.uuid)}
							Loading
						{:then machines}
							{#if machines && machines.length == 0}
								<p class="w-full rounded-lg bg-muted/30 p-2 text-center text-muted-foreground">
									There is no machine
								</p>
							{:else if machines && machines.length > 0}
								<Select.Root type="single">
									<Select.Trigger class="w-full">
										{selectedMachine.fqdn || 'Select'}
									</Select.Trigger>
									<Select.Content>
										<Select.Group>
											{#each machines as machine}
												<Select.Item
													value={machine.id}
													onclick={() => {
														selectedMachine = machine;
														resetSeletedOSDDevices();
													}}
												>
													{machine.fqdn}
												</Select.Item>
											{/each}
										</Select.Group>
									</Select.Content>
								</Select.Root>
							{/if}
						{/await}
					{/if}
				</div>

				<div class="grid w-full items-center gap-2">
					<Label>OSD Devices</Label>

					<span class="flex flex-wrap items-center gap-2">
						{#each createCephRequest.osdDevices as device}
							<Badge variant="secondary" class="flex w-fit gap-1 text-sm hover:cursor-pointer">
								{device}
							</Badge>
						{/each}
						{#if !selectedMachine.id}
							<p class="w-full rounded-lg bg-muted/30 p-2 text-center text-muted-foreground">
								Select a machine
							</p>
						{:else}
							{#await fetchOSDDevices(selectedMachine.id)}
								Loading...
							{:then osdDevices}
								{#if osdDevices && osdDevices.length == 0}
									<p class="w-full rounded-lg bg-muted/30 p-2 text-center text-muted-foreground">
										There is no device
									</p>
								{:else if osdDevices && osdDevices.length > 0}
									<Select.Root type="multiple" bind:value={createCephRequest.osdDevices}>
										<Select.Trigger>Select</Select.Trigger>
										<Select.Content>
											{#each osdDevices as device}
												<Select.Item value={device}>
													{device}
												</Select.Item>
											{/each}
										</Select.Content>
									</Select.Root>
								{/if}
							{/await}
						{/if}
					</span>
				</div>

				<div class="flex items-center justify-between gap-2">
					<Label>Development</Label>
					<Switch bind:checked={createCephRequest.development} />
				</div>
				<div class="grid gap-2 rounded-lg bg-muted/50 p-4">
					<p class="text-sm text-muted-foreground">Single Node mode when development is true.</p>
					<p class="text-xs font-light">
						Note that High Availability (HA) is not available in single node mode since Ceph-mon is
						configured for a single node and CephFS functionality is not implemented.
					</p>
				</div>

				<div class="grid w-full items-center gap-2">
					<Label>Virtual IPs</Label>
					<span class="flex flex-wrap items-center gap-2">
						{#each createKubernetesRequest.virtualIps as ip}
							<Badge
								variant="secondary"
								class="flex gap-1 text-sm hover:cursor-pointer"
								onclick={() => {
									createKubernetesRequest.virtualIps = createKubernetesRequest.virtualIps.filter(
										(_, i) => i !== createKubernetesRequest.virtualIps.indexOf(ip)
									);
								}}
							>
								{ip}
								<Icon icon="ph:x" class="h-3 w-3" />
							</Badge>
						{/each}
					</span>
					<div class="flex w-full items-center justify-between gap-2">
						<Input
							onkeydown={(e) => {
								if (e.key === 'Enter') {
									createKubernetesRequest.virtualIps = [
										...createKubernetesRequest.virtualIps,
										e.currentTarget.value
									];
									e.currentTarget.value = '';
								}
							}}
							onblur={(e) => {
								createKubernetesRequest.virtualIps = [
									...createKubernetesRequest.virtualIps,
									e.currentTarget.value
								];
								e.currentTarget.value = '';
							}}
							placeholder="Press Enter to Add Virtual IP"
						/>
					</div>
				</div>

				<div class="grid w-full items-center gap-2">
					<Label>Calico CIDR</Label>
					<Input class="w-full" bind:value={createKubernetesRequest.calicoCidr} />
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					toast.loading('Loading...');
					integrate();

					client
						.createCeph(createCephRequest)
						.then(async (cr) => {
							await client
								.createKubernetes(createKubernetesRequest)
								.then(async (kr) => {
									await client
										.setCephCSI(getSetCephCSIRequest(cr, kr, createCephRequest.development))
										.then(() => {
											toast.success(`Create '${kr.facilityName}' & '${cr.facilityName}' success`);
										})
										.catch((e) => {
											toast.error(`Fail to create storage classes: ${e.toString()}`);
										});
								})
								.catch((e) => {
									toast.error(`Fail to create Kubernetes: ${e.toString()}`);
								});
						})
						.catch((e) => {
							toast.error(`Fail to create Ceph: ${e.toString()}`);
						})
						.finally(() => {
							reset();
						});

					close();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root> -->
