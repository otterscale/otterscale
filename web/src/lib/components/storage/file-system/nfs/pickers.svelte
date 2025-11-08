<script lang="ts" module>
	import * as Picker from '$lib/components/custom/picker';
	import { m } from '$lib/paraglide/messages';

	import SubvolumeGroupPicker from '../../utils/subvolume-group-picker.svelte';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		selectedVolume,
		selectedSubvolumeGroupName = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroupName: string;
	} = $props();
</script>

<Picker.Root align="left">
	<Picker.Wrapper class="*:h-8">
		<Picker.Label>{m.group()}</Picker.Label>
		<SubvolumeGroupPicker
			{selectedScope}
			{selectedFacility}
			{selectedVolume}
			bind:selectedSubvolumeGroupName
		/>
	</Picker.Wrapper>
</Picker.Root>
