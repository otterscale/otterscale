<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { ModelService } from '$lib/api/model/v1/model_pb';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import * as Select from '$lib/components/ui/select/index.js';
	import { cn } from '$lib/utils';

	import type { ModeSource } from '../types';
</script>

<script lang="ts">
	let {
		value = $bindable(),
		modelSource = $bindable(),
		scope,
		namespace
	}: { value: string; modelSource?: ModeSource; scope: string; namespace: string } = $props();

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
		<Select.Trigger
			onclick={() => {
				modelSource = 'local' as ModeSource;
			}}
		>
			<div class="font-base flex items-center gap-2 text-sm text-primary">
				<Icon icon="ph:archive-fill" />
				local
			</div>
		</Select.Trigger>
		<Select.Content>
			{#each $modelArtifactOptions as option (option.value)}
				<Select.Item value={option.value}>
					<Icon
						icon={option.icon ? option.icon : 'ph:empty'}
						class={cn('size-5', option.icon ? 'visible' : 'invisible')}
					/>
					{option.label}
				</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>
{/if}
