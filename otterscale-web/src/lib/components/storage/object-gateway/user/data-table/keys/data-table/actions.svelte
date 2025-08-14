<script lang="ts" module>
	import type { User, User_Key } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		user,
		key,
		users: data = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		user: User;
		key: User_Key;
		users: Writable<User[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Delete {selectedScopeUuid} {selectedFacility} {user} {key} bind:data />
	</Layout.ActionItem>
</Layout.Actions>
