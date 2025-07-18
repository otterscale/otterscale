import { z, type ZodFirstPartySchemaTypes } from 'zod';

const BORDER_INPUT_CLASSNAME = 'flex items-center rounded-md border shadow-sm';
const UNFOCUS_INPUT_CLASSNAME = 'border-none shadow-none focus-visible:ring-0 bg-transparent';
const RING_INVALID_INPUT_CLASSNAME = 'ring-destructive ring-1'
const RING_VALID_INPUT_CLASSNAME = 'ring-1'
const typeToIcon: Record<string, string> = {
    color: 'ph:palette',
    'datetime-local': 'ph:clock',
    date: 'ph:calendar',
    time: 'ph:clock',
    url: 'ph:link',
    email: 'ph:mailbox',
    tel: 'ph:phone',
    boolean: 'ph:toggle-left',
    text: 'ph:textbox',
    number: 'ph:list-numbers',
    search: 'ph:magnifying-glass',
    password: 'ph:password'
};

class PasswordManager {
    isVisible = $state<boolean>(false);

    enable() {
        this.isVisible = true;
    }

    disable() {
        this.isVisible = false;
    }
}

type InputValidatorResponse = {
    isValid: boolean;
    errors: z.ZodIssue[]
}
class InputValidator {
    schema: ZodFirstPartySchemaTypes | undefined

    constructor(schema: ZodFirstPartySchemaTypes | undefined) {
        this.schema = schema
    }

    validate(input: any) {
        if (!this.schema) { return { isValid: true, errors: undefined } }

        const result = this.schema.safeParse(input)
        if (result.success) return { isValid: true, errors: [] } as InputValidatorResponse
        return { isValid: false, errors: result.error.errors } as InputValidatorResponse
    }
}

export {
    BORDER_INPUT_CLASSNAME, InputValidator, PasswordManager, RING_INVALID_INPUT_CLASSNAME, RING_VALID_INPUT_CLASSNAME, typeToIcon, UNFOCUS_INPUT_CLASSNAME
};
export type {
    InputValidatorResponse
};

