import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import type { Observable } from 'rxjs';

import { ConfigService } from '../config/config.service';
import type { ItemPatch, Pantry, PantryDetail, PantryItem } from '../../shared/models/pantry.model';

/** Thin HTTP client for the pantry endpoints. */
@Injectable({ providedIn: 'root' })
export class PantryApiService {
  private readonly http = inject(HttpClient);
  private readonly config = inject(ConfigService);

  /** Creates a pantry from a free-text name; the API derives the slug. */
  createPantry(name: string): Observable<Pantry> {
    return this.http.post<Pantry>(`${this.config.backendUrl}/pantries`, { name });
  }

  /** Loads a pantry with every item. */
  getPantry(slug: string): Observable<PantryDetail> {
    return this.http.get<PantryDetail>(`${this.config.backendUrl}/pantries/${slug}`);
  }

  /** Updates the status, view or category of a single item. */
  updateItem(slug: string, productId: string, patch: ItemPatch): Observable<PantryItem> {
    return this.http.patch<PantryItem>(
      `${this.config.backendUrl}/pantries/${slug}/items/${productId}`,
      patch,
    );
  }

  /** Resets every active product to the pending state. */
  resetActiveItems(slug: string): Observable<void> {
    return this.http.post<void>(`${this.config.backendUrl}/pantries/${slug}/items/reset`, null);
  }
}
