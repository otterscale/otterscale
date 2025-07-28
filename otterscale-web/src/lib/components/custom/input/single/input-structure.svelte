<!-- <script lang="ts" module>
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { FormValidator } from '$lib/components/custom/form';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils';
	import 'highlight.js/styles/github.css';
	import { getContext } from 'svelte';
	import Monaco from 'svelte-monaco';
	import { DialogStateController } from '../../utils.svelte';
</script>

<script lang="ts">
	let {
		id,
		value = $bindable(),
		language,
		required,
		preview
	}: {
		id?: string;
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

	let controller = $state(new DialogStateController(false));

	const isNotFilled = $derived(required && !value);

	const formValidator: FormValidator = getContext('FormValidator');
	$effect(() => {
		formValidator.set(id, isNotFilled);
	});
</script>

{#if preview}
	{#key value}
		{#if value}
			<div class="h-[100px] w-full">
				<Monaco
					options={{
						language,
						automaticLayout: true,

						padding: { top: 8, bottom: 8 },
						overviewRulerBorder: false,
						hideCursorInOverviewRuler: true,

						readOnly: true
					}}
					theme="vs-dark"
					bind:value
				/>
			</div>
		{/if}
	{/key}
{/if}

<AlertDialog.Root bind:open={controller.state}>
	<AlertDialog.Trigger
		class={cn(
			buttonVariants({ variant: 'outline' }),
			'ring-1',
			isNotFilled ? 'ring-destructive' : ''
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
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					value = temporaryValue;
					controller.close();
				}}>Confirm</AlertDialog.Action
			>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root> -->
