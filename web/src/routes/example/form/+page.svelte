<script lang="ts" module>
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput, Multiple as MultipleInput } from '$lib/components/custom/input';
	import { PalinTabs } from '$lib/components/custom/tabs';
</script>

<script lang="ts">
	import Icon from '@iconify/svelte';

	let field1 = $state(false);
	let field2 = $state(false);
	let field3 = $state(false);
</script>

<Form.Root class="max-h-none">
	<PalinTabs.Root value="basic">
		<PalinTabs.List>
			<PalinTabs.Trigger value="basic">
				<Icon icon="ph:atom" />
				Basic
			</PalinTabs.Trigger>
			<PalinTabs.Trigger value="advanced">
				<Icon icon="ph:planet" />
				Advanced
			</PalinTabs.Trigger>
		</PalinTabs.List>
		<PalinTabs.Content value="basic">
			<Form.Fieldset>
				<Form.Legend>Fieldset</Form.Legend>
				<Form.Description>
					This is a long description for the form fieldset. It provides detailed information about the purpose
					of this section.
				</Form.Description>
				<Form.Field>
					<Form.Label>
						Field
						{#snippet information()}
							This is a detailed description for the form field. It provides additional context about the
							purpose and usage of this field.
						{/snippet}
					</Form.Label>
					<Form.Help>This is a help text for Field. Please enter the required information.</Form.Help>
					<SingleInput.General required bind:invalid={field1} />
				</Form.Field>
			</Form.Fieldset>
		</PalinTabs.Content>

		<PalinTabs.Content value="advanced">
			<Form.Fieldset>
				<Form.Legend>Fieldset</Form.Legend>
				<Form.Description>
					This is a long description for the form fieldset. It provides detailed information about the purpose
					of this section.
				</Form.Description>
				<Form.Field>
					<Form.Label>
						Field
						{#snippet information()}
							This is a detailed description for the form field. It provides additional context about the
							purpose and usage of this field.
						{/snippet}
					</Form.Label>
					<Form.Help>This is a help text for Field. Please enter the required information.</Form.Help>
				</Form.Field>
				<Form.Field>
					<Form.Label>
						Field
						{#snippet information()}
							This is a detailed description for the form field. It provides additional context about the
							purpose and usage of this field.
						{/snippet}
					</Form.Label>
					<Form.Help>This is a help text for Field. Please enter the required information.</Form.Help>
					<SingleInput.General required bind:invalid={field2} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Separator>Or</Form.Separator>

			<Form.Fieldset>
				<Form.Legend>Fieldset</Form.Legend>
				<Form.Description>
					This is a long description for the form fieldset. It provides detailed information about the purpose
					of this section.
				</Form.Description>
				<Form.Field>
					<Form.Label>
						Field
						{#snippet information()}
							This is a detailed description for the form field. It provides additional context about the
							purpose and usage of this field.
						{/snippet}
					</Form.Label>
					<Form.Help>This is a help text for Field. Please enter the required information.</Form.Help>
					<MultipleInput.Root required bind:invalid={field3}>
						<MultipleInput.Viewer />
						<MultipleInput.Controller>
							<MultipleInput.Input />
							<MultipleInput.Add />
							<MultipleInput.Clear />
						</MultipleInput.Controller>
					</MultipleInput.Root>
				</Form.Field>
			</Form.Fieldset>
		</PalinTabs.Content>
	</PalinTabs.Root>

	<Form.Actions>
		<Form.Submit disabled={field1 || field2 || field3}>Confirm</Form.Submit>
	</Form.Actions>
</Form.Root>
