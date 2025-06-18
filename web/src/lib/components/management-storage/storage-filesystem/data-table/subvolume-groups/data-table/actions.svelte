<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import ActionEdit from './action-edit.svelte';
	import ActionDelete from './action-delete.svelte';
	import type { SubvolumeGroup } from './types';

	let { subvolumeGroup }: { subvolumeGroup: SubvolumeGroup } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionEdit {subvolumeGroup} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {subvolumeGroup} />
	</Layout.ActionItem>
</Layout.Actions>
