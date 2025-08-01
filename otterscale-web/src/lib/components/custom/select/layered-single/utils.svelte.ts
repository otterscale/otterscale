import type { AccessorType, AncestralOptionType, OptionType } from './types';

const getAllAncestralOptions = (
    options: OptionType[],
    parents: OptionType[] = []
): AncestralOptionType[] => {
    return options.flatMap(option => {
        const ancestorOption = [...parents, option];

        if (option.subOptions) {
            return getAllAncestralOptions(option.subOptions!, [...parents, option]);
        } else {
            return [ancestorOption];
        }
    });
};

function getAncestralOptionsMap(
    options: OptionType[]
): Record<string, OptionType[]> {
    let allAncestralOptionsMap: Record<string, OptionType[]> = {};

    const allAncestralOptions = getAllAncestralOptions(options)
    allAncestralOptions.forEach((line) => {
        allAncestralOptionsMap[JSON.stringify(line.map((node) => node.value))] = line;
    });

    return allAncestralOptionsMap;
}

class OptionManager {
    options: OptionType[];
    ancestralOptionsMap: Record<string, OptionType[]>
    accesor: AccessorType

    constructor(options: OptionType[], accesor: AccessorType) {
        this.ancestralOptionsMap = getAncestralOptionsMap(options);

        this.options = options
        this.accesor = accesor;
    }

    get selectedAncestralOption(): AncestralOptionType {
        return this.ancestralOptionsMap[JSON.stringify(this.accesor.value)]
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
        return JSON.stringify(this.accesor.value) === JSON.stringify([...parents.map((parent) => (parent.value)), option.value]);
    }

    handleSelect(option: OptionType, parents: OptionType[]) {
        this.accesor.value = [...parents.map((parent) => (parent.value)), option.value]
    }
}

function validate(required: boolean | undefined, optionManager: OptionManager) {
    return required && !optionManager.selectedAncestralOption
}

export {
    getAllAncestralOptions, getAncestralOptionsMap, OptionManager, validate
};
