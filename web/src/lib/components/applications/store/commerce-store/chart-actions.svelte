<script lang="ts" module>
	import type { Writable } from 'svelte/store';

	import type { Release } from '$lib/api/application/v1/application_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import { m } from '$lib/paraglide/messages';

	import Delete from './chart-action-delete-release.svelte';
	import Edit from './chart-action-edit-release.svelte';
	import Rollback from './chart-action-rollback-release.svelte';
</script>

<script lang="ts">
	let {
		release,
		scope,
		releases
	}: {
		release: Release;
		scope: string;
		releases: Writable<Release[]>;
	} = $props();

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Actions.List bind:open>
	<Actions.Label>
		{m.actions()}
	</Actions.Label>
	<Actions.Separator />
	<Actions.Item>
		<Edit {release} {scope} {releases} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Rollback {release} {scope} {releases} closeActions={close} />
	</Actions.Item>
	<Actions.Separator />
	<Actions.Item>
		<Delete {release} {scope} {releases} closeActions={close} />
	</Actions.Item>
</Actions.List>
