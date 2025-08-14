<script lang="ts" module>
	import type { User } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
	import Edit from './edit.svelte';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		user,
		users: data = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		user: User;
		users: Writable<User[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {selectedScopeUuid} {selectedFacility} {user} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {selectedScopeUuid} {selectedFacility} {user} bind:data />
	</Layout.ActionItem>
</Layout.Actions>
