/*
 * *****************************************************************************
 * Copyright (C) 2019-2022 Chrystian Huot <chrystian.huot@saubeo.solutions>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>
 * ****************************************************************************
 */

import { Component, OnDestroy } from '@angular/core';
import {
    FreeScannerAvoidOptions,
    FreeScannerBeepStyle,
    FreeScannerCategory,
    FreeScannerCategoryStatus,
    FreeScannerEvent,
    FreeScannerLivefeedMap,
    FreeScannerSystem,
} from '../freescanner';
import { FreeScannerService } from '../freescanner.service';

@Component({
    selector: 'freescanner-select',
    styleUrls: [
        '../common.scss',
        './select.component.scss',
    ],
    templateUrl: './select.component.html',
})
export class FreeScannerSelectComponent implements OnDestroy {
    categories: FreeScannerCategory[] | undefined;

    map: FreeScannerLivefeedMap = {};

    systems: FreeScannerSystem[] | undefined;

    tagsToggle: boolean | undefined;

    private eventSubscription = this.freeScannerService.event.subscribe((event: FreeScannerEvent) => this.eventHandler(event));

    constructor(private freeScannerService: FreeScannerService) { }

    avoid(options?: FreeScannerAvoidOptions): void {
        if (options?.all == true) {
            this.freeScannerService.beep(FreeScannerBeepStyle.Activate);

        } else if (options?.all == false) {
            this.freeScannerService.beep(FreeScannerBeepStyle.Deactivate);

        } else if (options?.system !== undefined && options?.talkgroup !== undefined) {
            this.freeScannerService.beep(this.map[options!.system.id][options!.talkgroup.id].active
                ? FreeScannerBeepStyle.Deactivate
                : FreeScannerBeepStyle.Activate
            );

        } else {
            this.freeScannerService.beep(options?.status ? FreeScannerBeepStyle.Activate : FreeScannerBeepStyle.Deactivate);
        }

        this.freeScannerService.avoid(options);
    }

    ngOnDestroy(): void {
        this.eventSubscription.unsubscribe();
    }

    toggle(category: FreeScannerCategory): void {
        if (category.status == FreeScannerCategoryStatus.On)
            this.freeScannerService.beep(FreeScannerBeepStyle.Deactivate);
        else
            this.freeScannerService.beep(FreeScannerBeepStyle.Activate);

        this.freeScannerService.toggleCategory(category);
    }

    private eventHandler(event: FreeScannerEvent): void {
        if (event.config) {
            this.tagsToggle = event.config.tagsToggle;
            this.systems = event.config.systems;
        }
        if (event.categories) this.categories = event.categories;
        if (event.map) this.map = event.map;
    }
}
