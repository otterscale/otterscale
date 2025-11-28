<script lang="ts">
	import Icon from '@iconify/svelte';

	import { cn } from '$lib/utils';

	import type { ModeSource } from '../types';
	import SelectCloudModel from './util-select-cloud-model.svelte';
	import SelectLocalModel from './util-select-local-model.svelte';

	let {
		value = $bindable(),
		modelSource = $bindable(),
		scope,
		namespace,
		required,
		invalid = $bindable()
	}: {
		value: string;
		modelSource?: ModeSource;
		scope: string;
		namespace: string;
		required: boolean;
		invalid: boolean;
	} = $props();

	$effect(() => {
		invalid = required && !value;
	});
</script>

<div
	class={cn(
		'flex items-center gap-2 rounded-lg border px-3 py-3 text-sm shadow-sm hover:cursor-default',
		invalid ? 'text-destructive ring-1 ring-destructive' : ''
	)}
>
	<Icon icon="ph:robot" />
	{#if value}
		{value}
	{:else}
		<p class="text-xs text-destructive/60">Select Model from Model Artifacts or HuggingFace.</p>
	{/if}
</div>

<div class="ml-auto flex gap-1">
	<SelectLocalModel bind:modelSource bind:value {scope} {namespace} />
	<SelectCloudModel bind:modelSource bind:value />
</div>
