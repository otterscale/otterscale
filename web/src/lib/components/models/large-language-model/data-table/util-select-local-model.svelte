<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { resolve } from '$app/paths';
	import { ModelService } from '$lib/api/model/v1/model_pb';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		value = $bindable(),
		scope,
		namespace,
		fromLocal = $bindable()
	}: { value: string; scope: string; namespace: string, fromLocal: boolean } = $props();

	const transport: Transport = getContext('transport');

	const modelClient = createClient(ModelService, transport);

	let isModelArtifactOptionsLoaded = $state(false);
	const modelArtifactOptions = writable<SingleSelect.OptionType[]>([]);
	async function fetchModelArtifactOptions() {
		const response = await modelClient.listModelArtifacts({
			scope: scope,
			namespace: namespace
		});
		modelArtifactOptions.set(
			response.modelArtifacts.map((modelArtifact) => ({
				value: modelArtifact.name,
				label: modelArtifact.name,
				icon: 'ph:robot'
			}))
		);
	}

	onMount(async () => {
		try {
			await fetchModelArtifactOptions();
			isModelArtifactOptionsLoaded = true;
		} catch (error) {
			console.debug('Failed to init data:', error);
		}
	});
</script>

{#if isModelArtifactOptionsLoaded}
	<Select.Root type="single" bind:value>
		<Select.Trigger>
			<Icon icon="ph:archive-fill" />
		</Select.Trigger>
		<Select.Content>
			{#each $modelArtifactOptions as option (option.value)}
				<Select.Item value={option.value} onclick={() => {
					fromLocal = true
				}}>
					<Icon
						icon={option.icon ? option.icon : 'ph:empty'}
						class={cn('size-5', option.icon ? 'visible' : 'invisible')}
					/>
					{option.label}
				</Select.Item>
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
			{/each}
		</Select.Content>
	</Select.Root>
{/if}
