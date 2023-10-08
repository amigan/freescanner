/**
 * @license
 * Copyright Saubeo Solutions All Rights Reserved
 *
 * This source code is proprietary and confidential
 * Unauthorized copying of this file, via any medium is strictly prohibited
 */

import { Routes } from '@angular/router';
import { routes as freeScannerRoutes } from './pages/freescanner';

export const routes: Routes = [
    ...freeScannerRoutes,
    {
        path: '**',
        pathMatch: 'full',
        redirectTo: '',
    },
];
