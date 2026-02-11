<script lang="ts">
	import { Eye } from '@lucide/svelte';
	import { stringify } from 'yaml';

	import * as Code from '$lib/components/custom/code';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Item from '$lib/components/ui/item';

	let {
		schema: apiSchema,
		object,
		onOpenChangeComplete
	}: {
		schema: any;
		object: Record<string, unknown>;
		onOpenChangeComplete?: () => void;
	} = $props();

	let open = $state(false);
</script>

<Dialog.Root bind:open {onOpenChangeComplete}>
	<Dialog.Trigger class="w-full">
		<Item.Root class="p-0 text-xs" size="sm">
			<Item.Media>
				<Eye />
			</Item.Media>
			<Item.Content>
				<Item.Title>View</Item.Title>
			</Item.Content>
		</Item.Root>
	</Dialog.Trigger>
	<Dialog.Content
		class="flex h-fit max-h-[77vh] max-w-[62vw] min-w-[50vw] flex-col justify-between"
	>
		<Dialog.Header>Configuration</Dialog.Header>
		<Dialog.Description>{apiSchema?.description}</Dialog.Description>
		<Code.Root
			variant="secondary"
			lang="yaml"
			code={stringify(object)}
			class="w-full border-none"
		/>
	</Dialog.Content>
</Dialog.Root>
