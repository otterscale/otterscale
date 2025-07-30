<script lang="ts" module>
	import * as Picker from '$lib/components/custom/picker';
	import FioOutputTypePicker from '$lib/components/bist/utils/fio-output-type-pickers.svelte';
</script>

<script lang="ts">
    let { selectedMode = $bindable() }: { selectedMode: string } = $props();
</script>

<Picker.Root align="left" class="mt-2">
	<Picker.Wrapper class="*:h-8">
		<Picker.Label>Mode</Picker.Label>
		<FioOutputTypePicker bind:selectedMode={selectedMode} />
	</Picker.Wrapper>
</Picker.Root>
