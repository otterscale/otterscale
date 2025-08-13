import type { Application_Chart } from "$lib/api/application/v1/application_pb";
import { getIcon } from '@iconify/svelte';

const valuesMapList: Record<string, { [key: string]: string }> = {
    "minio": {
        "service.type": "NodePort",
        "service.nodePorts.api": "30001",
        "service.nodePorts.console": "30002",
    },
    "nginx": {
        "service.type": "NodePort",
        "service.nodePorts.http": "31001",
        "service.nodePorts.https": "31002",
    },
    "grafana": {
        "service.type": "NodePort",
        "service.nodePorts.grafana": "32001",
    },
    "code-server-go": {
        "codeServer.password": "password"
    },
    "code-server-python": {
        "codeServer.password": "password"
    },
};


function fuzzLogosIcon(name: string, defaultName: string): string {
    if (!name) {
        return defaultName;
    }
    let icon = `logos:${name}-icon`;
    if (getIcon(icon)) {
        return icon;
    }
    icon = `logos:${name}`;
    if (getIcon(icon)) {
        return icon;
    }
    return defaultName;
}

class FilterManager {
    charts: Application_Chart[]

    searchedName: string = $state('')
    selectedKeywords: string[] = $state([])
    selectedMaintainerNames: string[] = $state([])

    constructor(charts: Application_Chart[]) {
        this.charts = charts
    }

    get isFiltered() {
        if (this.searchedName !== '' || this.selectedKeywords.length > 0 || this.selectedMaintainerNames.length > 0) {
            return true
        }
    }

    get filteredCharts() {
        return this.charts
            .filter((chart) => (this.searchedName ? chart.name.includes(this.searchedName) : true))
            .filter((chart) => (this.selectedKeywords.length ? chart.keywords.some((keyword) => (this.selectedKeywords.includes(keyword))) : true))
            .filter((chart) => (this.selectedMaintainerNames.length ? chart.maintainers.some((maintainer) => (this.selectedMaintainerNames.includes(maintainer.name))) : true))
    }

    isKeywordSelected(keyword: string) {
        return this.selectedKeywords.includes(keyword)
    }
    isMaintainerSelected(maintainerName: string) {
        return this.selectedMaintainerNames.includes(maintainerName)
    }

    toggleKeyword(keyword: string) {
        if (this.isKeywordSelected(keyword)) {
            this.selectedKeywords = this.selectedKeywords.filter(
                (selectedKeyword) => selectedKeyword !== keyword
            );
        } else {
            this.selectedKeywords.push(keyword);
        }
    }
    toggleMaintainer(maintainerName: string) {
        if (this.isMaintainerSelected(maintainerName)) {
            this.selectedMaintainerNames =
                this.selectedMaintainerNames.filter(
                    (selectedMaintainerName) => selectedMaintainerName !== maintainerName
                );
        } else {
            this.selectedMaintainerNames.push(maintainerName);
        }
    }

    resetName() {
        this.searchedName = ''
    }
    resetKeyword() {
        this.selectedKeywords = []
    }
    resetMaintainer() {
        this.selectedMaintainerNames = []
    }
    reset() {
        this.resetName()
        this.resetKeyword()
        this.resetMaintainer()
    }
}

class PaginationManager {
    activePage: number = $state(0)
    count: number = $state(0)
    perPage: number = 6
    siblingCount: number = 1

    constructor(charts: Application_Chart[]) {
        this.count = charts.length
    }
}

export { FilterManager, fuzzLogosIcon, PaginationManager, valuesMapList };

