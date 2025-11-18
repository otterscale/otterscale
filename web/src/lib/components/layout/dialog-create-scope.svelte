<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';

	import InstallCard from './install-card.svelte';
	import { plans } from './plans';
	import SheetCreateScope from './sheet-create-scope.svelte';

	let { open = $bindable(false) }: { open: boolean } = $props();

	let openSheet = $state(false);
	let selected = $state(-1);

	function handlePlanSelect(index: number) {
		open = false;
		openSheet = true;
		selected = index;
	}
</script>

<SheetCreateScope bind:open={openSheet} plan={plans[selected]} />

<Dialog.Root bind:open>
	<Dialog.Content showCloseButton={false} class="min-w-4xl overflow-hidden border-0 p-0">
		<InstallCard onSelect={handlePlanSelect} />
	</Dialog.Content>
</Dialog.Root>
