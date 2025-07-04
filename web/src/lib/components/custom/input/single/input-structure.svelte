<script lang="ts" module>
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import 'highlight.js/styles/github.css';
	import Monaco from 'svelte-monaco';
</script>

<script lang="ts">
	let {
		value = $bindable()
	}: {
		value: string;
	} = $props();

	const ORIGIN_VALUE = value;
	let temporaryValue = $state(ORIGIN_VALUE);
	function reset() {
		temporaryValue = ORIGIN_VALUE;
	}

	let open = $state(false);
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class={buttonVariants({ variant: 'outline' })}>
		Input/Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content class="h-[50vh] w-[50vw]">
		<AlertDialog.Header>
			<p class="w-full text-center text-xl font-bold">Editor</p>
		</AlertDialog.Header>
		<Monaco
			options={{ language: 'json', automaticLayout: true }}
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
