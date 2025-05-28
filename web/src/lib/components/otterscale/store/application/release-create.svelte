<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Input } from '$lib/components/ui/input';
	import { toast } from 'svelte-sonner';
	import { Label } from '$lib/components/ui/label';
	import {
		Nexus,
		type Application_Chart,
		type Application_Release,
		type CreateReleaseRequest,
		type Facility_Info
	} from '$gen/api/nexus/v1/nexus_pb';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { writable } from 'svelte/store';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { ReleaseValuesEdit } from '$lib/components/otterscale/index';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { valuesMapList } from './dataset';

	let {
		releases = $bindable(),
		chart
	}: {
		releases: Application_Release[];
		chart: Application_Chart;
	} = $props();

	const transport: Transport = getContext('transport');
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
			const scopes = await client.listScopes({});
			const kubernetesesPromises = scopes.scopes.map((scope) =>
				client.listKuberneteses({ scopeUuid: scope.uuid })
			);
			const responses = await Promise.all(kubernetesesPromises);
			const allKuberneteses = responses.flatMap((response) => response.kuberneteses);
			kubernetesesStore.set(allKuberneteses);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			kubernetesesLoading.set(false);
		}
	}

	const DEFAULT_VERSION = chart.versions[0];
	const DEFAULT_KUBERNETES = {} as Facility_Info;
	const DEFAULT_REQUEST = {
		name: '',
		valuesYaml: '',
		valuesMap: valuesMapList[chart.name],
		dryRun: false
	} as CreateReleaseRequest;

	let createReleaseRequest = $state(DEFAULT_REQUEST);
	let selectedVersion = $state(DEFAULT_VERSION);
	let selectedKubernetes = $state(DEFAULT_KUBERNETES);

	function integrate() {
		createReleaseRequest.chartRef = selectedVersion.chartRef;
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

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchChart(chart.name);
			await fetchKuberneteses();

			isMounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>
		{#if isMounted}
			Install
		{:else}
			<ComponentLoading />
		{/if}
	</AlertDialog.Trigger>
	<AlertDialog.Content interactOutsideBehavior="close">
		<AlertDialog.Header>
			<AlertDialog.Description class="space-y-2">
				<fieldset class="items-center rounded-lg border p-4">
					<legend class="text-sm">Basic</legend>
					<div class="grid gap-2">
						<span class="flex items-center justify-between gap-2">
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

						<span class="flex items-center justify-between gap-2">
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
													createReleaseRequest.scopeUuid = kubernetes.scopeUuid;
													createReleaseRequest.facilityName = kubernetes.facilityName;
												}}
												value={`${kubernetes.scopeName}/${kubernetes.facilityName}`}
											>
												{`${kubernetes.scopeName}/${kubernetes.facilityName}`}
											</Select.Item>
										{/each}
									</Select.Group>
								</Select.Content>
							</Select.Root>
						</span>
						<!-- <span>
							<Label>Scope</Label>
							<Input bind:value={createReleaseRequest.scopeUuid} />
						</span>
						<span>
							<Label>Facility</Label>
							<Input bind:value={createReleaseRequest.facilityName} />
						</span> -->
						<span class="flex items-center justify-between gap-2">
							<Label>Namespace</Label>
							<Input bind:value={createReleaseRequest.namespace} />
						</span>

						<span class="flex items-center justify-between gap-2">
							<Label>Dry Run</Label>
							<Switch id="enable_ssh" bind:checked={createReleaseRequest.dryRun} />
						</span>
					</div>
				</fieldset>
				<fieldset class="items-center rounded-lg border p-3">
					<legend class="text-sm">Advance</legend>

					<span class="flex items-center justify-between">
						<Label>Configuration</Label>
						<ReleaseValuesEdit
							chartRef={selectedVersion.chartRef}
							bind:valuesYaml={createReleaseRequest.valuesYaml}
							bind:valuesMap={createReleaseRequest.valuesMap}
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

					toast.promise(() => client.createRelease(createReleaseRequest), {
						loading: 'Loading...',
						success: (r) => {
							client.listReleases({}).then((r) => {
								releases = r.releases;
							});
							return `Create ${r.name} success`;
						},
						error: (e) => {
							let msg = `Fail to create ${createReleaseRequest.name}`;
							toast.error(msg, {
								description: (e as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY
							});
							return msg;
						}
					});

					reset();
					close();
				}}>Confirm</AlertDialog.Action
			>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
