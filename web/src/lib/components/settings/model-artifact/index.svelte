<script lang="ts">
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { ApplicationService, type Job } from '$lib/api/application/v1/application_pb';
	import type { ModelArtifact } from '$lib/api/model/v1/model_pb';
	import { ModelService } from '$lib/api/model/v1/model_pb';
	import { getJobStatus } from '$lib/components/applications/jobs/utils';
	import { Reloader, ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table';
	import * as Layout from '$lib/components/settings/layout';
	import { Spinner } from '$lib/components/ui/spinner';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { activeNamespace } from '$lib/stores';

	import Create from './create.svelte';
	import Delete from './delete.svelte';

	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const modelClient = createClient(ModelService, transport);
	const applicationClient = createClient(ApplicationService, transport);

	const modelArtifacts = writable<ModelArtifact[]>([]);
	async function fetchModelArtifacts() {
		const response = await modelClient.listModelArtifacts({
			scope: scope
		});
		modelArtifacts.set(response.modelArtifacts);
	}

	let jobs = $state([] as Job[]);
	async function fetchJobs() {
		const llmdResponse = await applicationClient.listJobs({
			scope: scope,
			namespace: 'llm-d'
		});
		const activeNamespaceResponse = await applicationClient.listJobs({
			scope: scope,
			namespace: $activeNamespace
		});
		jobs = [...llmdResponse.jobs, ...activeNamespaceResponse.jobs];
	}
	const jobMap = $derived(new Map(jobs.map((job) => [job.name, job])));

	async function fetch() {
		try {
			await Promise.all([fetchModelArtifacts(), fetchJobs()]);
		} catch (error) {
			console.error('Failed to fetch model artifacts and namespaces:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	let isMounted = $state(false);
	onMount(async () => {
		await fetch();
		isMounted = true;
		reloadManager.start();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

{#if isMounted}
	<Layout.Root>
		<Layout.Title>{m.model_artifact()}</Layout.Title>
		<Layout.Description>
			{m.model_artifact_setting_description()}
		</Layout.Description>
		<Layout.Controller>
			<Create {scope} {reloadManager} />
			<Reloader
				checked={reloadManager.state}
				onCheckedChange={(isChecked) => {
					if (isChecked) {
						reloadManager.restart();
					} else {
						reloadManager.stop();
					}
				}}
			/>
		</Layout.Controller>
		<Layout.Viewer>
			<div class="w-full rounded-lg border shadow-sm">
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>{m.name()}</Table.Head>
							<Table.Head>{m.namespace()}</Table.Head>
							<Table.Head>{m.model_name()}</Table.Head>
							<Table.Head>{m.status()}</Table.Head>
							<Table.Head>{m.phase()}</Table.Head>
							<Table.Head class="text-right">{m.size()}</Table.Head>
							<Table.Head>{m.volume()}</Table.Head>
							<Table.Head class="text-end">{m.create_time()}</Table.Head>
							<Table.Head></Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each $modelArtifacts as modelArtifact, index (index)}
							<Table.Row>
								<Table.Cell>{modelArtifact.name}</Table.Cell>
								<Table.Cell>{modelArtifact.namespace}</Table.Cell>
								<Table.Cell>{modelArtifact.modelName}</Table.Cell>
								<Table.Cell>
									{@const downloadModelArtifactJob = jobMap.get(modelArtifact.name)}
									{#if downloadModelArtifactJob}
										{@const status = getJobStatus(downloadModelArtifactJob)}
										{#if status === 'Running'}
											<span class="flex items-center gap-1 text-muted-foreground">
												<Spinner />
												{status}
											</span>
										{:else if ['Failed', 'FailureTarget'].includes(status)}
											<p class="text-destructive">
												{status}
											</p>
										{:else}
											{status}
										{/if}
									{/if}
								</Table.Cell>
								<Table.Cell>{modelArtifact.phase}</Table.Cell>
								<Table.Cell class="text-right">
									{@const { value, unit } = formatCapacity(modelArtifact.size)}
									{value}
									{unit}
								</Table.Cell>
								<Table.Cell>
									<Tooltip.Provider>
										<Tooltip.Root>
											<Tooltip.Trigger>
												<p class="max-w-30 truncate">{modelArtifact.volumeName}</p>
											</Tooltip.Trigger>
											<Tooltip.Content>
												{modelArtifact.volumeName}
											</Tooltip.Content>
										</Tooltip.Root>
									</Tooltip.Provider>
								</Table.Cell>
								<Table.Cell>
									<div class="flex justify-end">
										{#if modelArtifact.createdAt}
											<Tooltip.Provider>
												<Tooltip.Root>
													<Tooltip.Trigger>
														{formatTimeAgo(timestampDate(modelArtifact.createdAt))}
													</Tooltip.Trigger>
													<Tooltip.Content>
														{timestampDate(modelArtifact.createdAt)}
													</Tooltip.Content>
												</Tooltip.Root>
											</Tooltip.Provider>
										{/if}
									</div>
								</Table.Cell>
								<Table.Cell>
									<Delete {modelArtifact} {scope} {reloadManager} />
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		</Layout.Viewer>
	</Layout.Root>
{/if}
