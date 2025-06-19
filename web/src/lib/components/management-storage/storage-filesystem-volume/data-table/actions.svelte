<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import ActionAuthorize from './action-authorize.svelte';
	import ActionDelete from './action-delete.svelte';
	import ActionEdit from './action-edit.svelte';
	import { ActionSnapshot } from './action-snapshot';
	import type { FileSystem } from './types';

	let { fileSystem }: { fileSystem: FileSystem } = $props();
</script>

<Layout.Actions>
	<Layout.ActionSeparator />
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionAuthorize {fileSystem} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionSnapshot />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionEdit {fileSystem} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {fileSystem} />
	</Layout.ActionItem>
</Layout.Actions>
