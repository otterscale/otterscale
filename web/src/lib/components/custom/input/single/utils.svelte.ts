import { z, type ZodFirstPartySchemaTypes } from 'zod';

const BORDER_INPUT_CLASSNAME = 'border-input placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent text-base shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:outline-none focus-visible:ring-1 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm';
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
    BORDER_INPUT_CLASSNAME, InputValidator, PasswordManager, typeToIcon
};
export type {
    InputValidatorResponse
};

