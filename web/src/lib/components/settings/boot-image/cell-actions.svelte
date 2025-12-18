<script lang="ts" module>
	import type { Writable } from 'svelte/store';

	import {
		type Configuration,
		type Configuration_BootImage
	} from '$lib/api/configuration/v1/configuration_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Edit from './action-edit.svelte';
	import SetDefault from './action-set-default.svelte';
</script>

<script lang="ts">
	let {
		bootImage,
		configuration,
		reloadManager
	}: {
		bootImage: Configuration_BootImage;
		configuration: Writable<Configuration>;
		reloadManager: ReloadManager;
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
		<Edit {bootImage} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<SetDefault {bootImage} {configuration} closeActions={close} />
	</Actions.Item>
</Actions.List>
