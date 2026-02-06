<script lang="ts">
	import Plus from '@lucide/svelte/icons/plus';

	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { buttonVariants } from '$lib/components/ui/button/button.svelte';

	let createModalOpen = $state(false);
</script>

<AlertDialog.Root bind:open={createModalOpen}>
	<AlertDialog.Trigger disabled class={buttonVariants({ variant: 'outline' })}>
		<Plus class="opacity-60" size={16} />
	</AlertDialog.Trigger>
</AlertDialog.Root>
