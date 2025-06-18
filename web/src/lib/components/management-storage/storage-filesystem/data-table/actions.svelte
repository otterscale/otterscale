<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import ActionAuthorize from './action-authorize.svelte';
	import ActionDelete from './action-delete.svelte';
	import ActionEdit from './action-edit.svelte';
	import { ActionSnapshot } from './action-snapshot';
	import { Subvolumes } from './subvolumes';
	import { SubvolumeGroups } from './subvolume-groups';
	import type { FileSystem } from './types';

	let { fileSystem }: { fileSystem: FileSystem } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Views</Layout.ActionLabel>
		<Layout.ActionItem>
		<SubvolumeGroups />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Subvolumes />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionSnapshot />
	</Layout.ActionItem>
	<Layout.ActionSeparator />
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionEdit {fileSystem} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionAuthorize {fileSystem} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {fileSystem} />
	</Layout.ActionItem>
</Layout.Actions>
