<script lang="ts" module>
	import * as Picker from '$lib/components/custom/picker';
	import { m } from '$lib/paraglide/messages';

	import SubvolumeGroupPicker from '../../utils/subvolume-group-picker.svelte';
</script>

<script lang="ts">
	let {
		scope,
		volume,
		selectedSubvolumeGroupName = $bindable()
	}: {
		scope: string;
		volume: string;
		selectedSubvolumeGroupName: string;
	} = $props();
</script>

<Picker.Root align="left">
	<Picker.Wrapper class="*:h-8">
		<Picker.Label>{m.group()}</Picker.Label>
		<SubvolumeGroupPicker {scope} {volume} bind:selectedSubvolumeGroupName />
	</Picker.Wrapper>
</Picker.Root>
