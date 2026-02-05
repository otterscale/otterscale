export { resolver } from '@sjsf/form/resolvers/basic';
export { translation } from '@sjsf/form/translations/en';
import '@sjsf/form/fields/extra/enum-include';
import '@sjsf/form/fields/extra/multi-enum-include';
import '@sjsf/form/fields/extra/unknown-native-file-include';

import { overrideByRecord } from '@sjsf/form/lib/resolver';
import { theme as shadcnTheme } from '@sjsf/shadcn4-theme';

import CustomLayout from './layout/Layout.svelte';

export const theme = overrideByRecord(shadcnTheme, {
	layout: CustomLayout
});
import '@sjsf/shadcn4-theme/extra-widgets/textarea-include';
import '@sjsf/shadcn4-theme/extra-widgets/checkboxes-include';
import '@sjsf/shadcn4-theme/extra-widgets/radio-include';
import '@sjsf/shadcn4-theme/extra-widgets/file-include';
import '@sjsf/shadcn4-theme/extra-widgets/date-picker-include';
import '@sjsf/shadcn4-theme/extra-widgets/switch-include';

export { createFormIdBuilder as idBuilder } from '@sjsf/form/id-builders/modern';
export { createFormMerger as merger } from '@sjsf/form/mergers/modern';

import { addFormComponents, createFormValidator } from '@sjsf/ajv8-validator';
import type { ValidatorFactoryOptions } from '@sjsf/form';
import addFormats from 'ajv-formats';

export const validator = <T>(options: ValidatorFactoryOptions) =>
	createFormValidator<T>({
		...options,
		ajvPlugins: (ajv) => addFormComponents(addFormats(ajv))
	});
