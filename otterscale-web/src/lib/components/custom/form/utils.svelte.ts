import type { Invalidity } from './type';

class FormValidator {
    invalidity = $state<Invalidity>({});

    set(id: string | null | undefined, value: boolean | null | undefined) {
        if (id) {
            this.invalidity[id] = value;
        }
    }

    get isInvalid() {
        return Object.values(this.invalidity).some((v) => (v))
    }
}

export { FormValidator };
