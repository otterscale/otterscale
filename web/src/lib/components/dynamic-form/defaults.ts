export { createFormIdBuilder as idBuilder } from '@sjsf/form/id-builders/modern';
export { createFormMerger as merger } from '@sjsf/form/mergers/modern';
export { resolver } from '@sjsf/form/resolvers/basic';
export { translation } from '@sjsf/form/translations/en';
export { icons } from '@sjsf/lucide-icons';
export { theme } from '@sjsf/shadcn4-theme';

import { addFormComponents, createFormValidator } from '@sjsf/ajv8-validator';
import type { ValidatorFactoryOptions } from '@sjsf/form';
import addFormats from 'ajv-formats';
export const validator = <T>(options: ValidatorFactoryOptions) => (
    createFormValidator<T>({
        ...options,
        ajvPlugins: (ajv) => addFormComponents(addFormats(ajv))
    })
);

import '@sjsf/form/fields/extra/aggregated-include';
import '@sjsf/form/fields/extra/array-files-include';
import '@sjsf/form/fields/extra/array-native-files-include';
import '@sjsf/form/fields/extra/array-tags-include';
import '@sjsf/form/fields/extra/boolean-select-include';
import '@sjsf/form/fields/extra/enum-include';
import '@sjsf/form/fields/extra/file-include';
import '@sjsf/form/fields/extra/files-include';
import '@sjsf/form/fields/extra/multi-enum-include';
import '@sjsf/form/fields/extra/native-file-include';
import '@sjsf/form/fields/extra/native-files-include';
import '@sjsf/form/fields/extra/tags-include';
import '@sjsf/form/fields/extra/unknown-native-file-include';

import "@sjsf/basic-theme/extra-widgets/checkboxes-include"
import "@sjsf/basic-theme/extra-widgets/date-picker-include"
import "@sjsf/basic-theme/extra-widgets/file-include"
import "@sjsf/basic-theme/extra-widgets/multi-select-include"
import "@sjsf/basic-theme/extra-widgets/radio-include"
import "@sjsf/basic-theme/extra-widgets/range-include"
import "@sjsf/basic-theme/extra-widgets/textarea-include"

import '@sjsf/shadcn4-theme/extra-widgets/checkboxes-include';
import '@sjsf/shadcn4-theme/extra-widgets/combobox-include';
import '@sjsf/shadcn4-theme/extra-widgets/date-picker-include';
import '@sjsf/shadcn4-theme/extra-widgets/date-range-picker-include';
import '@sjsf/shadcn4-theme/extra-widgets/file-include';
import '@sjsf/shadcn4-theme/extra-widgets/multi-select-include';
import '@sjsf/shadcn4-theme/extra-widgets/radio-buttons-include';
import '@sjsf/shadcn4-theme/extra-widgets/radio-include';
import '@sjsf/shadcn4-theme/extra-widgets/range-include';
import '@sjsf/shadcn4-theme/extra-widgets/range-slider-include';
import '@sjsf/shadcn4-theme/extra-widgets/switch-include';
import '@sjsf/shadcn4-theme/extra-widgets/textarea-include';
