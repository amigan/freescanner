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

import { Subscription } from "rxjs";

export interface FreeScannerAvoidOptions {
    all?: boolean;
    call?: FreeScannerCall;
    minutes?: number;
    status?: boolean;
    system?: FreeScannerSystem;
    talkgroup?: FreeScannerTalkgroup;
}

export interface FreeScannerBeep {
    begin: number;
    end: number;
    frequency: number;
    type: OscillatorType;
}

export enum FreeScannerBeepStyle {
    Activate = 'activate',
    Deactivate = 'deactivate',
    Denied = 'denied',
}

export interface FreeScannerCall {
    audio?: {
        type: 'Buffer';
        data: number[];
    };
    audioName?: string;
    audioType?: string;
    dateTime: Date;
    frequencies?: FreeScannerCallFrequency[];
    frequency?: number;
    id: number;
    patches: number[];
    source?: number;
    sources?: FreeScannerCallSource[];
    system: number;
    talkgroup: number;
    talkgroupData?: FreeScannerTalkgroup;
    systemData?: FreeScannerSystem;
}

export interface FreeScannerCallFrequency {
    errorCount?: number;
    freq?: number;
    len?: number;
    pos?: number;
    spikeCount?: number;
}

export interface FreeScannerCallSource {
    pos?: number;
    src?: number;
}

export interface FreeScannerCategory {
    label: string;
    status: FreeScannerCategoryStatus;
    type: FreeScannerCategoryType;
}

export enum FreeScannerCategoryStatus {
    Off = 'off',
    On = 'on',
    Partial = 'partial',
}

export enum FreeScannerCategoryType {
    Group = 'group',
    Tag = 'tag',
}

export interface FreeScannerConfig {
    afs?: string;
    branding?: string;
    dimmerDelay: number | false;
    email?: string;
    groups: { [key: string]: { [key: number]: number[] } };
    keypadBeeps: FreeScannerKeypadBeeps | false;
    playbackGoesLive: boolean;
    showListenersCount: boolean;
    systems: FreeScannerSystem[];
    tags: { [key: string]: { [key: number]: number[] } };
    tagsToggle: boolean;
    time12hFormat: boolean;
}

export interface FreeScannerEvent {
    auth?: boolean;
    categories?: FreeScannerCategory[];
    call?: FreeScannerCall;
    config?: FreeScannerConfig;
    expired?: boolean;
    holdSys?: boolean;
    holdTg?: boolean;
    linked?: boolean;
    listeners?: number;
    livefeedMode?: FreeScannerLivefeedMode;
    map?: FreeScannerLivefeedMap;
    pause?: boolean;
    playbackList?: FreeScannerPlaybackList;
    playbackPending?: number;
    queue?: number;
    time?: number;
    tooMany?: boolean;
}

export interface FreeScannerKeypadBeeps {
    [FreeScannerBeepStyle.Activate]: FreeScannerBeep[];
    [FreeScannerBeepStyle.Deactivate]: FreeScannerBeep[];
    [FreeScannerBeepStyle.Denied]: FreeScannerBeep[];
}

export interface FreeScannerLivefeed {
    active: boolean;
    minutes: number | undefined;
    timer: Subscription | undefined;
}

export interface FreeScannerLivefeedMap {
    [key: number]: {
        [key: number]: FreeScannerLivefeed;
    };
}

export enum FreeScannerLivefeedMode {
    Offline = 'offline',
    Online = 'online',
    Playback = 'playback',
}

export interface FreeScannerPlaybackList {
    count: number;
    dateStart: Date;
    dateStop: Date;
    options: FreeScannerSearchOptions;
    results: FreeScannerCall[];
}

export interface FreeScannerSearchOptions {
    date?: Date;
    group?: string;
    limit: number;
    offset: number;
    sort: number;
    system?: number;
    tag?: string;
    talkgroup?: number;
}

export interface FreeScannerSystem {
    id: number;
    label: string;
    led?: 'blue' | 'cyan' | 'green' | 'magenta' | 'orange' | 'red' | 'white' | 'yellow';
    order?: number;
    talkgroups: FreeScannerTalkgroup[];
    units: FreeScannerUnit[];
}

export interface FreeScannerTalkgroup {
    frequency?: number;
    group: string;
    id: number;
    label: string;
    led?: 'blue' | 'cyan' | 'green' | 'magenta' | 'orange' | 'red' | 'white' | 'yellow';
    name: string;
    tag: string;
}

export interface FreeScannerUnit {
    id: number;
    label: string;
}
