<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { writable } from 'svelte/store';

	import * as Picker from '$lib/components/custom/picker';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let { selectedMode = $bindable() }: { selectedMode: string } = $props();

	const modes = writable<SingleSelect.OptionType[]>([
		{ value: 'get', label: 'Get', icon: 'ph:download-simple' },
		{ value: 'put', label: 'Put', icon: 'ph:upload-simple' },
		{ value: 'delete', label: 'Delete', icon: 'ph:trash' },
	]);
</script>

<Picker.Root align="left" class="mt-2">
	<Picker.Wrapper class="*:h-8">
		<Picker.Label>{m.mode()}</Picker.Label>
		<SingleSelect.Root options={modes} bind:value={selectedMode}>
			<SingleSelect.Trigger />
			<SingleSelect.Content>
				<SingleSelect.Options>
					<SingleSelect.Input />
					<SingleSelect.List>
						<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
						<SingleSelect.Group>
							{#each $modes as option}
								<SingleSelect.Item {option}>
									<Icon
										icon={option.icon ? option.icon : 'ph:empty'}
										class={cn('size-5', option.icon ? 'visible' : 'invisible')}
									/>
									{option.label}
									<SingleSelect.Check {option} />
								</SingleSelect.Item>
							{/each}
						</SingleSelect.Group>
					</SingleSelect.List>
				</SingleSelect.Options>
			</SingleSelect.Content>
		</SingleSelect.Root>
	</Picker.Wrapper>
</Picker.Root>
