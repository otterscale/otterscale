<script lang="ts">
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Code from '$lib/components/custom/code';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils';
	import 'highlight.js/styles/github.css';
	import Monaco from 'svelte-monaco';

	let {
		id,
		value = $bindable(),
		language,
		required,
		preview = true,
		invalid = $bindable()
	}: {
		id?: string;
		language: 'json';
		value: string;
		required?: boolean;
		preview?: boolean;
		invalid?: boolean | null | undefined;
	} = $props();

	let temporaryValue = $state(value);
	let open = $state(false);

	const isInvalid = $derived(required && !value);
	$effect(() => {
		invalid = isInvalid;
	});
</script>

{#if preview}
	{#key value}
		{#if value}
			<Code.Root class="w-full" lang={language} code={value} hideLines>
				<Code.CopyButton />
			</Code.Root>
		{/if}
	{/key}
{/if}

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger
		class={cn(
			buttonVariants({ variant: 'outline' }),
			'ring-1',
			isInvalid ? 'ring-destructive' : ''
		)}
	>
		{#if isInvalid}
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
			options={{
				language,
				automaticLayout: true,
				padding: { top: 8, bottom: 8 },
				overviewRulerBorder: false,
				hideCursorInOverviewRuler: true
			}}
			theme="vs-dark"
			bind:value={temporaryValue}
		/>
		<AlertDialog.Footer>
			<AlertDialog.Cancel
				onclick={() => {
					temporaryValue = value;
				}}
			>
				Cancel
			</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					value = temporaryValue;
					open = false;
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
