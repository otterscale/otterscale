import type { OptionType, AncestralOptionType, valuesSetterType } from './types';

class OptionManager {
    visibility = 1
    options: OptionType[];

    selectedAncestralOptions: AncestralOptionType[];
    valuesSetter: valuesSetterType;

    constructor(options: OptionType[], selectedAncestralOptions: AncestralOptionType[], valuesSetter: valuesSetterType) {
        this.options = options;
        this.selectedAncestralOptions = $state(selectedAncestralOptions);
        this.valuesSetter = valuesSetter;
    }

    get isSomeAncestralOptionsSelected(): boolean {
        return Array.isArray(this.selectedAncestralOptions) && this.selectedAncestralOptions.length > 0;
    }

    hasSubOptions(option: OptionType): boolean {
        if (!option.subOptions) {
            return false;
        }
        if (option.subOptions.length == 0) {
            return false;
        }
        return true;
    }

    isOptionSelected(option: OptionType, parents: OptionType[]): boolean {
        const selectedAncestralOption: AncestralOptionType = [...parents, option]
        return this.selectedAncestralOptions.some(
            (o) => JSON.stringify(o) === JSON.stringify(selectedAncestralOption)
        );
    }

    handleSelect(option: OptionType, parents: OptionType[]) {
        const selectedAncestralOption: AncestralOptionType = [...parents, option];

        if (Array.isArray(this.selectedAncestralOptions) && this.isOptionSelected(option, parents)) {
            this.selectedAncestralOptions = this.selectedAncestralOptions.filter((o) => JSON.stringify(o) !== JSON.stringify(selectedAncestralOption));
        } else {
            this.selectedAncestralOptions = [
                ...this.selectedAncestralOptions,
                selectedAncestralOption
            ];
        }

        this.valuesSetter(this.selectedAncestralOptions);
    }

    all() {
        const getAllAncestralOptions = (
            options: OptionType[],
            parents: OptionType[] = []
        ): AncestralOptionType[] => {
            return options.flatMap(option => {
                const ancestorOption = [...parents, option];

                if (this.hasSubOptions(option)) {
                    return getAllAncestralOptions(option.subOptions!, [...parents, option]);
                } else {
                    return [ancestorOption];
                }
            });
        };

        this.selectedAncestralOptions = getAllAncestralOptions(this.options);
        this.valuesSetter(this.selectedAncestralOptions);
    }

    clear() {
        this.selectedAncestralOptions = [];
        this.valuesSetter(this.selectedAncestralOptions);
    }
}

export {
    OptionManager
}