<script lang="ts" module>
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils';
	import 'highlight.js/styles/github.css';
	import Monaco from 'svelte-monaco';
	import { RING_INVALID_INPUT_CLASSNAME, RING_VALID_INPUT_CLASSNAME } from './utils.svelte';
</script>

<script lang="ts">
	let {
		value = $bindable(),
		language,
		required,
		preview
	}: {
		language: 'json';
		value: string;
		required?: boolean;
		preview?: boolean;
	} = $props();

	const ORIGIN_VALUE = value;
	let temporaryValue = $state(ORIGIN_VALUE);
	function reset() {
		temporaryValue = ORIGIN_VALUE;
	}

	let open = $state(false);

	const isNotFilled = $derived(required && !value);
</script>

{#if preview}
	{#key value}
		{#if value}
			<div class="h-[100px] w-full">
				<Monaco
					options={{
						language,
						automaticLayout: true,
						readOnly: true,
						padding: { top: 8, bottom: 8 },
						overviewRulerLanes: 0,
						overviewRulerBorder: false,
						hideCursorInOverviewRuler: true
					}}
					theme="vs-dark"
					bind:value
				/>
			</div>
		{/if}
	{/key}
{/if}

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger
		class={cn(
			buttonVariants({ variant: 'outline' }),
			isNotFilled ? RING_INVALID_INPUT_CLASSNAME : RING_VALID_INPUT_CLASSNAME
		)}
	>
		{#if isNotFilled}
			<p class={cn('text-destructive text-xs')}>Required</p>
		{:else}
			Input/Edit
		{/if}
	</AlertDialog.Trigger>
	<AlertDialog.Content class="h-[50vh] w-[50vw]">
		<AlertDialog.Header>
			<p class="w-full text-center text-xl font-bold">Editor</p>
		</AlertDialog.Header>
		<Monaco
			options={{ language, automaticLayout: true }}
			theme="vs-dark"
			bind:value={temporaryValue}
		/>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					value = temporaryValue;
					open = false;
				}}>Confirm</AlertDialog.Action
			>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
