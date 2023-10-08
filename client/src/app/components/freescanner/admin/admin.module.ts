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

import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { AppSharedModule } from '../../../shared/shared.module';
import { FreeScannerAdminComponent } from './admin.component';
import { FreeScannerAdminService } from './admin.service';
import { FreeScannerAdminConfigComponent } from './config/config.component';
import { FreeScannerAdminAccessComponent } from './config/access/access.component';
import { FreeScannerAdminApiKeysComponent } from './config/api-keys/api-keys.component';
import { FreeScannerAdminDirWatchComponent } from './config/dir-watch/dir-watch.component';
import { FreeScannerAdminDownstreamsComponent } from './config/downstreams/downstreams.component';
import { FreeScannerAdminGroupsComponent } from './config/groups/groups.component';
import { FreeScannerAdminOptionsComponent } from './config/options/options.component';
import { FreeScannerAdminSystemsSelectComponent } from './config/systems/select/select.component';
import { FreeScannerAdminSystemComponent } from './config/systems/system/system.component';
import { FreeScannerAdminSystemsComponent } from './config/systems/systems.component';
import { FreeScannerAdminTalkgroupComponent } from './config/systems/talkgroup/talkgroup.component';
import { FreeScannerAdminUnitComponent } from './config/systems/unit/unit.component';
import { FreeScannerAdminTagsComponent } from './config/tags/tags.component';
import { FreeScannerAdminLoginComponent } from './login/login.component';
import { FreeScannerAdminLogsComponent } from './logs/logs.component';
import { FreeScannerAdminTodosComponent } from './todos/todos.component';
import { FreeScannerAdminToolsComponent } from './tools/tools.component';
import { FreeScannerAdminImportExportConfigComponent } from './tools/import-export-config/import-export-config.component';
import { FreeScannerAdminImportTalkgroupsComponent } from './tools/import-talkgroups/import-talkgroups.component';
import { FreeScannerAdminImportUnitsComponent } from './tools/import-units/import-units.component';
import { FreeScannerAdminPasswordComponent } from './tools/password/password.component';

@NgModule({
    declarations: [
        FreeScannerAdminComponent,
        FreeScannerAdminConfigComponent,
        FreeScannerAdminAccessComponent,
        FreeScannerAdminApiKeysComponent,
        FreeScannerAdminDirWatchComponent,
        FreeScannerAdminDownstreamsComponent,
        FreeScannerAdminGroupsComponent,
        FreeScannerAdminImportExportConfigComponent,
        FreeScannerAdminImportTalkgroupsComponent,
        FreeScannerAdminImportUnitsComponent,
        FreeScannerAdminLoginComponent,
        FreeScannerAdminLogsComponent,
        FreeScannerAdminOptionsComponent,
        FreeScannerAdminPasswordComponent,
        FreeScannerAdminSystemComponent,
        FreeScannerAdminSystemsComponent,
        FreeScannerAdminSystemsSelectComponent,
        FreeScannerAdminTagsComponent,
        FreeScannerAdminTalkgroupComponent,
        FreeScannerAdminTodosComponent,
        FreeScannerAdminToolsComponent,
        FreeScannerAdminUnitComponent,
    ],
    entryComponents: [FreeScannerAdminSystemsSelectComponent],
    exports: [FreeScannerAdminComponent],
    imports: [AppSharedModule, HttpClientModule],
    providers: [FreeScannerAdminService],
})
export class FreeScannerAdminModule { }
