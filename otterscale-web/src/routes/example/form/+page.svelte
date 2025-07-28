<script lang="ts" module>
	import * as Form from '$lib/components/custom/form';
	import Button from '$lib/components/ui/button/button.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import { Single as SingleInput } from '$lib/components/custom/input';
</script>

<script lang="ts">
	let isInvalid = $state(false);
</script>

<Form.Root class="group max-h-none">
	<Form.Fieldset>
		<Form.Legend>Fieldset</Form.Legend>
		<Form.Description>
			This is a long description for the form fieldset. It provides detailed information about the
			purpose of this section.
		</Form.Description>
		<Form.Field>
			<Form.Label>
				Field
				{#snippet information()}
					This is a detailed description for the form field. It provides additional context about
					the purpose and usage of this field.
				{/snippet}
			</Form.Label>
			<SingleInput.General required id="field" />
			<Form.Help>This is a help text for Field. Please enter the required information.</Form.Help>
		</Form.Field>
	</Form.Fieldset>

	<Form.Separator>Or</Form.Separator>

	<Form.Fieldset>
		<Form.Legend>Fieldset</Form.Legend>
		<Form.Description>
			This is a long description for the form fieldset. It provides detailed information about the
			purpose of this section.
		</Form.Description>
		<Form.Field>
			<Form.Label>
				Field
				{#snippet information()}
					This is a detailed description for the form field. It provides additional context about
					the purpose and usage of this field.
				{/snippet}
			</Form.Label>
			<Form.Help>This is a help text for Field. Please enter the required information.</Form.Help>
		</Form.Field>
		<Form.Field>
			<Form.Label>
				Field
				{#snippet information()}
					This is a detailed description for the form field. It provides additional context about
					the purpose and usage of this field.
				{/snippet}
			</Form.Label>
			<Input />
			<Form.Help>This is a help text for Field. Please enter the required information.</Form.Help>
		</Form.Field>
	</Form.Fieldset>

	<Form.Actions>
		<Form.Submit>Confirm</Form.Submit>
	</Form.Actions>
</Form.Root>
