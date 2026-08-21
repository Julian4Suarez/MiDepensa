import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { firstValueFrom } from 'rxjs';

import type { AppConfig } from './config.model';

/**
 * Loads `/config.json` before the app bootstraps.
 *
 * The file is written by `entrypoint.sh` when the container starts, so the same
 * built image can point at any backend without being rebuilt.
 */
@Injectable({ providedIn: 'root' })
export class ConfigService {
  private readonly http = inject(HttpClient);
  private config: AppConfig | null = null;

  /** Fetches and validates the runtime configuration. */
  async load(): Promise<void> {
    const loaded = await firstValueFrom(this.http.get<AppConfig>('config.json'));

    if (!loaded?.BACKEND_URL || !/^https?:\/\/.+/.test(loaded.BACKEND_URL)) {
      throw new Error('config.json: BACKEND_URL must be an absolute http(s) URL');
    }
    this.config = { BACKEND_URL: loaded.BACKEND_URL.replace(/\/+$/, '') };
  }

  /** Base URL of the API. Throws when accessed before {@link load} resolves. */
  get backendUrl(): string {
    if (!this.config) {
      throw new Error('ConfigService used before initialization');
    }
    return this.config.BACKEND_URL;
  }
}
