<script lang="ts" module>
	import type { Writable } from 'svelte/store';

	import Edit from './action-edit.svelte';
	import SetDefault from './action-set-default.svelte';

	import { type Configuration, type Configuration_BootImage } from '$lib/api/configuration/v1/configuration_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		bootImage,
		configuration,
	}: {
		bootImage: Configuration_BootImage;
		configuration: Writable<Configuration>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Edit {bootImage} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<SetDefault {bootImage} {configuration} />
	</Layout.ActionItem>
</Layout.Actions>
