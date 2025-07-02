<script lang="ts">
	import { Tabs } from '$lib/components/custom/tabs/index';
	import { TestingTools } from './testing-tools';
	import Icon from '@iconify/svelte';
</script>

<Tabs.Root value="testing-tools">
	<Tabs.List class="*:gap-2 [*_&_svg]:size-5">
		<Tabs.Trigger value="testing-tools">
			<Icon icon="ph:cube" />
			TESTING TOOLS
		</Tabs.Trigger>
	</Tabs.List>
	<Tabs.Content value="testing-tools">
		<TestingTools />
	</Tabs.Content>
</Tabs.Root>
