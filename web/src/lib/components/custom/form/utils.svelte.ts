import type { Invalidaties } from './type'

class FormValidator {
    invalidaties = $state<Invalidaties>({});

    set(id: string | null | undefined, value: boolean | undefined) {
        if (id) {
            this.invalidaties[id] = value;
        }
    }

    get isInvalid() {
        return Object.values(this.invalidaties).some((v) => (v))
    }
}

export { FormValidator }
