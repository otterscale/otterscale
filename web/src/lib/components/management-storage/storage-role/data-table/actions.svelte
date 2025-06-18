<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import ActionDelete from './action-delete.svelte';
	import ActionEdit from './action-edit.svelte';
	import type { Role } from './types';

	let { role }: { role: Role } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionEdit {role} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {role} />
	</Layout.ActionItem>
</Layout.Actions>
