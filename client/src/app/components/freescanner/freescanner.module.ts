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

import { FullscreenOverlayContainer, OverlayContainer } from '@angular/cdk/overlay';
import { NgModule } from '@angular/core';
import { AppSharedModule } from '../../shared/shared.module';
import { FreeScannerComponent } from './freescanner.component';
import { FreeScannerService } from './freescanner.service';
import { FreeScannerMainComponent } from './main/main.component';
import { FreeScannerSupportComponent } from './main/support/support.component';
import { FreeScannerSearchComponent } from './search/search.component';
import { FreeScannerSelectComponent } from './select/select.component';

@NgModule({
    declarations: [
        FreeScannerComponent,
        FreeScannerMainComponent,
        FreeScannerSearchComponent,
        FreeScannerSelectComponent,
        FreeScannerSupportComponent,
    ],
    exports: [FreeScannerComponent],
    imports: [
        AppSharedModule,
    ],
    providers: [
        FreeScannerService,
        { provide: OverlayContainer, useClass: FullscreenOverlayContainer },
    ],
})
export class FreeScannerModule { }
