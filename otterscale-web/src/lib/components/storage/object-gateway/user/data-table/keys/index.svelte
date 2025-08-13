<script lang="ts" module>
	import type { User } from '$lib/api/storage/v1/storage_pb';
	import * as Sheet from '$lib/components/ui/sheet';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		user,
		users: users = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		user: User;
		users: Writable<User[]>;
	} = $props();
</script>

<div class="flex items-center justify-end gap-1">
	{user.keys.length}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[38vw]">
			<DataTable {selectedScopeUuid} {selectedFacility} {user} bind:users />
		</Sheet.Content>
	</Sheet.Root>
</div>
