<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';

	let { open = $bindable(false) }: { open: boolean } = $props();
</script>

<Sheet.Root bind:open>
	<Sheet.Content side="right">
		<Sheet.Header>
			<Sheet.Title>Title</Sheet.Title>
			<Sheet.Description>Description</Sheet.Description>
		</Sheet.Header>
		<Sheet.Footer>
			<Sheet.Close class={buttonVariants({ variant: 'outline' })}>Save changes</Sheet.Close>
		</Sheet.Footer>
	</Sheet.Content>
</Sheet.Root>
