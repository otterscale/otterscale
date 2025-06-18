<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import ActionDelete from './action-delete.svelte';
	import ActionEdit from './action-edit.svelte';
	import type { User } from './types';

	let { user }: { user: User } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionEdit {user} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {user} />
	</Layout.ActionItem>
</Layout.Actions>
