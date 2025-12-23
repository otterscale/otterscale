<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { resolve } from '$app/paths';
	import { ApplicationService, type Job } from '$lib/api/application/v1/application_pb';
	import { ModelService } from '$lib/api/model/v1/model_pb';
	import { getJobStatus } from '$lib/components/applications/jobs/utils';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let {
		modelName = $bindable(),
		persistentVolumeClaimName = $bindable(),
		fromPersistentVolumeClaim = $bindable(),
		scope,
		namespace
	}: {
		modelName: string;
		persistentVolumeClaimName: string;
		fromPersistentVolumeClaim: boolean;
		scope: string;
		namespace: string;
	} = $props();

	const transport: Transport = getContext('transport');

	const applicationClient = createClient(ApplicationService, transport);
	const modelClient = createClient(ModelService, transport);

	const modelArtifactOptions = writable<SingleSelect.OptionType[]>([]);
	async function fetchModelArtifactOptions() {
		const response = await modelClient.listModelArtifacts({
			scope: scope,
			namespace: namespace
		});
		modelArtifactOptions.set(
			response.modelArtifacts.map((modelArtifact) => ({
				value: modelArtifact.name,
				label: modelArtifact.modelName,
				icon: 'ph:robot'
			}))
		);
	}

	let jobs = $state([] as Job[]);
	async function fetchJobs() {
		const response = await applicationClient.listJobs({
			scope: scope,
			namespace: 'llm-d'
		});
		jobs = response.jobs;
	}
	const jobMap = $derived(new Map(jobs.map((job) => [job.name, job])));

	async function fetch() {
		try {
			await Promise.all([fetchModelArtifactOptions(), fetchJobs()]);
		} catch (error) {
			console.error('Failed to fetch model artifacts and namespaces:', error);
		}
	}

	let isLoaded = $state(false);
	onMount(async () => {
		try {
			await fetch();
			isLoaded = true;
		} catch (error) {
			console.debug('Failed to init data:', error);
		}
	});
</script>

{#if isLoaded}
	<Select.Root type="single">
		<Select.Trigger>
			<Icon icon="ph:archive-fill" />
		</Select.Trigger>
		<Select.Content>
			{#if $modelArtifactOptions.length > 0}
				{#each $modelArtifactOptions as option (option.value)}
					{@const downloadModelArtifactJob = jobMap.get(option.value)}
					{@const isAvailable =
						downloadModelArtifactJob && getJobStatus(downloadModelArtifactJob) === 'Complete'}
					<Select.Item
						value={option.value}
						onclick={() => {
							fromPersistentVolumeClaim = true;
							modelName = option.label;
							persistentVolumeClaimName = option.value;
						}}
						disabled={!isAvailable}
					>
						<div class="flex items-center gap-2">
							{#if isAvailable}
								<Icon
									icon={option.icon ? option.icon : 'ph:empty'}
									class={cn('size-5', option.icon ? 'visible' : 'invisible')}
								/>
							{:else}
								<Icon icon="ph:spinner-gap" class="size-5 animate-spin" />
							{/if}
							<div>
								<h4>{option.label}</h4>
								<p class="text-muted-foreground">{option.value}</p>
							</div>
						</div>
					</Select.Item>
				{/each}
			{:else}
				<Empty.Root>
					<Empty.Header>
						<Empty.Media variant="icon">
							<Icon icon="ph:robot" />
						</Empty.Media>
						<Empty.Title>{m.no_model_artifact()}</Empty.Title>
						<Empty.Description>
							{m.no_model_artifact_guide()}
						</Empty.Description>
					</Empty.Header>
					<Empty.Content>
						<Button
							href={resolve('/(auth)/scope/[scope]/settings/model-artifact', { scope: scope })}
						>
							{m.download()}
						</Button>
					</Empty.Content>
				</Empty.Root>
			{/if}
		</Select.Content>
	</Select.Root>
{/if}
