<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import Icon from '@iconify/svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import { toast } from 'svelte-sonner';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import {
		Nexus,
		type Application_Chart,
		type Application_Chart_Version,
		type CreateReleaseRequest,
		type Facility_Info
	} from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { writable } from 'svelte/store';

	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { ReleaseValuesEdit } from '$lib/components/otterscale/index';

	let {
		chart
	}: {
		chart: Application_Chart;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const chartStore = writable<Application_Chart>();
	const chartLoading = writable(true);
	async function fetchChart(chartName: string) {
		try {
			const response = await client.getChart({
				name: chartName
			});
			chartStore.set(response);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			chartLoading.set(false);
		}
	}

	const kubernetesesStore = writable<Facility_Info[]>([]);
	const kubernetesesLoading = writable(true);
	async function fetchKuberneteses() {
		try {
			const response = await client.listKuberneteses({});
			kubernetesesStore.set(response.kuberneteses);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			kubernetesesLoading.set(false);
		}
	}

	const DEFAULT_VERSION = chart.versions[0];
	const DEFAULT_KUBERNETES = {} as Facility_Info;
	const DEFAULT_REQUEST = { name: chart.name, dryRun: false } as CreateReleaseRequest;

	let createReleaseRequest = $state(DEFAULT_REQUEST);
	let selectedVersion = $state(DEFAULT_VERSION);
	let selectedKubernetes = $state(DEFAULT_KUBERNETES);

	function integrate() {
		createReleaseRequest.chartRef = selectedVersion.chartRef;
		// createReleaseRequest.scopeUuid = selectedKubernetes.scopeUuid;
		// createReleaseRequest.facilityName = selectedKubernetes.facilityName;
	}

	function reset() {
		createReleaseRequest = DEFAULT_REQUEST;
		selectedVersion = DEFAULT_VERSION;
		selectedKubernetes = DEFAULT_KUBERNETES;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	onMount(async () => {
		try {
			await fetchChart(chart.name);
			await fetchKuberneteses();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>Install</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Description>
				<fieldset class="items-center rounded-lg border p-4">
					<legend class="text-lg">Basic</legend>
					<div class="grid gap-2">
						<span class="flex items-center justify-between">
							<Label>Version</Label>
							<Select.Root type="single">
								<Select.Trigger class="w-fit">
									{selectedVersion.chartVersion ? selectedVersion.chartVersion : 'Select'}
								</Select.Trigger>
								<Select.Content>
									<Select.Group>
										{#each $chartStore.versions as version}
											<Select.Item
												onclick={() => {
													selectedVersion = version;
												}}
												value={version.chartVersion}
											>
												<span class="flex w-full items-end justify-between gap-1">
													{version.chartVersion}
													<Badge variant="outline" class="text-xs font-light"
														>{version.applicationVersion}</Badge
													>
												</span>
											</Select.Item>
										{/each}
									</Select.Group>
								</Select.Content>
							</Select.Root>
						</span>

						<span class="flex items-center justify-between">
							<Label>Kubernetes</Label>
							<Select.Root type="single">
								<Select.Trigger class="w-fit">
									{selectedKubernetes.scopeUuid && selectedKubernetes.facilityName
										? `${selectedKubernetes.scopeName}/${selectedKubernetes.facilityName}`
										: 'Select'}
								</Select.Trigger>
								<Select.Content>
									<Select.Group>
										{#each $kubernetesesStore as kubernetes}
											<Select.Item
												onclick={() => {
													selectedKubernetes = kubernetes;
												}}
												value={`${selectedKubernetes.scopeName}/${selectedKubernetes.facilityName}`}
											>
												`${selectedKubernetes.scopeName}/${selectedKubernetes.facilityName}`
											</Select.Item>
										{/each}
									</Select.Group>
								</Select.Content>
							</Select.Root>
						</span>
						<span>
							<Label>Scope</Label>
							<Input bind:value={createReleaseRequest.scopeUuid} />
						</span>
						<span>
							<Label>Facility</Label>
							<Input bind:value={createReleaseRequest.facilityName} />
						</span>
						<span>
							<Label>Namespace</Label>
							<Input bind:value={createReleaseRequest.namespace} />
						</span>

						<span class="flex items-center justify-between">
							<Label>Dry Run</Label>
							<Switch id="enable_ssh" bind:checked={createReleaseRequest.dryRun} />
						</span>
					</div>
				</fieldset>
				<fieldset class="items-center rounded-lg border p-3">
					<legend class="text-lg">Advance</legend>

					<span class="flex items-center justify-between">
						<Label>Configuration</Label>
						<ReleaseValuesEdit
							chartRef={selectedVersion.chartRef}
							bind:valuesYaml={createReleaseRequest.valuesYaml}
						/>
					</span>
				</fieldset>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					integrate();
					console.log(createReleaseRequest);
					// client.createRelease(createReleaseRequest).then((r) => {
					// 	toast.info(`Create ${r.name}.`);
					// });
					reset();
					close();
				}}>Confirm</AlertDialog.Action
			>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
