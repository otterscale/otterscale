<script lang="ts">
	import type { FormState } from '@sjsf/form';
	import { Content, Form, setFormContext } from '@sjsf/form';

	interface Props {
		form: FormState<Record<string, unknown>>;
	}

	let { form }: Props = $props();

	setFormContext(form);
</script>

<Form>
	<Content />
	<slot />
</Form>
