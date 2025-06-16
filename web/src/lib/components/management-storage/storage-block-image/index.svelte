<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import { ImageDataTable } from './image-data-table';
	import { NamespaceDataTable } from './namespace-data-table';
</script>

<main class="space-y-4">
	<ImageDataTable />

	<Label class="bg-muted flex w-full justify-center rounded-lg p-4 text-xl">Namespace</Label>
	<NamespaceDataTable />
</main>
