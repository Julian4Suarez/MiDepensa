import { provideHttpClient } from '@angular/common/http';
import {
  type ApplicationConfig,
  inject,
  provideAppInitializer,
  provideZoneChangeDetection,
} from '@angular/core';
import { provideRouter, withComponentInputBinding } from '@angular/router';
import { provideIonicAngular } from '@ionic/angular/standalone';

import { routes } from './app.routes';
import { ConfigService } from './core/config/config.service';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    // withComponentInputBinding maps route params onto component inputs, so
    // PantryPage receives `slug` without injecting ActivatedRoute.
    provideRouter(routes, withComponentInputBinding()),
    // "md" keeps a single look on every platform instead of switching to iOS
    // styling on Apple devices.
    provideIonicAngular({ mode: 'md' }),
    provideHttpClient(),
    // Blocks bootstrap until config.json is fetched and validated.
    provideAppInitializer(() => inject(ConfigService).load()),
  ],
};
