import type { InputType } from './types';

import { typeToIcon, PasswordManager } from './utils.svelte';

import { default as InputGeneral } from './input-general.svelte';
import { default as InputPassword } from './input-password.svelte';
import { default as InputColor } from './input-color.svelte';
import { default as InputBoolean } from './input-boolean.svelte';

export {
    type InputType,
    //
    typeToIcon,
    //
    PasswordManager,
    //
    InputGeneral as General,
    InputPassword as Password,
    InputColor as Color,
    InputBoolean as Boolean,
};
